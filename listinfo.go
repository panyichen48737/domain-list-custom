package main

import (
	"bufio"
	"errors"
	"os"
	"sort"
	"strings"

	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
)

// includeRule is a single `include:` line: the source file it pulls in together
// with the attribute constraints written on that same line.
// `include:x @cn`        -> must = ["@cn"],  ban = []
// `include:x @-!cn`      -> must = [],       ban = ["@!cn"] (upstream fn syntax:
//                           take everything that does NOT carry `!cn`)
// `include:x @cn @-ads`  -> must = ["@cn"],  ban = ["@ads"]
// Constraints of one line are combined, they are not applied one after another.
type includeRule struct {
	source fileName
	must   []attribute
	ban    []attribute
}

// ListInfo is the information structure of a single file in data directory.
// It includes all types of rules of the file, as well as servel types of
// sturctures of same items for convenience in later process.
type ListInfo struct {
	Name                    fileName
	HasInclusion            bool
	Inclusions              []includeRule
	FullTypeList            []*router.Domain
	KeywordTypeList         []*router.Domain
	RegexpTypeList          []*router.Domain
	AttributeRuleUniqueList []*router.Domain
	DomainTypeList          []*router.Domain
	DomainTypeUniqueList    []*router.Domain
	AttributeRuleListMap    map[attribute][]*router.Domain
	GeoSite                 *router.GeoSite
}

// NewListInfo return a ListInfo
func NewListInfo() *ListInfo {
	return &ListInfo{
		FullTypeList:            make([]*router.Domain, 0, 10),
		KeywordTypeList:         make([]*router.Domain, 0, 10),
		RegexpTypeList:          make([]*router.Domain, 0, 10),
		AttributeRuleUniqueList: make([]*router.Domain, 0, 10),
		DomainTypeList:          make([]*router.Domain, 0, 10),
		DomainTypeUniqueList:    make([]*router.Domain, 0, 10),
		AttributeRuleListMap:    make(map[attribute][]*router.Domain),
	}
}

