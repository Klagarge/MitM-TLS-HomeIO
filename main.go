package main

import (
	"crypto/rand"
	"crypto/tls"
	"io"
	"log"
	"net"
)

var fakeClient = true
var fakeServer = true

type discretInput struct {
	slaveId int
	address int
}

var monitorDiscreteInputs map[discretInput]bool
var transactions map[int]bool

func main() {
	var certServ tls.Certificate
	var err error

	monitorDiscreteInputs = make(map[discretInput]bool)
	transactions = make(map[int]bool)

	addMonitorDiscreteInput(5, 15, false)
	addMonitorDiscreteInput(5, 14, true)
	addMonitorDiscreteInput(5, 13, true)

	if fakeServer {
		certServ, err = tls.LoadX509KeyPair("fakeCertificateServer.crt", "fakeCertificateServer.pem")
	} else {
		certServ, err = tls.LoadX509KeyPair("HomeIoServerTLS.crt", "private_key_server.pem")
	}
	if err != nil {
		log.Fatalf("bridge: loadkeys: %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{certServ}}
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

	var clientCert tls.Certificate
	var err error

	if fakeClient {
		clientCert, err = tls.LoadX509KeyPair("fakeCertificateClient.crt", "fakeCertificateClient.pem")
	} else {
		clientCert, err = tls.LoadX509KeyPair("HomeIoClientTLS.crt", "private_key_client.pem")
	}
	if err != nil {
		log.Printf("bridge: load client keys: %s", err)
		return
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{clientCert},
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

func addMonitorDiscreteInput(slaveId int, address int, valueShouldBe bool) {
	monitorDiscreteInputs[discretInput{slaveId, address}] = valueShouldBe
}

func monitorClientToServer(src net.Conn, dst net.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := src.Read(buf)
		if err != nil {
			if err != io.EOF {
				println("bridge: error reading from client to server: ", err)
			}
			return
		}

		var functionCode = buf[7]

		if functionCode == 2 {
			var transactionId = int(buf[0])<<8 + int(buf[1])
			var unitId = buf[6]
			var startAddress = int(buf[8])<<8 + int(buf[9])
			var key = discretInput{int(unitId), startAddress}
			_, prs := monitorDiscreteInputs[key]
			if prs {
				transactions[transactionId] = monitorDiscreteInputs[key]
				println("Monitoring discrete input: ", key.slaveId, ":", key.address, " to ", monitorDiscreteInputs[key])
			}
		}

		n, err = dst.Write(buf[:n])
		if err != nil {
			println("bridge: error writing to client to server: ", err)
			return
		}
	}
}

func monitorServerToClient(src net.Conn, dst net.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := src.Read(buf)
		if err != nil {
			if err != io.EOF {
				println("bridge: error reading from server to client: ", err)
			}
			return
		}

		transactionId := int(buf[0])<<8 + int(buf[1])
		_, prs := transactions[transactionId]
		if prs {
			//c := buf[8]
			inputStatus := buf[9]
			valueToSet := transactions[transactionId]
			delete(transactions, transactionId)
			if valueToSet != (inputStatus == 1) {
				buf[9] = boolToByte(valueToSet)
				println("Modified packet id: ", transactionId, " to ", valueToSet)
			}

		}

		n, err = dst.Write(buf[:n])
		if err != nil {
			println("bridge: error writing to server to client: ", err)
			return
		}
	}

}

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}
