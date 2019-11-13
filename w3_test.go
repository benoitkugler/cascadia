package cascadia

import (
	"encoding/json"
	"io/ioutil"
	"testing"
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

func TestValidSelectors(t *testing.T) {
	c, err := ioutil.ReadFile("test_ressources/valid_selectors.json")
	if err != nil {
		t.Fatal(err)
	}
	var tests []validSelector
	if err = json.Unmarshal(c, &tests); err != nil {
		t.Fatal(err)
	}
	for _, test := range tests {
		if test.Xfail {
			continue
		}
		_, err := ParseGroupWithPseudoElements(test.Selector)
		if err != nil {
			t.Fatalf("%s -> unable to parse valid selector : %s : %s", test.Name, test.Selector, err)
		}
	}
}
