package main

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

// banAttrChild is the included list used by the include-syntax tests below.
// It holds one rule of each kind: no attribute, @cn, and @!cn.
const banAttrChild = `plain.com
attrcn.com @cn
attrbanned.com @!cn
`

// writeListInfoMap writes name->content into a temporary data directory and
// parses every file into a ListInfoMap, the same way main.go does.
func writeListInfoMap(t *testing.T, files map[string]string) ListInfoMap {
	t.Helper()
	dir := t.TempDir()
	lm := make(ListInfoMap)
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("writing %s: %v", name, err)
		}
		if err := lm.Marshal(path); err != nil {
			t.Fatalf("parsing %s: %v", name, err)
		}
	}
	return lm
}

// plainRules returns the exported rules of a list in the plaintext format,
// sorted so tests do not depend on emission order.
func plainRules(t *testing.T, list *ListInfo) []string {
	t.Helper()
	list.ToGeoSite(nil)
	lines := make([]string, 0, len(list.GeoSite.Domain))
	for _, line := range strings.Split(string(list.ToPlainText()), "\n") {
		if line != "" {
			lines = append(lines, line)
		}
	}
	slices.Sort(lines)
	return lines
}

func flatten(t *testing.T, lm ListInfoMap) {
	t.Helper()
	if err := lm.FlattenAndGenUniqueDomainList(); err != nil {
		t.Fatalf("FlattenAndGenUniqueDomainList: %v", err)
	}
}

// `include:child @-!cn` takes everything from child except rules tagged @!cn,
// and unlike the @cn form it also takes the rules carrying no attribute at all.
func TestIncludeExcludingAttribute(t *testing.T) {
	lm := writeListInfoMap(t, map[string]string{
		"child":  banAttrChild,
		"parent": "include:child @-!cn\n",
	})
	flatten(t, lm)

	got := plainRules(t, lm["PARENT"])
	want := []string{"domain:attrcn.com:@cn", "domain:plain.com"}
	if !slices.Equal(got, want) {
		t.Errorf("parent rules = %v, want %v", got, want)
	}
}

// Several ban attributes on one line, as used by upstream v2fly data.
func TestIncludeExcludingMultipleAttributes(t *testing.T) {
	lm := writeListInfoMap(t, map[string]string{
		"child":  banAttrChild,
		"parent": "include:child @-!cn @-ads\n",
	})
	flatten(t, lm)

	got := plainRules(t, lm["PARENT"])
	want := []string{"domain:attrcn.com:@cn", "domain:plain.com"}
	if !slices.Equal(got, want) {
		t.Errorf("parent rules = %v, want %v", got, want)
	}
}

// Several wanted attributes on one line are combined with AND, the way upstream
// v2fly combines MustAttrs, so only the rule carrying both survives.
func TestIncludeWithMultipleAttributesRequiresAll(t *testing.T) {
	lm := writeListInfoMap(t, map[string]string{
		"child": `cnonly.com @cn
cnads.com @cn @ads
adsonly.com @ads
`,
		"parent": "include:child @cn @ads\n",
	})
	flatten(t, lm)

	got := plainRules(t, lm["PARENT"])
	want := []string{"domain:cnads.com:@cn,@ads"}
	if !slices.Equal(got, want) {
		t.Errorf("parent rules = %v, want %v", got, want)
	}
}

// A wanted attribute combined with a banned one keeps the rules that carry the
// wanted one and not the banned one.
func TestIncludeWithWantedAndBannedAttribute(t *testing.T) {
	lm := writeListInfoMap(t, map[string]string{
		"child": `plain.com
attrcn.com @cn
attrcnads.com @cn @ads
`,
		"parent": "include:child @cn @-ads\n",
	})
	flatten(t, lm)

	got := plainRules(t, lm["PARENT"])
	want := []string{"domain:attrcn.com:@cn"}
	if !slices.Equal(got, want) {
		t.Errorf("parent rules = %v, want %v", got, want)
	}
}

// An attribute written without a name is a syntax error, not a silently
// ignored filter: dropping it would turn the line into "take everything".
func TestIncludeWithEmptyAttributeIsRejected(t *testing.T) {
	for _, line := range []string{
		"include:child @",
		"include:child @-",
		"include:child @cn @@",
	} {
		t.Run(line, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "parent")
			if err := os.WriteFile(path, []byte(line+"\n"), 0644); err != nil {
				t.Fatalf("writing parent: %v", err)
			}
			lm := make(ListInfoMap)
			if err := lm.Marshal(path); err == nil {
				t.Fatalf("%q was accepted, want an error", line)
			}
		})
	}
}

