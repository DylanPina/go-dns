package dns

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type DNSHeader struct {
	ID uint16 // Packet Identifier

	// Flags is a 16-bit field that includes QR, OPCODE, AA, TC, RD, RA, Z, and RCODE
	/*

		| QR | Opcode (4 bits) | AA | TC | RD | RA | Z (3 bits) | RCODE (4 bits) |

		|----|------------------|----|----|----|----|------------|----------------|

		 15   14-11             10   9    8    7    6-4          3-0
	*/
	FLAGS uint16

	QDCOUNT uint16 // Question Count (QDCOUNT)
	ANCOUNT uint16 // Answer Record Count (ANCOUNT)
	NSCOUNT uint16 // Authority Record Count (NSCOUNT)
	ARCOUNT uint16 // Additional Record Count (ARCOUNT)
}

type DNSQuestion struct {
	Name  string // The domain name being queried
	Type  uint16 // Type of the query (A, AAAA, etc.)
	Class uint16 // Class of the query (IN, CH, HS, etc.)
}

func (h *DNSHeader) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	fields := []any{h.ID, h.FLAGS, h.QDCOUNT, h.ANCOUNT, h.NSCOUNT, h.ARCOUNT}

	for _, field := range fields {
		if err := binary.Write(buf, binary.BigEndian, field); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func DecodeDNSHeader(data []byte) (*DNSHeader, error) {
	if len(data) < 12 {
		return nil, errors.New("malformed DNS header")
	}

	return &DNSHeader{
		ID:      binary.BigEndian.Uint16(data[0:2]),
		FLAGS:   binary.BigEndian.Uint16(data[2:4]),
		QDCOUNT: binary.BigEndian.Uint16(data[4:6]),
		ANCOUNT: binary.BigEndian.Uint16(data[6:8]),
		NSCOUNT: binary.BigEndian.Uint16(data[8:10]),
		ARCOUNT: binary.BigEndian.Uint16(data[10:12]),
	}, nil
}

func DecodeDNSQuestion(data []byte) (*DNSQuestion, error) {
	if len(data) < 4 {
		return nil, errors.New("malformed DNS question")
	}

	name, n := decodeDomainName(data)
	if n < 0 {
		return nil, errors.New("malformed domain name")
	}

	return &DNSQuestion{
		Name:  name,
		Type:  binary.BigEndian.Uint16(data[n : n+2]),
		Class: binary.BigEndian.Uint16(data[n+2 : n+4]),
	}, nil
}

func (q *DNSQuestion) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Encode the domain name
	nameParts := bytes.Split([]byte(q.Name), []byte("."))
	for _, part := range nameParts {
		if err := buf.WriteByte(byte(len(part))); err != nil {
			return nil, err
		}
		if _, err := buf.Write(part); err != nil {
			return nil, err
		}
	}

	// Null byte to terminate the domain name
	if err := buf.WriteByte(0); err != nil {
		return nil, err
	}
	// Encode the type
	if err := binary.Write(buf, binary.BigEndian, q.Type); err != nil {
		return nil, err
	}
	// Encode the class
	if err := binary.Write(buf, binary.BigEndian, q.Class); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func decodeDomainName(data []byte) (string, int) {
	var name string
	var n int

	for i := 0; i < len(data); {
		length := int(data[i])
		if length == 0 {
			n++
			break
		}

		if length&0xC0 == 0xC0 {
			offset := int(binary.BigEndian.Uint16(data[i:])) & 0x3FFF
			decode, _ := decodeDomainName(data[offset:])
			name += decode
			n += 2
			break
		}

		if i+length >= len(data) {
			return "", -1
		}

		name += string(data[i+1:i+length+1]) + "."
		i += length + 1
		n += length + 1
	}

	return name[:len(name)-1], n
}
