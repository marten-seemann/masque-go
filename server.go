package masque

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/dunglas/httpsfv"

	"github.com/lucas-clemente/quic-go/http3"
)

const flowIDHeader = "Datagram-Flow-Id"

// Server is a MASQUE server.
// It wraps a http3.Server. This allows running MASQUE alongside a regular HTTP server.
type Server struct {
	*http3.Server

	handler http.Handler
}

func (s *Server) setup() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "CONNECT-UDP" {
			s.Server.Handler.ServeHTTP(w, r)
			return
		}
		// TODO: check for the masque scheme
		flowIDItem, err := httpsfv.UnmarshalItem(r.Header[flowIDHeader])
		if err != nil {
			w.WriteHeader(400)
		}
		flowID, ok := flowIDItem.Value.(int64)
		if !ok {
			w.WriteHeader(400)
		}
		fmt.Println("Flow ID:", flowID)
		w.WriteHeader(200)
	})
	s.handler = mux
}

func (s *Server) Serve(conn net.PacketConn) error {
	log.Printf("Listening for incoming connections on %s.", conn.LocalAddr())
	s.setup()
	return s.Server.Serve(conn)
}

func (s *Server) ListenAndServe() error {
	log.Printf("Listening for incoming connections on %s.", s.Addr)
	s.setup()
	return s.Server.ListenAndServe()
}

func (s *Server) Close() error {
	return s.Server.Close()
}
