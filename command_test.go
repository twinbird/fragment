package main

import (
	"bytes"
	"net"
	"os"
	"testing"
	"time"
)

const RECV_BUF_SIZE = 1024

func TestMain(m *testing.M) {
	li, err := handleStart(11211)
	if err != nil {
		panic(err)
	}
	code := m.Run()
	err = handleStop(li)
	if err != nil {
		panic(err)
	}
	os.Exit(code)
}

func makeConnection(t *testing.T) net.Conn {
	tcp_addr, err := net.ResolveTCPAddr("tcp", "localhost:11211")
	if err != nil {
		t.Fatal(err)
	}
	con, err := net.DialTCP("tcp", nil, tcp_addr)
	if err != nil {
		t.Fatal(err)
	}
	return con
}

func sendQuit(t *testing.T, con net.Conn) {
	// send quit
	t.Log("send:quit")
	_, err := con.Write([]byte("quit"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetAndGet(t *testing.T) {
	con := makeConnection(t)
	defer con.Close()
	con.SetReadDeadline(time.Now().Add(10 * time.Second))
	con.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// set
	t.Log("set name 12345 0 8")
	_, err := con.Write([]byte("set name 12345 0 8"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("twinbird")
	_, err = con.Write([]byte("twinbird"))
	if err != nil {
		t.Fatal(err)
	}

	recvBuf := make([]byte, RECV_BUF_SIZE)
	rlen, err := con.Read(recvBuf)
	if err != nil {
		t.Fatal(err)
	}
	recvBuf = recvBuf[:rlen]
	t.Logf("recv:%s", string(recvBuf))

	expectVer := []byte("STORED\n")

	if bytes.Compare(expectVer, recvBuf) != 0 {
		t.Errorf("set command response error. Expect:%x, Actual:%x\n",
			expectVer, recvBuf)
	}

	// get
	t.Log("send:get name")
	_, err = con.Write([]byte("get name"))
	if err != nil {
		t.Fatal(err)
	}
	getRecvBuf := make([]byte, RECV_BUF_SIZE)
	grlen, err := con.Read(getRecvBuf)
	if err != nil {
		t.Fatal(err)
	}
	getRecvBuf = getRecvBuf[:grlen]
	t.Logf("recv:%s", string(getRecvBuf))

	expectGetVer := []byte("VALUE name 12345 8\ntwinbird\nEND")
	if bytes.Compare(expectGetVer, getRecvBuf) != 0 {
		t.Errorf("get command error. Expect:%x, Actual:%x\n",
			expectGetVer, getRecvBuf)
	}
}

func TestVersion(t *testing.T) {
	con := makeConnection(t)
	defer con.Close()
	con.SetReadDeadline(time.Now().Add(10 * time.Second))
	con.SetWriteDeadline(time.Now().Add(10 * time.Second))

	t.Log("send:version")
	_, err := con.Write([]byte("version"))
	if err != nil {
		t.Fatal(err)
	}

	recvBuf := make([]byte, RECV_BUF_SIZE)
	rlen, err := con.Read(recvBuf)
	if err != nil {
		t.Fatal(err)
	}
	recvBuf = recvBuf[:rlen]
	t.Logf("recv:%s", string(recvBuf))

	expectVer := []byte("0.0.1\n")

	if bytes.Compare(expectVer, recvBuf) != 0 {
		t.Errorf("Version command error. Expect:%x, Actual:%x\n",
			expectVer, recvBuf)
	}
}
