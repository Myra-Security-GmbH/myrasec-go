package myrasec

const (
	RecordTypeA     = "A"
	RecordTypeAAAA  = "AAAA"
	RecordTypeCNAME = "CNAME"
)

func (rec DNSRecord) CanBeProtected() bool {
	return rec.RecordType == RecordTypeA ||
		rec.RecordType == RecordTypeAAAA ||
		rec.RecordType == RecordTypeCNAME
}
