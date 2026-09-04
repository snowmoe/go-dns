package dns

type Header struct {
	ID      uint16
	Flags   Flags
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

type Flags struct {
	QR                 bool
	OPCode             uint8
	Authoritative      bool
	Truncated          bool
	RecursionDesired   bool
	RecursionAvailable bool
	ResponseCode       uint8
}

func (c *Cursor) ParseHeader() Header {
	return Header{
		ID:      c.parseID(),
		Flags:   c.parseFlags(),
		QDCount: c.parseQDCount(),
		ANCount: c.parseANCount(),
		NSCount: c.parseNSCount(),
		ARCount: c.parseARCount(),
	}
}

func (c *Cursor) parseID() uint16      { return c.takeUint16() }
func (c *Cursor) parseQDCount() uint16 { return c.takeUint16() }
func (c *Cursor) parseANCount() uint16 { return c.takeUint16() }
func (c *Cursor) parseNSCount() uint16 { return c.takeUint16() }
func (c *Cursor) parseARCount() uint16 { return c.takeUint16() }

func (c *Cursor) parseFlags() Flags {
	f := c.takeUint16()

	qr := bits(f, 15, 1) == 1
	op := bits(f, 11, 4)
	aa := bits(f, 10, 1) == 1
	tc := bits(f, 9, 1) == 1
	rd := bits(f, 8, 1) == 1
	ra := bits(f, 7, 1) == 1
	rc := bits(f, 0, 4)

	return Flags{
		QR:                 qr,
		OPCode:             uint8(op),
		Authoritative:      aa,
		Truncated:          tc,
		RecursionDesired:   rd,
		RecursionAvailable: ra,
		ResponseCode:       uint8(rc),
	}
}