// Regression: `include:child @cn` still selects only the rules tagged @cn.
func TestIncludeWithAttribute(t *testing.T) {
	lm := writeListInfoMap(t, map[string]string{
		"child":  banAttrChild,
		"parent": "include:child @cn\n",
	})
	flatten(t, lm)

	got := plainRules(t, lm["PARENT"])
	want := []string{"domain:attrcn.com:@cn"}
	if !slices.Equal(got, want) {
		t.Errorf("parent rules = %v, want %v", got, want)
	}
}

// Regression: plain `include:child` still takes every rule.
func TestIncludeWithoutAttribute(t *testing.T) {
	lm := writeListInfoMap(t, map[string]string{
		"child":  banAttrChild,
		"parent": "include:child\n",
	})
	flatten(t, lm)

	got := plainRules(t, lm["PARENT"])
	want := []string{"domain:attrbanned.com:@!cn", "domain:attrcn.com:@cn", "domain:plain.com"}
	if !slices.Equal(got, want) {
		t.Errorf("parent rules = %v, want %v", got, want)
	}
}

// A subdomain listed before its parent in the file must still be dropped;
// Flatten sorts by label count before feeding the trie.
func TestDomainDedupKeepsParentOnly(t *testing.T) {
	lm := writeListInfoMap(t, map[string]string{
		"dedup": "abc.example.com\nexample.com\n123.example.com\n",
	})
	flatten(t, lm)

	got := plainRules(t, lm["DEDUP"])
	want := []string{"domain:example.com"}
	if !slices.Equal(got, want) {
		t.Errorf("dedup rules = %v, want %v", got, want)
	}
}

// A circular include must be reported, not spun on forever.
func TestIncludeCycleIsReported(t *testing.T) {
	lm := writeListInfoMap(t, map[string]string{
		"aaa": "include:bbb\n",
		"bbb": "include:aaa\n",
	})

	done := make(chan error, 1)
	go func() { done <- lm.FlattenAndGenUniqueDomainList() }()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("circular include returned nil, want an error")
		}
		for _, name := range []string{"AAA", "BBB"} {
			if !strings.Contains(err.Error(), name) {
				t.Errorf("error should name the lists involved, got: %v", err)
			}
		}
	case <-time.After(10 * time.Second):
		t.Fatal("circular include did not return, still looping")
	}
}

// An `include:` naming a file that is not in the data directory must be
// reported instead of stalling every list that depends on it.
func TestMissingIncludeTargetIsReported(t *testing.T) {
	lm := writeListInfoMap(t, map[string]string{
		"aaa": "include:nosuchfile\n",
	})

	err := lm.FlattenAndGenUniqueDomainList()
	if err == nil {
		t.Fatal("missing include target returned nil, want an error")
	}
	if !strings.Contains(err.Error(), "NOSUCHFILE") {
		t.Errorf("error should name the missing target, got: %v", err)
	}
}

// A rule value may itself contain colons, as a regexp matching a URL does. Split
// on every colon left neither branch of parseTypeRule matching and dropped the
// rule from the output without a word.
func TestRuleValueWithColonsIsKept(t *testing.T) {
	lm := writeListInfoMap(t, map[string]string{
		"regex": "regexp:^https?://ads\\.example\\.com\n",
	})
	flatten(t, lm)

	got := plainRules(t, lm["REGEX"])
	want := []string{`regexp:^https?://ads\.example\.com`}
	if !slices.Equal(got, want) {
		t.Errorf("regex rule = %v, want %v", got, want)
	}
}

// A tab between a rule and its attribute separates them just like a space, so
// the attribute is parsed rather than glued onto the end of the value.
func TestRuleSeparatedByTabKeepsAttribute(t *testing.T) {
	lm := writeListInfoMap(t, map[string]string{
		"tabbed": "domain:tabbed.com\t@cn\n",
	})
	flatten(t, lm)

	got := plainRules(t, lm["TABBED"])
	want := []string{"domain:tabbed.com:@cn"}
	if !slices.Equal(got, want) {
		t.Errorf("tab-separated rule = %v, want %v", got, want)
	}
}

// An attribute written without a name on a rule line is rejected, the same way
// the include side rejects it, instead of emitting a rule that no filter selects.
func TestRuleWithEmptyAttributeIsRejected(t *testing.T) {
	for _, line := range []string{
		"domain:emptyattr.com @",
		"domain:emptyattr.com @cn @",
	} {
		t.Run(line, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "list")
			if err := os.WriteFile(path, []byte(line+"\n"), 0644); err != nil {
				t.Fatalf("writing list: %v", err)
			}
			lm := make(ListInfoMap)
			if err := lm.Marshal(path); err == nil {
				t.Fatalf("%q was accepted, want an error", line)
			}
		})
	}
}
