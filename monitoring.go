package main

import (
	"io"
	"net"
)

type discretInput struct {
	slaveId int
	address int
}

var monitorDiscreteInputs map[discretInput]bool
var transactions map[int]bool

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
