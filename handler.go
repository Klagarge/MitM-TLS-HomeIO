package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"mitm-bridge/modbus-lib"
	"os"
)

type Handler struct {
}

// NewHandler creates a new handler with the given home and electrical simulation.
func NewHandler() *Handler {
	return &Handler{}
}

func (m *Handler) HandleCoils(req *modbus.CoilsRequest) (res []bool, err error) {
	return
}

func (m *Handler) HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) (res []bool, err error) {
	// Load the client certificate with the private key
	clientCert, err := tls.LoadX509KeyPair("HomeIoClientTLS.crt", "private_key_client.pem")
	if err != nil {
		fmt.Printf("failed to load client key pair: %v\n", err)
		return
	}

	// Create the Modbus TCP+TLS client instance
	clientConn, err = modbus.NewClient(&modbus.ClientConfiguration{
		URL:           "tcp+tls://" + clientIP + ":5802",
		TLSClientCert: &clientCert,
		TLSRootCAs:    &x509.CertPool{},
	})
	if err != nil {
		fmt.Printf("failed to start modbus TCP+TLS instance: %v\n", err)
		clientConn = nil
		return
	}

	// Open the Modbus client connection
	err = clientConn.Open()
	if err != nil {
		clientConn = nil
		fmt.Printf("failed to open client connection: %v\n", err)
	}

	unitId := req.UnitId
	err = clientConn.SetUnitId(unitId)
	if err != nil {
		fmt.Printf("failed to set Unit Id: %v\n", err)
		os.Exit(2)
	}

	for address := req.Addr; address < req.Addr+req.Quantity; address++ {
		value, err := clientConn.ReadDiscreteInput(address)
		if err != nil {
			fmt.Printf("failed to read Discret Input: %v\n", err)
		}
		res = append(res, value)

	}

	err = clientConn.Close()
	if err != nil {
		fmt.Printf("failed to close connection: %v\n", err)
	}

	return
}

func (m *Handler) HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) (res []uint16, err error) {

	return
}

func (m *Handler) HandleInputRegisters(req *modbus.InputRegistersRequest) (res []uint16, err error) {

	return
}
