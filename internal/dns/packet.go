package dns

type DNSPacket struct {
	Header   DNSHeader
	Question DNSQuestion
	Answer   DNSAnswer
}
