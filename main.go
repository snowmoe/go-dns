package main

import (
	"fmt"
	"net"

	"github.com/snowmoe/go-dns/internal/dns"
)

// dns packet parsing yay

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

	msg := c.DecodeMessage()

	fmt.Printf(`
   id: %d
flags: %+v
`,
		msg.Header.ID,
		msg.Header.Flags,
	)

	for _, q := range msg.Questions {
		fmt.Printf(`
 name: %s
 type: %d
class: %d
`,
			q.Name,
			q.Type,
			q.Class,
		)
	}

	for _, r := range msg.Resources {

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
}
