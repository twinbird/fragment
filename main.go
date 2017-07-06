package main

import (
	"bytes"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"
)

const (
	MaxValueByteSize    = 1000000
	MaxKeyByteSize      = 256
	MaxIntCharacterSize = 10
	MaxCommandLength    = 7
)

var db inmemoryDB

func main() {
	listener, err := handleStart(11211)
	if err != nil {
		log.Fatal(err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	for _ = range sigCh {
		break
	}

	handleStop(listener)
}

func handleStop(listener *net.TCPListener) error {
	return listener.Close()
}

func handleStart(port int) (*net.TCPListener, error) {
	if err := db.initialize(); err != nil {
		log.Fatal(err)
	}

	portStr := ":" + strconv.Itoa(port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", portStr)
	if err != nil {
		return nil, err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			con, err := listener.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			go handleClient(con)
		}
	}()
	return listener, nil
}

func handleClient(con net.Conn) {
	defer con.Close()
	con.SetReadDeadline(time.Now().Add(1 * time.Minute))
	con.SetWriteDeadline(time.Now().Add(1 * time.Minute))

	// command + key + parameter(32bit integer) * 3 + data + space * 4 + '\r\n' * 2
	inBuf := make([]byte,
		MaxCommandLength+
			MaxKeyByteSize+
			MaxIntCharacterSize*3+
			MaxValueByteSize+
			4+
			2*2)
	buf := make([]byte, 0)
	buffer := bytes.NewBuffer(buf)

	for {
		mlen, err := con.Read(inBuf)
		if mlen == 0 {
			log.Println("connection closed")
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
		_, err = buffer.Write(inBuf[:mlen])
		if err != nil {
			log.Println(err)
			break
		}
		command, err := parseCommand(buffer.Bytes())
		if err != nil {
			log.Println(err)
			break
		}

		if command.parsed == true {
			err = handleCommand(con, command)
			if err != nil {
				log.Println(err)
				break
			}
			buffer.Reset()
		}
	}
}

func handleCommand(con net.Conn, command *Command) error {
	switch {
	case bytes.Equal(command.command, []byte("set")):
		handleSet(con, command)
	case bytes.Equal(command.command, []byte("add")):
		handleAdd(con, command)
	case bytes.Equal(command.command, []byte("replace")):
		handleReplace(con, command)
	case bytes.Equal(command.command, []byte("get")):
		handleGet(con, command)
	case bytes.Equal(command.command, []byte("delete")):
		handleDelete(con, command)
	case bytes.Equal(command.command, []byte("version")):
		handleVersion(con, command)
	}
	return nil
}