// ProcessList processes each line of every single file in the data directory
// and generates a ListInfo of each file.
func (l *ListInfo) ProcessList(file *os.File) error {
	scanner := bufio.NewScanner(file)
	// Parse a file line by line to generate ListInfo
	for scanner.Scan() {
		line := scanner.Text()
		if isEmpty(line) {
			continue
		}
		line = removeComment(line)
		if isEmpty(line) {
			continue
		}
		parsedRule, err := l.parseRule(line)
		if err != nil {
			return err
		}
		if parsedRule == nil {
			continue
		}
		l.classifyRule(parsedRule)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// parseRule parses a single rule
func (l *ListInfo) parseRule(line string) (*router.Domain, error) {
	line = strings.TrimSpace(line)

	if line == "" {
		return nil, errors.New("empty line")
	}

	// Parse `include` rule first, eg: `include:google`, `include:google @cn @gfw`
	if strings.HasPrefix(line, "include:") {
		return nil, l.parseInclusion(line)
	}

	// Fields, not Split(" "): a tab between the rule and its attributes would
	// otherwise stay glued to the rule and turn `@cn` into part of the value.
	parts := strings.Fields(line)
	ruleWithType := strings.TrimSpace(parts[0])
	if ruleWithType == "" {
		return nil, errors.New("empty rule")
	}

	var rule router.Domain
	if err := l.parseTypeRule(ruleWithType, &rule); err != nil {
		return nil, err
	}

	for _, attrString := range parts[1:] {
		if attrString = strings.TrimSpace(attrString); attrString != "" {
			attr, err := l.parseAttribute(attrString)
			if err != nil {
				return nil, err
			}
			rule.Attribute = append(rule.Attribute, attr)
		}
	}

	return &rule, nil
}

func (l *ListInfo) parseInclusion(inclusion string) error {
	inclusionVal := strings.TrimPrefix(strings.TrimSpace(inclusion), "include:")
	l.HasInclusion = true
	inclusionValSlice := strings.Split(inclusionVal, "@")
	inc := includeRule{source: fileName(strings.ToUpper(strings.TrimSpace(inclusionValSlice[0])))}
	// support new inclusion syntax, eg: `include:google @cn @gfw`
	// and the ban syntax used upstream, eg: `include:tencent @-!cn`
	//
	// An attribute with no name, as in `include:child @` or `include:child @-`,
	// is rejected rather than ignored: ignoring it leaves a rule with neither a
	// wanted nor a banned attribute, which is the "take everything" form of the
	// line and would silently widen the include to the whole child list.
	for _, attr := range inclusionValSlice[1:] {
		attr = strings.ToLower(strings.TrimSpace(attr))
		if ban, isBan := strings.CutPrefix(attr, "-"); isBan {
			if ban = strings.TrimSpace(ban); ban == "" {
				return errors.New("invalid ban attribute in inclusion: " + inclusion)
			}
			inc.ban = append(inc.ban, attribute("@"+ban))
			continue
		}
		if attr == "" {
			return errors.New("invalid attribute in inclusion: " + inclusion)
		}
		// Added in this format: '@cn'
		inc.must = append(inc.must, attribute("@"+attr))
	}
	l.Inclusions = append(l.Inclusions, inc)
	return nil
}

func (l *ListInfo) parseTypeRule(domain string, rule *router.Domain) error {
	// SplitN, not Split: a regexp value may itself contain colons, as in
	// `regexp:^https?://example\.com`, and splitting on every one of them used
	// to leave neither branch below matching, silently dropping the rule.
	kv := strings.SplitN(domain, ":", 2)
	switch len(kv) {
	case 1: // line without type prefix
		rule.Type = router.Domain_RootDomain
		rule.Value = strings.ToLower(strings.TrimSpace(kv[0]))
	case 2: // line with type prefix
		ruleType := strings.TrimSpace(kv[0])
		ruleVal := strings.TrimSpace(kv[1])
		rule.Value = strings.ToLower(ruleVal)
		switch strings.ToLower(ruleType) {
		case "full":
			rule.Type = router.Domain_Full
		case "domain":
			rule.Type = router.Domain_RootDomain
		case "keyword":
			rule.Type = router.Domain_Plain
		case "regexp":
			rule.Type = router.Domain_Regex
			rule.Value = ruleVal
		default:
			return errors.New("unknown domain type: " + ruleType)
		}
	}
	return nil
}

func (l *ListInfo) parseAttribute(attr string) (*router.Domain_Attribute, error) {
	if attr[0] != '@' {
		return nil, errors.New("invalid attribute: " + attr)
	}
	// Trim out attribute prefix `@` character. An attribute with no name, as in
	// `domain:example.com @`, is rejected for the same reason the include side
	// rejects it: it produces a rule that no include filter can select.
	key := strings.ToLower(strings.TrimSpace(attr[1:]))
	if key == "" {
		return nil, errors.New("invalid attribute: " + attr)
	}

	var attribute router.Domain_Attribute
	attribute.Key = key
	attribute.TypedValue = &router.Domain_Attribute_BoolValue{BoolValue: true}
	return &attribute, nil
}

// classifyRule classifies a single rule and write into *ListInfo
func (l *ListInfo) classifyRule(rule *router.Domain) {
	if len(rule.Attribute) > 0 {
		l.AttributeRuleUniqueList = append(l.AttributeRuleUniqueList, rule)
		attrsString := attributesKey(rule) // attrsString will be "@cn@ads" if there are more than one attributes
		l.AttributeRuleListMap[attrsString] = append(l.AttributeRuleListMap[attrsString], rule)
	} else {
		switch rule.Type {
		case router.Domain_Full:
			l.FullTypeList = append(l.FullTypeList, rule)
		case router.Domain_RootDomain:
			l.DomainTypeList = append(l.DomainTypeList, rule)
		case router.Domain_Plain:
			l.KeywordTypeList = append(l.KeywordTypeList, rule)
		case router.Domain_Regex:
			l.RegexpTypeList = append(l.RegexpTypeList, rule)
		}
	}
}

// Flatten flattens the rules in a file that have "include" syntax
// in data directory, and adds those need-to-included rules into it.
// This feature supports the "include:filename@attribute" syntax.
// It also generates a domain trie of domain-typed rules for each file
// to remove duplications of them.
func (l *ListInfo) Flatten(lm *ListInfoMap) error {
	if l.HasInclusion {
		for _, inc := range l.Inclusions {
			includedList := (*lm)[inc.source]
			if len(inc.must) == 0 && len(inc.ban) == 0 {
				l.includeAll(includedList)
			} else {
				l.includeSelecting(includedList, inc.must, inc.ban)
			}
		}
	}

	// Ascending label count, so a parent suffix always reaches the trie before
	// the subdomains it is meant to swallow.
	sort.Slice(l.DomainTypeList, func(i, j int) bool {
		return strings.Count(l.DomainTypeList[i].GetValue(), ".") < strings.Count(l.DomainTypeList[j].GetValue(), ".")
	})

	trie := NewDomainTrie()
	for _, domain := range l.DomainTypeList {
		success, err := trie.Insert(domain.GetValue())
		if err != nil {
			return err
		}
		if success {
			l.DomainTypeUniqueList = append(l.DomainTypeUniqueList, domain)
		}
	}

	return nil
}

// includeAll implements `include:filename`: take every rule of the included list.
func (l *ListInfo) includeAll(includedList *ListInfo) {
	l.FullTypeList = append(l.FullTypeList, includedList.FullTypeList...)
	l.DomainTypeList = append(l.DomainTypeList, includedList.DomainTypeList...)
	l.KeywordTypeList = append(l.KeywordTypeList, includedList.KeywordTypeList...)
	l.RegexpTypeList = append(l.RegexpTypeList, includedList.RegexpTypeList...)
	l.AttributeRuleUniqueList = append(l.AttributeRuleUniqueList, includedList.AttributeRuleUniqueList...)
	for attr, domainList := range includedList.AttributeRuleListMap {
		l.AttributeRuleListMap[attr] = append(l.AttributeRuleListMap[attr], domainList...)
	}
}

// includeSelecting implements an include line that carries attributes, covering
// both `include:x @cn` and the upstream `include:x @-!cn` form. A rule qualifies
// when it carries none of the banned attributes and all of the wanted ones,
// which is what upstream v2fly's isMatchAttrFilters does with MustAttrs.
// Notice: a rule can be taken by more than one include line, in which case it
// lands in the lists below once per line. That is harmless here, because every
// list this program publishes goes through `sort --ignore-case -u` in build.yml
// before it is written to the `domains` branch.
func (l *ListInfo) includeSelecting(includedList *ListInfo, must []attribute, ban []attribute) {
	// Rules carrying no attribute at all qualify only when no wanted attribute
	// was given: `include:x @-!cn` takes them, `include:x @cn` does not.
	if len(must) == 0 {
		l.FullTypeList = append(l.FullTypeList, includedList.FullTypeList...)
		l.DomainTypeList = append(l.DomainTypeList, includedList.DomainTypeList...)
		l.KeywordTypeList = append(l.KeywordTypeList, includedList.KeywordTypeList...)
		l.RegexpTypeList = append(l.RegexpTypeList, includedList.RegexpTypeList...)
	}

	for _, rule := range includedList.AttributeRuleUniqueList {
		if hasAnyAttribute(rule, ban) {
			continue
		}
		if !hasAllAttributes(rule, must) {
			continue
		}
		l.AttributeRuleUniqueList = append(l.AttributeRuleUniqueList, rule)
		attrsString := attributesKey(rule)
		l.AttributeRuleListMap[attrsString] = append(l.AttributeRuleListMap[attrsString], rule)
	}
}

// hasAttribute reports whether the rule carries the given attribute, which is
// written in the key form, eg. "@!cn".
func hasAttribute(rule *router.Domain, attrWanted attribute) bool {
	for _, attr := range rule.GetAttribute() {
		if attribute("@"+attr.GetKey()) == attrWanted {
			return true
		}
	}
	return false
}

// hasAnyAttribute reports whether the rule carries at least one of the given
// attributes.
func hasAnyAttribute(rule *router.Domain, attrWanted []attribute) bool {
	for _, attr := range attrWanted {
		if hasAttribute(rule, attr) {
			return true
		}
	}
	return false
}

// hasAllAttributes reports whether the rule carries every one of the given
// attributes. An empty list is satisfied by every rule.
func hasAllAttributes(rule *router.Domain, attrWanted []attribute) bool {
	for _, attr := range attrWanted {
		if !hasAttribute(rule, attr) {
			return false
		}
	}
	return true
}

// attributesKey joins the attributes of a rule into the key used by
// AttributeRuleListMap, eg. "@cn@ads".
func attributesKey(rule *router.Domain) attribute {
	var attrsString attribute
	for _, attr := range rule.GetAttribute() {
		attrsString += attribute("@" + attr.GetKey())
	}
	return attrsString
}

// ToGeoSite converts every ListInfo into a router.GeoSite structure.
// It also excludes rules with certain attributes in certain files that
// user specified in command line when runing the program.
func (l *ListInfo) ToGeoSite(excludeAttrs map[fileName]map[attribute]bool) {
	geosite := new(router.GeoSite)
	geosite.CountryCode = string(l.Name)
	geosite.Domain = append(geosite.Domain, l.FullTypeList...)
	geosite.Domain = append(geosite.Domain, l.DomainTypeUniqueList...)
	geosite.Domain = append(geosite.Domain, l.RegexpTypeList...)

	for _, keywordRule := range l.KeywordTypeList {
		if len(strings.TrimSpace(keywordRule.GetValue())) > 0 {
			geosite.Domain = append(geosite.Domain, keywordRule)
		}
	}

	if excludeAttrs != nil && excludeAttrs[l.Name] != nil {
		excludeAttrsMap := excludeAttrs[l.Name]
		for _, domain := range l.AttributeRuleUniqueList {
			ifKeep := true
			for _, attr := range domain.GetAttribute() {
				if excludeAttrsMap[attribute(attr.GetKey())] {
					ifKeep = false
					break
				}
			}
			if ifKeep {
				geosite.Domain = append(geosite.Domain, domain)
			}
		}
	} else {
		geosite.Domain = append(geosite.Domain, l.AttributeRuleUniqueList...)
	}
	l.GeoSite = geosite
}

// ToPlainText convert router.GeoSite structure to plaintext format.
func (l *ListInfo) ToPlainText() []byte {
	plaintextBytes := make([]byte, 0, 1024*512)

	for _, rule := range l.GeoSite.Domain {
		ruleVal := strings.TrimSpace(rule.GetValue())
		if len(ruleVal) == 0 {
			continue
		}

		var ruleString string
		switch rule.Type {
		case router.Domain_Full:
			ruleString = "full:" + ruleVal
		case router.Domain_RootDomain:
			ruleString = "domain:" + ruleVal
		case router.Domain_Plain:
			ruleString = "keyword:" + ruleVal
		case router.Domain_Regex:
			ruleString = "regexp:" + ruleVal
		}

		if len(rule.Attribute) > 0 {
			ruleString += ":"
			for _, attr := range rule.Attribute {
				ruleString += "@" + attr.GetKey() + ","
			}
			ruleString = strings.TrimRight(ruleString, ",")
		}
		// Output format is: type:domain.tld:@attr1,@attr2
		plaintextBytes = append(plaintextBytes, []byte(ruleString+"\n")...)
	}

	return plaintextBytes
}
