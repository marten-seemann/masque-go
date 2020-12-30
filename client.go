package masque

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync/atomic"

	"github.com/lucas-clemente/quic-go/http3"
)

type Client struct {
	server *net.UDPAddr
	rt     http3.RoundTripper

	lastFlowID uint32 // use as an atomic
}

func NewClient(tlsConf *tls.Config, masqueServer *net.UDPAddr) *Client {
	return &Client{
		server: masqueServer,
		rt: http3.RoundTripper{
			TLSClientConfig: tlsConf,
			EnableDatagrams: true,
		},
	}
}

// Connect establishes a PacketConn to addr,
// proxied via MASQUE server using the MASQUE protocol.
func (c *Client) Connect(addr *net.UDPAddr) (net.PacketConn, error) {
	url, err := url.Parse("masque://" + c.server.String() + "/")
	if err != nil {
		return nil, err
	}
	fmt.Println(url)
	req := &http.Request{
		URL:    url,
		Method: "CONNECT-UDP",
		Header: http.Header{},
	}
	flowID := atomic.AddUint32(&c.lastFlowID, 2)
	req.Header.Add(flowIDHeader, fmt.Sprintf("%d", flowID))
	rsp, err := c.rt.RoundTripOpt(req, http3.RoundTripOpt{SkipSchemeCheck: true})
	if err != nil {
		return nil, err
	}
	fmt.Printf("Received response: %#v\n", rsp)
	return nil, nil
}
