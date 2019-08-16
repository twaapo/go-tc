package tc

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIpt(t *testing.T) {
	tests := map[string]struct {
		val  Ipt
		err1 error
		err2 error
	}{
		"empty":           {err1: fmt.Errorf("Ipt options are missing")},
		"simple":          {val: Ipt{Table: "testTable", Hook: 42, Index: 1984}},
		"invalidArgument": {val: Ipt{Tm: &Tcft{Install: 1}}, err1: ErrNoArgAlter},
	}
	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalIpt(&testcase.val)
			if err1 != nil {
				if testcase.err1 != nil && testcase.err1.Error() == err1.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}

			val := Ipt{}
			err2 := unmarshalIpt(data, &val)
			if err2 != nil {
				if testcase.err2 != nil && testcase.err2.Error() == err2.Error() {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, testcase.val); diff != "" {
				t.Fatalf("Ipt missmatch (want +got):\n%s", diff)
			}
		})
	}
}