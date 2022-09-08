package action

import (
	"fmt"
	"testing"
)

func TestEmail(t *testing.T) {
	testCases := []struct {
		input string
		valid bool
	}{
		{"", false},
		{"te", false},
		{"test", false},
		{"test@gmail.com", true},
		{"test@gmail.com ", false},
		{" test@gmail.com", false},
		{" test@gmail.com ", false},
		{"test@gmail.com\u00A0", false},
		{"test\u00A0@gmail.com", false},
		{"test\t@gmail.com", false},
		{"\ttest@gmail.com", false},
		{"test@", false},
		{"@gmail.com", false},
		{"homer.@.simpsons.com", false},
		{"homer+1@simpsons.com", true},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			err := Email(tc.input)
			if tc.valid && err != nil {
				t.Errorf("expected %q to be valid, got %v", tc.input, err)
			} else if !tc.valid && err == nil {
				t.Errorf("expected %q to be invalid, got %v", tc.input, err)
			}
		})
	}
}
