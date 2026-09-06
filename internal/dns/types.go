package dns

type Cursor struct {
	Pos  int
	Data []byte
}

type Type uint16
type Class uint16

type TXT []string
