package sha1

import (
	"crypto/sha1"
	"reflect"
	"testing"
)

func FuzzHash(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		got := Hash(b)
		expected := sha1.Sum(b)
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})
}
