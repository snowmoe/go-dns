package dns

import "strings"

type Cursor struct {
	Pos  int
	Data []byte
}

type Type uint16
type Class uint16

type TXT []string

func (c Class) String() string {
	switch c {
	case 1:
		return "IN"
	case 2:
		return "CS"
	case 3:
		return "CH"
	case 4:
		return "HS"
	case 255:
		return "*"
	default:
		return ""
	}
}

func (t Type) String() string {
	switch t {
	case 1:
		return "A"
	case 2:
		return "NS"
	case 5:
		return "CNAME"
	case 12:
		return "PTR"
	case 15:
		return "MX"
	case 16:
		return "TXT"
	default:
		return ""
	}
}

func (t TXT) String() string {
	var parts []string
	for _, r := range t {
		parts = append(parts, "\""+r+"\"")
	}
	return strings.Join(parts, ",")
}
