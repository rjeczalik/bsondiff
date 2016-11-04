package bsondiff_test

import (
	"reflect"
	"testing"

	"github.com/rjeczalik/bsondiff"
)

func TestDiff(t *testing.T) {
	cases := map[string]struct {
		old  map[string]interface{}
		new  map[string]interface{}
		diff map[string]interface{}
	}{
		"simple diff": {
			old: map[string]interface{}{
				"key": "value",
			},
			new: map[string]interface{}{
				"key":     "value",
				"new key": "new value",
			},
			diff: map[string]interface{}{
				"$set": map[string]interface{}{
					"new key": "new value",
				},
			},
		},
	}

	for name, cas := range cases {
		t.Run(name, func(t *testing.T) {
			var diff map[string]interface{}

			if err := bsondiff.Diff(cas.old, cas.new, &diff); err != nil {
				t.Fatalf("Diff()=%s", err)
			}

			if !reflect.DeepEqual(diff, cas.diff) {
				t.Fatalf("got %+v, want %+v", diff, cas.diff)
			}
		})
	}
}
