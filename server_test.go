package masque

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"github.com/lucas-clemente/quic-go/http3"

	"github.com/marten-seemann/masque-go/internal/testdata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	var (
		conn net.PacketConn
		url  string
		s    *Server
	)

	BeforeEach(func() {
		s = &Server{
			Server: &http3.Server{
				Server: &http.Server{TLSConfig: testdata.GetTLSConfig()},
			},
		}
		addr, err := net.ResolveUDPAddr("udp", "localhost:0")
		Expect(err).ToNot(HaveOccurred())
		conn, err = net.ListenUDP("udp", addr)
		Expect(err).ToNot(HaveOccurred())
		url = fmt.Sprintf("https://localhost:%d", conn.LocalAddr().(*net.UDPAddr).Port)
	})

	AfterEach(func() {
		Expect(conn.Close()).To(Succeed())
	})

	It("creates a new server", func() {
		go s.Serve(conn)
		defer s.Close()

		cl := http3.RoundTripper{TLSClientConfig: &tls.Config{RootCAs: testdata.GetRootCA()}}
		req, err := http.NewRequest(http.MethodGet, url, nil)
		Expect(err).ToNot(HaveOccurred())
		rsp, err := cl.RoundTrip(req)
		Expect(err).ToNot(HaveOccurred())
		Expect(rsp.StatusCode).To(Equal(404))
	})

	It("closes", func() {
		done := make(chan struct{})
		go func() {
			defer GinkgoRecover()
			defer close(done)
			s.Serve(conn)
		}()

		s.Close()
		Eventually(done).Should(BeClosed())
	})

	It("allows using a HTTP server and a MASQUE server at the same time", func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(418)
		})
		s.Handler = mux
		go s.Serve(conn)
		defer s.Close()

		cl := http3.RoundTripper{TLSClientConfig: &tls.Config{RootCAs: testdata.GetRootCA()}}
		req, err := http.NewRequest(http.MethodGet, url, nil)
		Expect(err).ToNot(HaveOccurred())
		rsp, err := cl.RoundTrip(req)
		Expect(err).ToNot(HaveOccurred())
		Expect(rsp.StatusCode).To(Equal(418))
	})
})
