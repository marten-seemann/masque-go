package masque

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/dunglas/httpsfv"

	"github.com/lucas-clemente/quic-go/http3"
)

const flowIDHeader = "Datagram-Flow-Id"

type Server struct {
	h3server http3.Server
}

func NewServer(addr string, tlsConf *tls.Config) *Server {
	return NewServerWithHandler(addr, tlsConf, http.DefaultServeMux)
}

func NewServerWithHandler(addr string, tlsConf *tls.Config, handler http.Handler) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "CONNECT-UDP" {
			handler.ServeHTTP(w, r)
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

func (s *Server) Close() error {
	return s.h3server.Close()
}
