package main

import (
	"fmt"
	"net"
	"snowmoe/dns/internal/dns"
)

// if you're reading this then stop reading the code
// this is below a PoC atp, it's a cobbled together.. thing
// i'm sick of bytes, i see bytes in my dreams, i breath in bytes
// i eat bytes, i am bytes, bytes[bytes[bytes[bytes]]]

// i thought this would be cool, I thought it would be a nice little PoC
// to understand dns at the packet level

// dns packets = BULLSHIT = BULLSHIT PARSING CODE FMLFMLFMLFML

func main() {
	c := dns.Cursor{
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

	c.ParseHeader()

	q := c.ReadQuestion()
	r := c.ReadResource()

	fmt.Printf(`
 name: %s
 type: %d
class: %d
`,
		q.Name,
		q.Type,
		q.Class,
	)

	fmt.Printf(`
name: %s
 len: %d
data: %d

`,
		r.Name,
		r.Length,
		r.Data,
	)
}
