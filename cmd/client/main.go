package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"

	"github.com/marten-seemann/masque-go"
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

	cl := masque.NewClient(&tls.Config{InsecureSkipVerify: true}, s)
	_, err = cl.Connect(remote)
	log.Fatal(err)
}
