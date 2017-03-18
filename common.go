package ofxgo

import (
	"github.com/golang/go/src/encoding/xml"
)

// Represents a top-level OFX message set (i.e. BANKMSGSRSV1)
type Message interface {
	Name() string         // The name of the OFX element this set represents
	Valid() (bool, error) // Called before a Message is marshaled and after
	// it's unmarshaled to ensure the request or response is valid
}

type BankAcct struct {
	XMLName  xml.Name // BANKACCTTO or BANKACCTFROM
	BankId   String   `xml:"BANKID"`
	BranchId String   `xml:"BRANCHID,omitempty"` // Unused in USA
	AcctId   String   `xml:"ACCTID"`
	AcctType String   `xml:"ACCTTYPE"`          // One of CHECKING, SAVINGS, MONEYMRKT, CREDITLINE, CD
	AcctKey  String   `xml:"ACCTKEY,omitempty"` // Unused in USA
}

type CCAcct struct {
	XMLName xml.Name // CCACCTTO or CCACCTFROM
	AcctId  String   `xml:"ACCTID"`
	AcctKey String   `xml:"ACCTKEY,omitempty"` // Unused in USA
}
