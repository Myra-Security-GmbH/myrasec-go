package myrasec

import (
	"testing"
)

func TestCanBeProtected(t *testing.T) {
	record := DNSRecord{
		RecordType: "A",
	}

	if !record.CanBeProtected() {
		t.Errorf("Expected to get [%t] for DNS record with type [%s]", true, record.RecordType)
	}

	record.RecordType = "AAAA"
	if !record.CanBeProtected() {
		t.Errorf("Expected to get [%t] for DNS record with type [%s]", true, record.RecordType)
	}

	record.RecordType = "CNAME"
	if !record.CanBeProtected() {
		t.Errorf("Expected to get [%t] for DNS record with type [%s]", true, record.RecordType)
	}

	record.RecordType = "MX"
	if record.CanBeProtected() {
		t.Errorf("Expected to get [%t] for DNS record with type [%s]", false, record.RecordType)
	}

	record.RecordType = "TXT"
	if record.CanBeProtected() {
		t.Errorf("Expected to get [%t] for DNS record with type [%s]", false, record.RecordType)
	}

	record.RecordType = "SRV"
	if record.CanBeProtected() {
		t.Errorf("Expected to get [%t] for DNS record with type [%s]", false, record.RecordType)
	}

	record.RecordType = "NS"
	if record.CanBeProtected() {
		t.Errorf("Expected to get [%t] for DNS record with type [%s]", false, record.RecordType)
	}

	record.RecordType = "PTR"
	if record.CanBeProtected() {
		t.Errorf("Expected to get [%t] for DNS record with type [%s]", false, record.RecordType)
	}

	record.RecordType = "CAA"
	if record.CanBeProtected() {
		t.Errorf("Expected to get [%t] for DNS record with type [%s]", false, record.RecordType)
	}
}
