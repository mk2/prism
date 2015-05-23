package prism

import (
	"testing"
)

func TestB2s(t *testing.T) {

	s := "testing"
	sb := []byte(s)

	cs := b2s(sb)

	if s != cs {
		t.Errorf("Unmatch: %s : %s", s, cs)
	}

}
