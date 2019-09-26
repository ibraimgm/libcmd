package libcmd_test

import (
	"errors"
	"testing"

	"github.com/ibraimgm/libcmd"
)

func TestIsParserErr(t *testing.T) {
	tests := []struct {
		err      error
		expected bool
	}{
		{err: nil, expected: false},
		{err: errors.New("my error"), expected: false},
	}

	for i, test := range tests {
		actual := libcmd.IsParserErr(test.err)

		if actual != test.expected {
			t.Errorf("Case %d, expected '%v', but comparison returned '%v'", i, test.expected, actual)
		}
	}
}
