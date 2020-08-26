package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBpf(t *testing.T) {
	tests := map[string]struct {
		val  Bpf
		err1 error
		err2 error
	}{
		"simple": {val: Bpf{Ops: bytesPtr([]byte{0x6, 0x0, 0x0, 0x0, 0xff, 0xff, 0xff, 0xff}),
			OpsLen:  uint16Ptr(0x1),
			ClassID: uint32Ptr(0x10001),
			Flags:   uint32Ptr(0x1)}},
		"da obj /tmp/bpf.o sec foo": {val: Bpf{FD: uint32Ptr(8), Name: stringPtr("bpf.o:[foo]"), Flags: uint32Ptr(0x1), FlagsGen: uint32Ptr(0x0)}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalBpf(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := Bpf{}
			err2 := unmarshalBpf(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Bpf missmatch (want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalBpf(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
