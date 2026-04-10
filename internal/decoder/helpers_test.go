package decoder

import "testing"

func TestAsUint64(t *testing.T) {
	if asUint64(uint16(10)) != 10 {
		t.Fatal("uint16 conversion failed")
	}
	if asUint64([]byte{0x01, 0x02}) != 258 {
		t.Fatal("bytes conversion failed")
	}
}

func TestAsIP(t *testing.T) {
	if got := asIP([]byte{10, 0, 0, 1}); got != "10.0.0.1" {
		t.Fatalf("unexpected ip %s", got)
	}
}
