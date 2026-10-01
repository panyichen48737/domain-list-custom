package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
)

// ListInfoMap is the map of files in data directory and ListInfo
type ListInfoMap map[fileName]*ListInfo

// Marshal processes a file in data directory and generates ListInfo for it.
func (lm *ListInfoMap) Marshal(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	list := NewListInfo()
	listName := fileName(strings.ToUpper(filepath.Base(path)))
	list.Name = listName
	if err := list.ProcessList(file); err != nil {
		return err
	}

	(*lm)[listName] = list
	return nil
}

// FlattenAndGenUniqueDomainList flattens the included lists and
// generates a domain trie for each file in data directory to
// make the items of domain type list unique.
func (lm *ListInfoMap) FlattenAndGenUniqueDomainList() error {
	inclusionLevel := make([]map[fileName]bool, 0, 20)
	okayList := make(map[fileName]bool)
	inclusionLevelAllLength, loopTimes := 0, 0

	for inclusionLevelAllLength < len(*lm) {
		inclusionMap := make(map[fileName]bool)

		if loopTimes == 0 {
			for _, listinfo := range *lm {
				if listinfo.HasInclusion {
					continue
				}
				inclusionMap[listinfo.Name] = true
			}
		} else {
			for _, listinfo := range *lm {
				if !listinfo.HasInclusion || okayList[listinfo.Name] {
					continue
				}

				// A list can only be flattened once every list it includes has
				// been flattened itself.
				ready := true
				for _, inc := range listinfo.Inclusions {
					if !okayList[inc.source] {
						ready = false
						break
					}
				}
				if ready {
					inclusionMap[listinfo.Name] = true
				}
			}
		}

		// No list could be resolved in this pass. Either the includes form a
		// cycle, or some `include:` names a file that is not in the data
		// directory. Without this check the loop below never advances
		// inclusionLevelAllLength and spins forever, appending an empty map
		// on every iteration until the run is killed.
		if len(inclusionMap) == 0 {
			return fmt.Errorf("cannot resolve the include graph, it is either circular or references a missing file:\n%s",
				describeUnresolved(*lm, okayList))
		}

		for filename := range inclusionMap {
			okayList[filename] = true
		}

		inclusionLevel = append(inclusionLevel, inclusionMap)
		inclusionLevelAllLength += len(inclusionMap)
		loopTimes++
	}

	for _, inclusionMap := range inclusionLevel {
		for inclusionFilename := range inclusionMap {
			if err := (*lm)[inclusionFilename].Flatten(lm); err != nil {
				return err
			}
		}
	}

	return nil
}

// describeUnresolved lists the lists that could not be flattened yet together
// with the include targets they are still waiting for. It is only used to make
// the stalled-include error actionable.
func describeUnresolved(lm ListInfoMap, okayList map[fileName]bool) string {
	stuck := make([]string, 0, len(lm))
	for name, listinfo := range lm {
		if okayList[name] {
			continue
		}
		waiting := make([]string, 0, len(listinfo.Inclusions))
		for _, inc := range listinfo.Inclusions {
			if !okayList[inc.source] {
				waiting = append(waiting, string(inc.source))
			}
		}
		sort.Strings(waiting)
		stuck = append(stuck, fmt.Sprintf("  %s waits for: %s", name, strings.Join(waiting, ", ")))
	}
	sort.Strings(stuck)
	return strings.Join(stuck, "\n")
}

// ToProto generates a router.GeoSite for each file in data directory
// and returns a router.GeoSiteList.
// Entries are emitted in name order so the serialized dat file stays
// byte-identical across runs instead of following Go's map iteration order.
func (lm *ListInfoMap) ToProto(excludeAttrs map[fileName]map[attribute]bool) *router.GeoSiteList {
	names := make([]string, 0, len(*lm))
	for name := range *lm {
		names = append(names, string(name))
	}
	sort.Strings(names)

	protoList := new(router.GeoSiteList)
	for _, name := range names {
		listinfo := (*lm)[fileName(name)]
		listinfo.ToGeoSite(excludeAttrs)
		protoList.Entry = append(protoList.Entry, listinfo.GeoSite)
	}
	return protoList
}

// ToPlainText returns a map of exported lists that user wants
// and the contents of them in byte format.
func (lm *ListInfoMap) ToPlainText(exportListsMap []string) (map[string][]byte, error) {
	filePlainTextBytesMap := make(map[string][]byte)
	for _, filename := range exportListsMap {
		if listinfo := (*lm)[fileName(strings.ToUpper(filename))]; listinfo != nil {
			plaintextBytes := listinfo.ToPlainText()
			filePlainTextBytesMap[filename] = plaintextBytes
		} else {
			fmt.Println("Notice: " + filename + ": no such exported list in the directory, skipped.")
		}
	}
	return filePlainTextBytesMap, nil
}
