package tc

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExtractFqCodelXStats(t *testing.T) {
	tests := map[string]struct {
		data     []byte
		expected *FqCodelXStats
		err      error
	}{
		"Qdisc": {data: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expected: &FqCodelXStats{Type: 0, Qd: &FqCodelQdStats{}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			stats := &FqCodelXStats{}
			if err := extractFqCodelXStats(testcase.data, stats); err != nil {
				if testcase.err != nil && testcase.err.Error() == err.Error() {
					// we received the expected error. everything is fine
					return
				}
				t.Fatalf("received error '%v', but expected '%v'", err, testcase.err)
			}
			if diff := cmp.Diff(stats, testcase.expected); diff != "" {
				t.Fatalf("TestExtractFqCodelXStats missmatch (-want +got):\n%s", diff)
			}
		})
	}
}
