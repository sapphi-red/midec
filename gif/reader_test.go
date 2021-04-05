package gif

import (
	"os"
	"testing"
)

const fixtureFolder = "../fixtures/gif/"

func Test_isAnimated(t *testing.T) {
	t.Parallel()

	runIsAnimated := func(filename string) (bool, error) {
		fp, err := os.Open(fixtureFolder + filename)
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
		{"1.gif", true, false},
		{"2.gif", true, false},
		{"3.gif", false, false},
	}

	for _, tc := range testcases {
		t.Run(tc.filename, func(t *testing.T) {
			actualIsAnimated, actualErr := runIsAnimated(tc.filename)
			if tc.expectedIsAnimated != actualIsAnimated {
				t.Errorf("IsAnimated = %t; want %t", actualIsAnimated, tc.expectedIsAnimated)
			}
			if tc.expectedHasError != (actualErr != nil) {
				t.Errorf("HasError = %v; want %t", actualErr, tc.expectedHasError)
			}
		})
	}
}
