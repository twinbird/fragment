package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var db inmemoryDB

func main() {
	listener, err := handleStart(11211)
	if err != nil {
		log.Fatal(err)
	}

	sig_ch := make(chan os.Signal, 1)
	exit_ch := make(chan int)
	signal.Notify(
		sig_ch,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	waitInterruptLoop(sig_ch, exit_ch)

	_ = <-exit_ch
	handleStop(listener)
}

func waitInterruptLoop(sig_ch chan os.Signal, exit_ch chan int) {
	for {
		s := <-sig_ch
		switch s {
		case os.Interrupt:
			exit_ch <- 0
		case syscall.SIGHUP:
			exit_ch <- 0
		case syscall.SIGINT:
			exit_ch <- 0
		case syscall.SIGTERM:
			exit_ch <- 0
		case syscall.SIGQUIT:
			exit_ch <- 0
		}
	}
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
	buf := make([]byte, 1024)

	for {
		mlen, err := con.Read(buf)
		if mlen == 0 {
			log.Println("connection closed")
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
		command, err := parseCommand(buf, mlen)
		if err != nil {
			log.Println(err)
			break
		}

		err = handleCommand(con, command)
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func handleCommand(con net.Conn, command *Command) error {
	switch string(command.command) {
	case "set":
		handleSet(con, command)
	case "add":
		handleAdd(con, command)
	case "replace":
		handleReplace(con, command)
	case "get":
		handleGet(con, command)
	case "delete":
		handleDelete(con, command)
	case "version":
		handleVersion(con, command)
	}
	return nil
}
