package png

import (
	"os"
	"testing"
)

const testdataFolder = "../testdata/png/"

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
		{"animated.png", true, false},
		{"static.png", false, false},
		{"invalid-header.png", false, true},
		{"invalid-chunk-header1.png", false, true},
		{"invalid-chunk-header2.png", false, true},
		{"invalid-unknown-chunk.png", false, true},
		{"invalid-actl-chunk.png", false, true},
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
