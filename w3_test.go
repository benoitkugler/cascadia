package cascadia

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"golang.org/x/net/html"
)

type invalidSelector struct {
	Name     string `json:"name,omitempty"`
	Selector string `json:"selector,omitempty"`
}

type validSelector struct {
	invalidSelector
	Expect  []string `json:"expect,omitempty"`
	Exclude []string `json:"exclude,omitempty"`
	Level   int      `json:"level,omitempty"`
	Xfail   bool     `json:"xfail,omitempty"`
}

func TestInvalidSelectors(t *testing.T) {
	c, err := ioutil.ReadFile("test_ressources/invalid_selectors.json")
	if err != nil {
		t.Fatal(err)
	}
	var tests []invalidSelector
	if err = json.Unmarshal(c, &tests); err != nil {
		t.Fatal(err)
	}
	for _, test := range tests {
		_, err := ParseGroupWithPseudoElements(test.Selector)
		if err == nil {
			t.Fatalf("%s -> expected error on invalid selector : %s", test.Name, test.Selector)
		}
	}
}

func parseReference() *html.Node {
	f, err := os.Open("test_ressources/content.xhtml")
	if err != nil {
		log.Fatal(err)
	}
	node, err := html.Parse(f)
	if err != nil {
		log.Fatal(err)
	}
	return node
}

func getId(n *html.Node) string {
	for _, attr := range n.Attr {
		if attr.Key == "id" {
			return attr.Val
		}
	}
	return ""
}

func isEqual(m map[string]bool, l []string) bool {
	if len(m) != len(l) {
		return false
	}
	for _, s := range l {
		if !m[s] {
			return false
		}
	}
	return true
}

func TestValidSelectors(t *testing.T) {
	c, err := ioutil.ReadFile("test_ressources/valid_selectors.json")
	if err != nil {
		t.Fatal(err)
	}
	var tests []validSelector
	if err = json.Unmarshal(c, &tests); err != nil {
		t.Fatal(err)
	}
	doc := parseReference()
	for _, test := range tests {
		if test.Xfail {
			t.Logf("skiped test %s", test.Name)
			continue
		}
		sels, err := ParseGroupWithPseudoElements(test.Selector)
		if err != nil {
			t.Fatalf("%s -> unable to parse valid selector : %s : %s", test.Name, test.Selector, err)
		}
		matchingIds := map[string]bool{}
		for _, sel := range sels {
			if sel.PseudoElement() != "" {
				continue // pseudo element doesn't count as a match in this test since they are not part of the document
			}
			for _, node := range Selector(sel.Match).MatchAll(doc) {
				matchingIds[getId(node)] = true
			}
		}
		if !isEqual(matchingIds, test.Expect) {
			t.Fatalf("%s : expected %v got %v", test.Name, test.Expect, matchingIds)
		}

	}
}
