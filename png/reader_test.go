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
		{"1.png", true, false},
		{"2.png", true, false},
		{"3.png", false, false},
	}

	for _, tc := range testcases {
		t.Run(tc.filename, func(t *testing.T) {
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
