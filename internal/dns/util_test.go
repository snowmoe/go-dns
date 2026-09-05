package dns

import "testing"

func TestBits(t *testing.T) {
	cases := []struct {
		name  string
		value uint16
		shift int
		width int
		want  uint16
	}{
		{"bit 15", 0b1000000000000000, 15, 1, 1},
		{"bit 0", 0b0000000000000001, 0, 1, 1},
		{"bits 8-11", 0b0000111100000000, 8, 4, 0b1111},
		{"bits 8-11 leak", 0b1111101011111111, 8, 4, 0b1010},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := bits(c.value, c.shift, c.width)
			if got != c.want {
				t.Errorf("got %b want %b", got, c.want)
			}
		})
	}
}

func TestTakeUint16(t *testing.T) {
	cases := []struct {
		name    string
		bytes   []byte
		pos     int
		want    uint16
		wantPos int
	}{
		{"bytes 0-1", []byte{0xAB, 0xCD, 0x12, 0x34, 0xFF, 0xFF}, 0, 0xABCD, 2},
		{"bytes 2-3", []byte{0xFF, 0xFF, 0x12, 0x34, 0xFF, 0xFF}, 2, 0x1234, 4},
		{"bytes 4-5", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x43, 0x21}, 4, 0x4321, 6},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cur := &Cursor{
				Pos:  c.pos,
				Data: c.bytes,
			}

			got := cur.takeUint16()

			if got != c.want {
				t.Errorf("got %d want %d", got, c.want)
			}

			if cur.Pos != c.wantPos {
				t.Errorf("got pos %d want %d", cur.Pos, c.wantPos)
			}
		})
	}
}

func TestTakeUint16Sequential(t *testing.T) {
	c := &Cursor{
		Pos:  0,
		Data: []byte{0x12, 0x34, 0x56, 0x78, 0xFF, 0xFF},
	}

	first := c.takeUint16()
	if first != 0x1234 {
		t.Errorf("first call: got %X want 1234", first)
	}

	if c.Pos != 2 {
		t.Errorf("first call: got pos %d want 2", c.Pos)
	}

	second := c.takeUint16()
	if second != 0x5678 {
		t.Errorf("second call: got %X want 5678", second)
	}

	if c.Pos != 4 {
		t.Errorf("second call: got pos %d want 4", c.Pos)
	}
}

func TestTakeUint32(t *testing.T) {
	cases := []struct {
		name    string
		bytes   []byte
		pos     int
		want    uint32
		wantPos int
	}{
		{"bytes 0-3", []byte{0xAB, 0xCD, 0x12, 0x34, 0xFF, 0xFF}, 0, 0xABCD1234, 4},
		{"bytes 2-5", []byte{0xFF, 0xFF, 0x12, 0x34, 0x56, 0x78}, 2, 0x12345678, 6},
		{"bytes 1-4", []byte{0xFF, 0x87, 0x65, 0x43, 0x21, 0xFF}, 1, 0x87654321, 5},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cur := &Cursor{
				Pos:  c.pos,
				Data: c.bytes,
			}

			got := cur.takeUint32()

			if got != c.want {
				t.Errorf("got %d want %d", got, c.want)
			}

			if cur.Pos != c.wantPos {
				t.Errorf("got pos %d want %d", cur.Pos, c.wantPos)
			}
		})
	}
}

func TestTakeUint32Sequential(t *testing.T) {
	c := &Cursor{
		Pos:  2,
		Data: []byte{0xFF, 0xFF, 0x12, 0x34, 0x56, 0x78, 0xAB, 0xCD, 0x43, 0x21},
	}

	first := c.takeUint32()
	if first != 0x12345678 {
		t.Errorf("first call: got %X want 12345678", first)
	}

	if c.Pos != 6 {
		t.Errorf("first call: got pos %d want 6", c.Pos)
	}

	second := c.takeUint32()
	if second != 0xABCD4321 {
		t.Errorf("second call: got %X want ABCD4321", second)
	}

	if c.Pos != 10 {
		t.Errorf("second call: got pos %d want 10", c.Pos)
	}
}
