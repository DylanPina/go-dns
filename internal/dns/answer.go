package dns

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type DNSAnswer struct {
	Name   string // The domain name being answered
	Type   uint16 // Type of the answer (A, AAAA, etc.)
	Class  uint16 // Class of the answer (IN, CH, HS, etc.)
	TTL    uint32 // Time to Live (TTL) in seconds
	Length uint16 // Length of the data field
	Data   []byte // The data of the answer (IP address, etc.)
}

func DecodeDNSAnswer(data []byte) (*DNSAnswer, error) {
	if len(data) < 12 {
		return nil, errors.New("malformed DNS answer")
	}

	name, n := decodeDomainName(data)
	if n < 0 {
		return nil, errors.New("malformed domain name")
	}

	length := binary.BigEndian.Uint16(data[n+8 : n+10])

	answer := &DNSAnswer{
		Name:   name,
		Type:   binary.BigEndian.Uint16(data[n : n+2]),
		Class:  binary.BigEndian.Uint16(data[n+2 : n+4]),
		TTL:    binary.BigEndian.Uint32(data[n+4 : n+8]),
		Length: length,
		Data:   data[n+10 : n+10+int(length)],
	}

	return answer, nil
}

func (a *DNSAnswer) Encode() ([]byte, error) {
	buf := make([]byte, 0)

	// Encode the domain name
	nameParts := bytes.Split([]byte(a.Name), []byte("."))
	for _, part := range nameParts {
		buf = append(buf, byte(len(part)))
		buf = append(buf, part...)
	}

	// Null byte to terminate the domain name
	buf = append(buf, 0)

	// Encode the type
	buf = append(buf, byte(a.Type>>8), byte(a.Type))

	// Encode the class
	buf = append(buf, byte(a.Class>>8), byte(a.Class))

	// Encode the TTL
	buf = append(buf, byte(a.TTL>>24), byte(a.TTL>>16), byte(a.TTL>>8), byte(a.TTL))

	// Encode the length
	buf = append(buf, byte(a.Length>>8), byte(a.Length))

	// Append the data
	buf = append(buf, a.Data...)

	return buf, nil
}
