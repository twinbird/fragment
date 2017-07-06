package main

import (
	"bytes"
	"fmt"
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

type setCommandParam struct {
	key     []byte
	value   []byte
	flags   int
	exptime int
}

func setCommand(t *testing.T, con net.Conn, param *setCommandParam) {
	outBuf := new(bytes.Buffer)

	// write command
	fmt.Fprintf(outBuf, "set %s %d %d %d\r\n%s\r\n",
		param.key, param.flags, param.exptime, len(param.value), param.value)

	t.Log(outBuf)
	_, err := con.Write(outBuf.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	// read response
	t.Log("read start")
	recvBuf := make([]byte, RECV_BUF_SIZE)
	rlen, err := con.Read(recvBuf)
	if err != nil {
		t.Fatal(err)
	}
	recvBuf = recvBuf[:rlen]
	t.Logf("recv:%s", string(recvBuf))

	// check read response
	expectVer := []byte("STORED\r\n")
	if bytes.Compare(expectVer, recvBuf) != 0 {
		t.Errorf("set command response error. Expect:%x, Actual:%x\n",
			expectVer, recvBuf)
	}
}

func TestSetAndGet(t *testing.T) {
	con := makeConnection(t)
	defer con.Close()
	con.SetReadDeadline(time.Now().Add(10 * time.Second))
	con.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// set
	setprm := &setCommandParam{
		key:     []byte("name"),
		value:   []byte("twinbird"),
		flags:   12345,
		exptime: 0,
	}
	setCommand(t, con, setprm)

	// get
	t.Log("send:get name")
	_, err := con.Write([]byte("get name\r\n"))
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

	expectGetVer := []byte("VALUE name 12345 8\r\ntwinbird\r\nEND\r\n")
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
	_, err := con.Write([]byte("version\r\n"))
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

	expectVer := []byte("0.0.1\r\n")

	if bytes.Equal(expectVer, recvBuf) == false {
		t.Errorf("Version command error. Expect:%x, Actual:%x\n",
			expectVer, recvBuf)
	}
}
