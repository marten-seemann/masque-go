package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"strings"

	"github.com/marten-seemann/masque-go"
	"github.com/marten-seemann/masque-go/internal/testdata"
)

func main() {
	server := flag.String("server", "", "MASQUE server (IP:port)")
	flag.Parse()

	if len(*server) == 0 {
		log.Fatal("missing MASQUE server")
	}

	s, err := net.ResolveUDPAddr("udp", *server)
	if err != nil {
		log.Fatal(err)
	}
	remote, err := net.ResolveUDPAddr("udp", "localhost:6121")
	if err != nil {
		log.Fatal(err)
	}

	tlsConf := &tls.Config{
		RootCAs:    testdata.GetRootCA(),
		ServerName: strings.Split(*server, ":")[0],
	}
	cl := masque.NewClient(tlsConf, s)
	_, err = cl.Connect(remote)
	log.Fatal(err)
}
