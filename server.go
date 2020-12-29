package masque

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/lucas-clemente/quic-go/http3"
)

const flowIDHeader = "Datagram-Flow-Id"

type Server struct {
	h3server http3.Server
}

func NewServer(addr string, tlsConf *tls.Config) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "CONNECT-UDP" {
			w.WriteHeader(400)
		}
		// TODO: check for the masque scheme
		flowIDStr := r.Header.Get(flowIDHeader)
		if len(flowIDStr) == 0 {
			w.WriteHeader(400)
		}
		flowID, err := strconv.Atoi(flowIDStr)
		if err != nil {
			w.WriteHeader(400)
		}
		fmt.Println("Flow ID:", flowID)
		w.WriteHeader(200)
	})
	return &Server{
		h3server: http3.Server{
			Server: &http.Server{
				Addr:      addr,
				TLSConfig: tlsConf,
				Handler:   mux,
			},
			EnableDatagrams: true,
		},
	}
}

func (s *Server) Serve() error {
	log.Printf("Listening for incoming connections on %s.", s.h3server.Addr)
	return s.h3server.ListenAndServe()
}
