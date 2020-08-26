package tc

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAction(t *testing.T) {
	tests := map[string]struct {
		val  Action
		err1 error
		err2 error
	}{
		"empty":               {err1: fmt.Errorf("kind is missing")},
		"unknown Kind":        {val: Action{Kind: "test", Index: 123}, err1: fmt.Errorf("unknown kind 'test'")},
		"bpf Without Options": {val: Action{Kind: "bpf"}, err1: ErrNoArg},
		"simple Bpf":          {val: Action{Kind: "bpf", Bpf: &ActBpf{FD: uint32Ptr(12), Name: stringPtr("simpleTest"), Parms: &ActBpfParms{Action: 2, Index: 4}}}},
		"connmark":            {val: Action{Kind: "connmark", ConnMark: &Connmark{Parms: &ConnmarkParam{Index: 42, Action: 1}}}},
		"csum":                {val: Action{Kind: "csum", CSum: &Csum{Parms: &CsumParms{Index: 1, Capab: 2}}}},
		"defact":              {val: Action{Kind: "defact", Defact: &Defact{Parms: &DefactParms{Index: 42, Action: 1}}}},
		"ife":                 {val: Action{Kind: "ife", Ife: &Ife{Parms: &IfeParms{Index: 42, Action: 1}}}},
		"ipt":                 {val: Action{Kind: "ipt", Ipt: &Ipt{Table: stringPtr("testTable"), Hook: uint32Ptr(42), Index: uint32Ptr(1984)}}},
		"mirred":              {val: Action{Kind: "mirred", Mirred: &Mirred{Parms: &MirredParam{Index: 42, Action: 1}}}},
		"mirred+stats":        {val: Action{Kind: "mirred", Mirred: &Mirred{Parms: &MirredParam{Index: 42, Action: 1}}, Stats: &GenStats{Basic: &GenBasic{Bytes: 8, Packets: 1}, RateEst: &GenRateEst{BytePerSecond: 42, PacketPerSecond: 3}, Queue: &GenQueue{QueueLen: 5, Backlog: 6, Drops: 1, Requeues: 3, Overlimits: 1}, RateEst64: &GenRateEst64{BytePerSecond: 12, PacketPerSecond: 1}, BasicHw: &GenBasic{Bytes: 42, Packets: 3}}}},
		"nat":                 {val: Action{Kind: "nat", Nat: &Nat{Parms: &NatParms{Index: 42, Action: 1}}}},
		"police":              {val: Action{Kind: "police", Police: &Police{AvRate: 1337, Result: 42}}},
		"sample":              {val: Action{Kind: "sample", Sample: &Sample{Parms: &SampleParms{Index: 42, Action: 1}}}},
		"vlan":                {val: Action{Kind: "vlan", VLan: &VLan{Parms: &VLanParms{Index: 42, Action: 1}}}},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			data, err1 := marshalActions([]*Action{&testcase.val})
			if err1 != nil {
				if !errors.Is(testcase.err1, err1) {
					return
				}
				t.Fatalf("Unexpected error: %v", err1)
			}
			val := []*Action{}
			err2 := unmarshalActions(data, &val)
			if err2 != nil {
				if !errors.Is(testcase.err2, err2) {
					return
				}
				t.Fatalf("Unexpected error: %v", err2)

			}
			if diff := cmp.Diff(val, []*Action{&testcase.val}); diff != "" {
				t.Fatalf("Action missmatch (-want +got):\n%s", diff)
			}
		})
	}
	t.Run("nil", func(t *testing.T) {
		_, err := marshalAction(nil)
		if !errors.Is(err, ErrNoArg) {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
