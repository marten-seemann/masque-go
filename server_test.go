package masque

import (
	"crypto/tls"
	"net/http"

	"github.com/lucas-clemente/quic-go/http3"

	"github.com/marten-seemann/masque-go/internal/testdata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	It("creates a new server", func() {
		s := NewServer("localhost:12345", testdata.GetTLSConfig())
		go s.Serve()

		cl := http3.RoundTripper{TLSClientConfig: &tls.Config{RootCAs: testdata.GetRootCA()}}
		req, err := http.NewRequest(http.MethodGet, "https://localhost:12345/", nil)
		Expect(err).ToNot(HaveOccurred())
		rsp, err := cl.RoundTrip(req)
		Expect(err).ToNot(HaveOccurred())
		Expect(rsp.StatusCode).To(Equal(404))
	})

	It("allows using a HTTP server and a MASQUE server at the same time", func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(418)
		})
		s := NewServerWithHandler("localhost:12346", testdata.GetTLSConfig(), mux)
		go s.Serve()

		cl := http3.RoundTripper{TLSClientConfig: &tls.Config{RootCAs: testdata.GetRootCA()}}
		req, err := http.NewRequest(http.MethodGet, "https://localhost:12346/", nil)
		Expect(err).ToNot(HaveOccurred())
		rsp, err := cl.RoundTrip(req)
		Expect(err).ToNot(HaveOccurred())
		Expect(rsp.StatusCode).To(Equal(418))
	})
})
