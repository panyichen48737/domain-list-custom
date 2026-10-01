package main

import (
	"slices"
	"testing"
)

// DomainTrie is fed in ascending label-count order by Flatten, so a parent
// domain always enters the trie before its subdomains. The cases below keep
// that ordering.
func TestDomainTrieInsert(t *testing.T) {
	tests := []struct {
		name     string
		domains  []string
		wantKept []string
	}{
		{
			name:     "subdomains of an inserted parent are redundant",
			domains:  []string{"example.com", "abc.example.com", "123.example.com"},
			wantKept: []string{"example.com"},
		},
		{
			name:     "the same domain is kept once",
			domains:  []string{"example.com", "example.com"},
			wantKept: []string{"example.com"},
		},
		{
			name:     "different top level domains do not collide",
			domains:  []string{"example.com", "example.cn", "example.com.cn"},
			wantKept: []string{"example.cn", "example.com", "example.com.cn"},
		},
		{
			name:     "a longer unrelated branch is kept",
			domains:  []string{"foo.com", "foo.bar.com"},
			wantKept: []string{"foo.bar.com", "foo.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trie := NewDomainTrie()
			kept := make([]string, 0, len(tt.domains))
			for _, domain := range tt.domains {
				inserted, err := trie.Insert(domain)
				if err != nil {
					t.Fatalf("Insert(%q): %v", domain, err)
				}
				if inserted {
					kept = append(kept, domain)
				}
			}
			slices.Sort(kept)
			if !slices.Equal(kept, tt.wantKept) {
				t.Errorf("kept = %v, want %v", kept, tt.wantKept)
			}
		})
	}
}

func TestDomainTrieInsertEmpty(t *testing.T) {
	if _, err := NewDomainTrie().Insert(""); err == nil {
		t.Error("Insert(\"\") returned nil, want an error")
	}
}
