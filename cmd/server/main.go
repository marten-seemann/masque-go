package main

import (
	"flag"
	"log"

	"github.com/marten-seemann/masque-go"
	"github.com/marten-seemann/masque-go/internal/testdata"
)

func main() {
	addr := flag.String("addr", "", "MASQUE server (IP:port)")
	flag.Parse()

	if len(*addr) == 0 {
		log.Fatal("missing MASQUE server")
	}

	server := masque.NewServer(*addr, testdata.GetTLSConfig())
	server.Serve()
}
