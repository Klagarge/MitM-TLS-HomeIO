package main

import (
	"crypto/rand"
	"crypto/tls"
	"log"
	"net"
)

func main() {
	var certServ *tls.Certificate
	var err error

	monitorDiscreteInputs = make(map[discretInput]bool)
	transactions = make(map[int]bool)

	addMonitorDiscreteInput(5, 15, false)
	addMonitorDiscreteInput(5, 14, true)
	addMonitorDiscreteInput(5, 13, true)

	certServ, err = GenerateSelfSignedCertificate("Bender")

	if err != nil {
		log.Fatalf("bridge: loadkeys: %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{*certServ}}
	config.Rand = rand.Reader
	service := "0.0.0.0:5803"
	listener, err := tls.Listen("tcp", service, &config)
	if err != nil {
		log.Fatalf("bridge: listen: %s", err)
	}
	log.Print("bridge: listening")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("bridge: accept: %s", err)
			break
		}
		log.Printf("bridge: accepted from %s", conn.RemoteAddr())
		go handleClient(conn)
	}
}

func handleClient(clientConn net.Conn) {
	defer clientConn.Close()

	var clientCert *tls.Certificate
	var err error

	clientCert, err = GenerateSelfSignedCertificate("Bender")
	if err != nil {
		log.Printf("bridge: load client keys: %s", err)
		return
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{*clientCert},
		InsecureSkipVerify: true,
	}

	serverConn, err := tls.Dial("tcp", "192.168.39.110:5802", tlsConfig)
	if err != nil {
		log.Printf("bridge: unable to connect to server: %s", err)
		return
	}
	defer serverConn.Close()

	go monitorClientToServer(clientConn, serverConn)
	monitorServerToClient(serverConn, clientConn)
}
