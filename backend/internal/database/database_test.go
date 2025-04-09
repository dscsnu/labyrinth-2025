package database

import (
	"strconv"
	"testing"
)

func TestGenRand(t *testing.T) {

	s, err := genRand()

	if err != nil {

		t.Fatal(err)

	}

	integer, err := strconv.Atoi(s)
	if err != nil {

		t.Fatal(err)

	}

	if integer < 100000 || integer > 999999 {

		t.Fatal("Generate number out of bounds")

	}
}
