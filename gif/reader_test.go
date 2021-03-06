package gif

import (
	"os"
	"testing"
)

const testdataFolder = "../testdata/gif/"

func Test_isAnimated(t *testing.T) {
	t.Parallel()

	runIsAnimated := func(filename string) (bool, error) {
		fp, err := os.Open(testdataFolder + filename)
		if err != nil {
			panic(err)
		}
		return isAnimated(fp)
	}

	testcases := []struct {
		filename           string
		expectedIsAnimated bool
		expectedHasError   bool
	}{
		{"loop.gif", true, false},
		{"animated.gif", true, false},
		{"static1.gif", false, false},
		{"static2.gif", false, false},
		{"static-plaintextextension.gif", false, false},
		{"invalid-header-length.gif", false, true},
		{"invalid-header-packedfields.gif", false, true},
		{"invalid-block-header.gif", false, true},
		{"invalid-block-unknown.gif", false, true},
		{"invalid-block-unknown2.gif", false, true},
		{"invalid-block-unknown3.gif", false, true},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.filename, func(t *testing.T) {
			t.Parallel()

			actualIsAnimated, actualErr := runIsAnimated(tc.filename)
			if tc.expectedIsAnimated != actualIsAnimated {
				t.Errorf("IsAnimated = %t; want %t", actualIsAnimated, tc.expectedIsAnimated)
			}
			if tc.expectedHasError != (actualErr != nil) {
				t.Errorf("Error = %v; want HasError = %t", actualErr, tc.expectedHasError)
			}
		})
	}
}
