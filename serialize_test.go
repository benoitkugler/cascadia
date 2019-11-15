package cascadia

import (
	"reflect"
	"testing"
)

func assembleTests() []string {
	var testSer []string
	for _, test := range selectorTests {
		testSer = append(testSer, test.selector)
	}
	for _, test := range testsPseudo {
		testSer = append(testSer, test.selector)
	}
	for _, test := range validSelectors {
		testSer = append(testSer, test.Selector)
	}
	return testSer
}

func TestSerialize(t *testing.T) {
	for _, selector := range assembleTests() {
		s, err := ParseGroupWithPseudoElements(selector)
		if err != nil {
			t.Fatalf("error compiling %q: %s", selector, err)
		}
		serialized := s.String()
		s2, err := ParseGroupWithPseudoElements(serialized)
		if err != nil {
			t.Fatalf("error compiling %q: %s %#v (original : %s)", serialized, err, s, selector)
		}
		if !reflect.DeepEqual(s, s2) {
			t.Fatalf("can't retrieve selector from serialized : %s (original : %s, sel : %#v)", serialized, selector, s)
		}
	}
}
