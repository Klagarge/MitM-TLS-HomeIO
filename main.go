package main

import (
	"crypto/tls"
	"fmt"
	"mitm-bridge/modbus-lib"
	"os"
	"sync"
	"time"
)

var (
	clientIP   = "192.168.39.110"
	clientConn *modbus.ModbusClient
	connMutex  sync.Mutex
)

func main() {
	// Load the server certificate authority
	serverCertPool, err := modbus.LoadCertPool("CA-OT-Security.crt")
	if err != nil {
		fmt.Printf("failed to load server certificate authority (CA): %v\n", err)
		os.Exit(1)
	}

	// Load the server certificate with the private key
	serverCert, err := tls.LoadX509KeyPair("HomeIoServerTLS.crt", "private_key_server.pem")
	if err != nil {
		fmt.Printf("failed to load server key pair: %v\n", err)
		os.Exit(1)
	}

	handler := NewHandler()

	// Create the Modbus TCP+TLS server instance
	tlsConnServer, err := modbus.NewServer(&modbus.ServerConfiguration{
		URL:           "tcp+tls://0.0.0.0:5803",
		Timeout:       30 * time.Second,
		MaxClients:    10,
		TLSServerCert: &serverCert,
		TLSClientCAs:  serverCertPool,
	}, handler)
	if err != nil {
		fmt.Printf("failed to start modbus TCP+TLS server instance: %v\n", err)
		os.Exit(1)
	}

	err = tlsConnServer.Start()
	if err != nil {
		fmt.Printf("failed to start server: %v\n", err)
		os.Exit(1)
	}

	for {
		time.Sleep(25 * time.Millisecond)
	}
}
