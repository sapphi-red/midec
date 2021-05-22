package isobmff

import (
	"os"
	"testing"
)

const testdataFolder = "../testdata/isobmff/"

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
		{"animated.avif", true, false},
		{"static.avif", false, false},
		{"static.heif", false, false},
		{"movie.mp4", false, false},
		{"invalid-filetypebox1.avif", false, true},
		{"invalid-filetypebox2.avif", false, true},
		{"invalid-filetypebox3.avif", false, true},
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
