package midec_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/sapphi-red/midec"
)

func Test_ReadAdvancer_Advance(t *testing.T) {
	t.Parallel()

	runAdvance := func(byteLen, advanceLen int) error {
		empty := make([]byte, byteLen)

		advancer := midec.NewReadAdvancer(bytes.NewReader(empty))
		return advancer.Advance(uint(advanceLen))
	}

	testcases := []struct {
		byteLen           int
		advanceLen int
		expectedHasError   bool
	}{
		{0, 1, true},
		{1, 1, false},
		{256 * 3 + 1, 256 * 3 + 1, false},
		{256 * 3 + 1, 256 * 3 + 2, true},
	}

	for _, tc := range testcases {
		name := fmt.Sprintf("%d %d", tc.byteLen, tc.advanceLen)
		t.Run(name, func(t *testing.T) {
			actualErr := runAdvance(tc.byteLen, tc.advanceLen)
			if tc.expectedHasError != (actualErr != nil) {
				t.Errorf("Error = %v; want HasError = %t", actualErr, tc.expectedHasError)
			}
		})
	}
}
