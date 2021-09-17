package cascadia

import (
	"reflect"
	"testing"
)

func TestSerialize(t *testing.T) {
	var testSer []string
	for _, test := range selectorTests {
		testSer = append(testSer, test.selector)
	}
	for _, test := range testsPseudo {
		testSer = append(testSer, test.selector)
	}
	for _, test := range loadValidSelectors(t) {
		if test.Xfail {
			continue
		}
		testSer = append(testSer, test.Selector)
	}

	xfails := map[string]struct{}{
		// we dont correctly escape in Serialize
		`.foo\:bar`:          {},
		`.test\.foo\[5\]bar`: {},
		`#\#foo\:bar`:        {},
		`#test\.foo\[5\]bar`: {},
	}

	for _, test := range testSer {
		s, err := ParseGroupWithPseudoElements(test)
		if err != nil {
			t.Fatalf("error compiling %q: %s", test, err)
		}
		if _, xfail := xfails[test]; xfail {
			t.Logf("Skipping %s", test)
			continue
		}

		serialized := s.String()
		s2, err := ParseGroupWithPseudoElements(serialized)
		if err != nil {
			t.Errorf("error compiling %q: %s %T (original : %s)", serialized, err, s, test)
		}
		if !reflect.DeepEqual(s, s2) {
			t.Errorf("can't retrieve selector from serialized : %s (original : %s, sel : %#v)", serialized, test, s)
		}
	}
}
