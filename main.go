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

type ByteReader struct {
	Pos    int
	Packet []byte
	Labels map[int]string
}

func main() {
	r := ByteReader{
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

	r.Packet = resp

	r.parseHeader()

	fmt.Println("name:", r.parseQName())
	// fmt.Println(r.Pos)

	fmt.Println("qtype:", r.parseQType())

	fmt.Println("qclass:", r.parseQClass())

	fmt.Println("aname:", r.parseName())

	r.skipBytes(10)

	r.printPos()

	fmt.Println("ip:", r.parseRData())

	r.printPos()
}

func (r *ByteReader) parseHeader() {
	// tbc
	r.Pos = 12
}

func (r *ByteReader) parseQType() uint16 {
	qType := binary.BigEndian.Uint16(r.Packet[r.Pos : r.Pos+2])
	r.Pos += 2

	return qType
}

func (r *ByteReader) parseQClass() uint16 {
	qClass := binary.BigEndian.Uint16(r.Packet[r.Pos : r.Pos+2])
	r.Pos += 2

	return qClass
}

func (r *ByteReader) parseQName() string {
	var (
		parts []string
		start = r.Pos
	)

	for {
		len := int(r.Packet[r.Pos])

		if len == 0 {
			r.Pos += 1
			break
		}

		parts = append(parts, string(r.Packet[r.Pos+1:r.Pos+len+1]))
		r.Pos = r.Pos + len + 1
	}

	label := strings.Join(parts, ".")
	r.Labels[start] = label

	return label
}

func (r *ByteReader) parseName() string {
	b := r.Packet[r.Pos]

	// top two bytes = 11 = pointer
	if b&0xC0 == 0xC0 {
		label := r.Labels[int(r.Packet[r.Pos+1])]
		r.Pos += 2
		return label
	} else {
		return ""
	}
}

func (r *ByteReader) parseRData() string {
	ip := r.Packet[r.Pos : r.Pos+4]
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

func (r *ByteReader) skipBytes(n int) {
	r.Pos += n
}

func (r *ByteReader) printPos() {
	fmt.Println("position:", r.Pos)
}
