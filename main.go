package main

import (
	"fmt"
	"net"

	"github.com/snowmoe/go-dns/internal/dns"
)

// dns packet parsing yay

func main() {
	c := dns.Cursor{}

	header := []byte{
		0xAB, 0xCD,
		0x01, 0x00,
		0x00, 0x01,
		0x00, 0x00,
		0x00, 0x00,
		0x00, 0x00,
	}

	// A google.com
	question := []byte{
		6, 'g', 'o', 'o', 'g', 'l', 'e',
		3, 'c', 'o', 'm',
		0,
		0x00, 0x01,
		0x00, 0x01,
	}

	packet := append(header, question...)

	conn, err := net.Dial("udp", "1.1.1.1:53")
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

	fmt.Printf("\nquestion section [%d]\n", len(msg.Questions))
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

	fmt.Printf("\nresource section [%d]\n", len(msg.Resources))
	for _, r := range msg.Resources {

		fmt.Printf(`
	name: %s
	 ttl: %d
	 len: %d
	data: %d
`,
			r.Name,
			r.TTL,
			r.Length,
			r.Data,
		)
	}
}
