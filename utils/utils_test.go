package utils

import (
	"testing"
)

func TestMatchSide(t *testing.T) {
	side := ""

	if err := MatchSide(&side, 1); err != nil {
		t.Fatal(err)
	}

	if side != Buy.String() {
		t.Fatal("Side by qty failed.")
	}

	if err := MatchSide(&side, -1); err == nil {
		t.Fatal("Side qty miss-match check failed.")
	}
}
