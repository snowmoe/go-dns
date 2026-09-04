package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

// if you're reading this then stop reading the code
// this is below a PoC atp, it's a cobbled together.. thing
// i'm sick of bytes, i see bytes in my dreams, i breath in bytes
// i eat bytes, i am bytes, bytes[bytes[bytes[bytes]]]

// i thought this would be cool, I thought it would be a nice little PoC
// to understand dns at the packet level

// dns packets = BULLSHIT = BULLSHIT PARSING CODE FMLFMLFMLFML

type Cursor struct {
	Pos    int
	Data   []byte
	Labels map[int]string
}

func main() {
	c := Cursor{
		Pos:    0,
		Labels: make(map[int]string),
	}

	header := []byte{
		0xAB, 0xCD,
		0x01, 0x00,
		0x00, 0x01,
		0x00, 0x00,
		0x00, 0x00,
		0x00, 0x00,
	}

	question := []byte{
		6, 'g', 'o', 'o', 'g', 'l', 'e', // 7
		3, 'c', 'o', 'm', // 11
		0,          // 12
		0x00, 0x01, // 25, 26
		0x00, 0x01, // 27, 28
	}

	packet := append(header, question...)

	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	conn.Write(packet)
	resp := make([]byte, 512)
	_, err = conn.Read(resp)

	c.Data = resp

	c.parseHeader()

	fmt.Println("name:", c.parseQName())
	// fmt.Println(c.Pos)

	fmt.Println("qtype:", c.parseQType())

	fmt.Println("qclass:", c.parseQClass())

	fmt.Println("aname:", c.parseName())

	c.skipBytes(10)

	c.printPos()

	fmt.Println("ip:", c.parseRData())

	c.printPos()
}

func (c *Cursor) parseHeader() {
	// tbc
	c.Pos = 12
}

func (c *Cursor) parseQType() uint16 {
	qType := binary.BigEndian.Uint16(c.Data[c.Pos : c.Pos+2])
	c.Pos += 2

	return qType
}

func (c *Cursor) parseQClass() uint16 {
	qClass := binary.BigEndian.Uint16(c.Data[c.Pos : c.Pos+2])
	c.Pos += 2

	return qClass
}

func (c *Cursor) parseQName() string {
	var (
		parts []string
		start = c.Pos
	)

	for {
		len := int(c.Data[c.Pos])

		if len == 0 {
			c.Pos += 1
			break
		}

		parts = append(parts, string(c.Data[c.Pos+1:c.Pos+len+1]))
		c.Pos = c.Pos + len + 1
	}

	label := strings.Join(parts, ".")
	c.Labels[start] = label

	return label
}

func (c *Cursor) parseName() string {
	b := c.Data[c.Pos]

	// top two bytes = 11 = pointer
	if b&0xC0 == 0xC0 {
		label := c.Labels[int(c.Data[c.Pos+1])]
		c.Pos += 2
		return label
	} else {
		return ""
	}
}

func (c *Cursor) parseRData() string {
	ip := c.Data[c.Pos : c.Pos+4]
	c.Pos += 3
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

func (c *Cursor) skipBytes(n int) {
	c.Pos += n
}

func (c *Cursor) printPos() {
	fmt.Println("position:", c.Pos)
}
