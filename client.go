package ofxgo

import (
	"errors"
	"io"
	"net/http"
	"time"
)

type Client struct {
	// Request fields to overwrite with the client's values. If nonempty,
	// defaults are used
	SpecVersion string // VERSION in header
	AppId       string // SONRQ>APPID
	AppVer      string // SONRQ>APPVER

	// Don't insert newlines or indentation when marshalling to SGML/XML
	NoIndent bool
}

var defaultClient Client

func (c *Client) OfxVersion() string {
	if len(c.SpecVersion) > 0 {
		return c.SpecVersion
	} else {
		return "203"
	}
}

func (c *Client) Id() String {
	if len(c.AppId) > 0 {
		return String(c.AppId)
	} else {
		return String("OFXGO")
	}
}

func (c *Client) Version() String {
	if len(c.AppVer) > 0 {
		return String(c.AppVer)
	} else {
		return String("0001")
	}
}

func (c *Client) IndentRequests() bool {
	return !c.NoIndent
}

func RawRequest(URL string, r io.Reader) (*http.Response, error) {
	response, err := http.Post(URL, "application/x-ofx", r)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New("OFXQuery request status: " + response.Status)
	}

	return response, nil
}

// Request marshals a Request object into XML, makes an HTTP request against
// it's URL, and then unmarshals the response into a Response object.
//
// Before being marshaled, some of the the Request object's values are
// overwritten, namely those dictated by the Client's configuration (Version,
// AppId, AppVer fields), and the client's curren time (DtClient). These are
// updated in place in the supplied Request object so they may later be
// inspected by the caller.
func (c *Client) Request(r *Request) (*Response, error) {
	r.Signon.DtClient = Date(time.Now())

	// Overwrite fields that the client controls
	r.Version = c.OfxVersion()
	r.Signon.AppId = c.Id()
	r.Signon.AppVer = c.Version()
	r.indent = c.IndentRequests()

	b, err := r.Marshal()
	if err != nil {
		return nil, err
	}

	response, err := RawRequest(r.URL, b)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	ofxresp, err := ParseResponse(response.Body)
	if err != nil {
		return nil, err
	}
	return ofxresp, nil
}