package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/lucas-clemente/quic-go/http3"

	"github.com/marten-seemann/masque-go"
	"github.com/marten-seemann/masque-go/internal/testdata"
)

func main() {
	addr := flag.String("addr", "", "MASQUE server (IP:port)")
	flag.Parse()

	if len(*addr) == 0 {
		log.Fatal("missing MASQUE server")
	}

	server := masque.Server{
		Server: &http3.Server{
			Server: &http.Server{
				Addr:      *addr,
				TLSConfig: testdata.GetTLSConfig(),
			},
		},
	}
	server.ListenAndServe()
}
