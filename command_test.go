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
	if bytes.Equal(expectVer, recvBuf) == false {
		t.Errorf("set command response error. Expect:%x, Actual:%x\n",
			expectVer, recvBuf)
	}
}

func addCommand(t *testing.T, con net.Conn, param *setCommandParam, exist bool) {
	buf := new(bytes.Buffer)

	// write command
	fmt.Fprintf(buf, "add %s %d %d %d\r\n%s\r\n",
		param.key, param.flags, param.exptime, len(param.value), param.value)

	t.Log(buf)
	_, err := con.Write(buf.Bytes())
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
	var expectVer []byte
	if exist == true {
		expectVer = []byte("NOT_STORED\r\n")
	} else {
		expectVer = []byte("STORED\r\n")
	}
	if bytes.Equal(expectVer, recvBuf) == false {
		t.Errorf("add command response error. Expect:%x, Actual:%x\n",
			expectVer, recvBuf)
	}
}

func deleteCommand(t *testing.T, con net.Conn, param *setCommandParam, exist bool) {
	buf := new(bytes.Buffer)

	// write command
	fmt.Fprintf(buf, "delete %s\r\n", param.key)

	t.Log(buf)
	_, err := con.Write(buf.Bytes())
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
	var expectVer []byte
	if exist == true {
		expectVer = []byte("DELETED\r\n")
	} else {
		expectVer = []byte("NOT_FOUND\r\n")
	}
	if bytes.Equal(expectVer, recvBuf) == false {
		t.Errorf("delete command response error. Expect:%x, Actual:%x\n", expectVer, recvBuf)
	}
}

func replaceCommand(t *testing.T, con net.Conn, param *setCommandParam, exist bool) {
	buf := new(bytes.Buffer)

	// write command
	fmt.Fprintf(buf, "replace %s %d %d %d\r\n%s\r\n",
		param.key, param.flags, param.exptime, len(param.value), param.value)

	t.Log(buf)
	_, err := con.Write(buf.Bytes())
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
	var expectVer []byte
	if exist == false {
		expectVer = []byte("NOT_STORED\r\n")
	} else {
		expectVer = []byte("STORED\r\n")
	}
	if bytes.Equal(expectVer, recvBuf) == false {
		t.Errorf("replace command response error. Expect:%x, Actual:%x\n",
			expectVer, recvBuf)
	}
}

func notFoundGetCommand(t *testing.T, con net.Conn, param *setCommandParam) {
	buf := new(bytes.Buffer)

	// write command
	fmt.Fprintf(buf, "get %s\r\n", param.key)

	t.Log(buf)
	_, err := con.Write(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	// read response
	getRecvBuf := make([]byte, RECV_BUF_SIZE)
	grlen, err := con.Read(getRecvBuf)
	if err != nil {
		t.Fatal(err)
	}
	getRecvBuf = getRecvBuf[:grlen]
	t.Logf("recv:%s", string(getRecvBuf))

	// check response
	buf.Reset()
	fmt.Fprintf(buf, "NOT_FOUND\r\n")
	if bytes.Equal(buf.Bytes(), getRecvBuf) == false {
		t.Errorf("get command error. Expect:%x, Actual:%x\n", buf, getRecvBuf)
	}
}

func getCommand(t *testing.T, con net.Conn, param *setCommandParam) {
	buf := new(bytes.Buffer)

	// write command
	fmt.Fprintf(buf, "get %s\r\n", param.key)

	t.Log(buf)
	_, err := con.Write(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	// read response
	getRecvBuf := make([]byte, RECV_BUF_SIZE)
	grlen, err := con.Read(getRecvBuf)
	if err != nil {
		t.Fatal(err)
	}
	getRecvBuf = getRecvBuf[:grlen]
	t.Logf("recv:%s", string(getRecvBuf))

	// check response
	buf.Reset()
	if len(param.value) == 0 {
		fmt.Fprintf(buf, "END\r\n")
	} else {
		fmt.Fprintf(buf, "VALUE %s %d %d\r\n%s\r\nEND\r\n", param.key, param.flags, len(param.value), param.value)
	}
	if bytes.Equal(buf.Bytes(), getRecvBuf) == false {
		t.Errorf("get command error. Expect:%x, Actual:%x\n", buf, getRecvBuf)
	}
}

func TestSetAndGet(t *testing.T) {
	con := makeConnection(t)
	defer con.Close()
	con.SetReadDeadline(time.Now().Add(10 * time.Second))
	con.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// set
	setprm := &setCommandParam{
		key:     []byte("SetAndGetKey"),
		value:   []byte("SetAndGetValue"),
		flags:   12345,
		exptime: 0,
	}
	setCommand(t, con, setprm)

	// get
	getCommand(t, con, setprm)
}

func TestSetAndReplaceAndGet(t *testing.T) {
	con := makeConnection(t)
	defer con.Close()
	con.SetReadDeadline(time.Now().Add(10 * time.Second))
	con.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// set
	setprm := &setCommandParam{
		key:     []byte("SetAndReplaceAndGetKey"),
		value:   []byte("SetAndReplaceAndGetValue"),
		flags:   12345,
		exptime: 0,
	}
	setCommand(t, con, setprm)

	// replace
	setprm.value = []byte("replaced")
	replaceCommand(t, con, setprm, true)

	// get
	getCommand(t, con, setprm)
}

func TestNonSetReplace(t *testing.T) {
	con := makeConnection(t)
	defer con.Close()
	con.SetReadDeadline(time.Now().Add(10 * time.Second))
	con.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// replace
	setprm := &setCommandParam{
		key:     []byte("NonSetReplaceKey"),
		value:   []byte("NonSetReplaceValue"),
		flags:   12345,
		exptime: 0,
	}
	replaceCommand(t, con, setprm, false)
}

func TestSetAndOverwriteAndGet(t *testing.T) {
	con := makeConnection(t)
	defer con.Close()
	con.SetReadDeadline(time.Now().Add(10 * time.Second))
	con.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// set
	setprm := &setCommandParam{
		key:     []byte("SetAndOverwriteAndGetKey"),
		value:   []byte("SetAndOverwriteAndGetKeyValue"),
		flags:   12345,
		exptime: 0,
	}
	setCommand(t, con, setprm)

	// overwrite
	setprm.value = []byte("overwrite")
	setCommand(t, con, setprm)

	// get
	getCommand(t, con, setprm)
}

func TestNonSetGet(t *testing.T) {
	con := makeConnection(t)
	defer con.Close()
	con.SetReadDeadline(time.Now().Add(10 * time.Second))
	con.SetWriteDeadline(time.Now().Add(10 * time.Second))

	setprm := &setCommandParam{
		key:     []byte("NonSetGet"),
		value:   []byte(""),
		flags:   0,
		exptime: 0,
	}
	notFoundGetCommand(t, con, setprm)
}

func TestAddAndGet(t *testing.T) {
	con := makeConnection(t)
	defer con.Close()
	con.SetReadDeadline(time.Now().Add(10 * time.Second))
	con.SetWriteDeadline(time.Now().Add(10 * time.Second))

	setprm := &setCommandParam{
		key:     []byte("AddAndGetKey"),
		value:   []byte("AddAndGetValue"),
		flags:   12345,
		exptime: 0,
	}
	addCommand(t, con, setprm, false)
	getCommand(t, con, setprm)
}

func TestAddAndExpiredGet(t *testing.T) {
	con := makeConnection(t)
	defer con.Close()
	con.SetReadDeadline(time.Now().Add(10 * time.Second))
	con.SetWriteDeadline(time.Now().Add(10 * time.Second))

	setprm := &setCommandParam{
		key:     []byte("AddAndExpiredGetKey"),
		value:   []byte("AddAndExpiredGetValue"),
		flags:   12345,
		exptime: 1,
	}
	addCommand(t, con, setprm, false)

	// wait 1 second for expire
	time.Sleep(2 * time.Second)

	notFoundGetCommand(t, con, setprm)
}

func TestSetAndExpiredGet(t *testing.T) {
	con := makeConnection(t)
	defer con.Close()
	con.SetReadDeadline(time.Now().Add(10 * time.Second))
	con.SetWriteDeadline(time.Now().Add(10 * time.Second))

	setprm := &setCommandParam{
		key:     []byte("AddAndExpiredGetKey"),
		value:   []byte("AddAndExpiredGetValue"),
		flags:   12345,
		exptime: 1,
	}
	setCommand(t, con, setprm)

	// wait 1 second for expire
	time.Sleep(2 * time.Second)

	notFoundGetCommand(t, con, setprm)
}

func TestDuplicateAdd(t *testing.T) {
	con := makeConnection(t)
	defer con.Close()
	con.SetReadDeadline(time.Now().Add(10 * time.Second))
	con.SetWriteDeadline(time.Now().Add(10 * time.Second))

	setprm := &setCommandParam{
		key:     []byte("add-key"),
		value:   []byte("add-value"),
		flags:   12345,
		exptime: 0,
	}
	addCommand(t, con, setprm, false)
	addCommand(t, con, setprm, true)
}

func TestSetAndDeleteAndGet(t *testing.T) {
	con := makeConnection(t)
	defer con.Close()
	con.SetReadDeadline(time.Now().Add(10 * time.Second))
	con.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// set
	setprm := &setCommandParam{
		key:     []byte("SetAndDeleteAndGetKey"),
		value:   []byte("SetAndDeleteAndGetValue"),
		flags:   12345,
		exptime: 0,
	}
	setCommand(t, con, setprm)

	// delete
	deleteCommand(t, con, setprm, true)

	// get
	setprm.value = []byte("")
	notFoundGetCommand(t, con, setprm)
}

func TestNonExistDelete(t *testing.T) {
	con := makeConnection(t)
	defer con.Close()
	con.SetReadDeadline(time.Now().Add(10 * time.Second))
	con.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// delete
	setprm := &setCommandParam{
		key:     []byte("NonExistDeleteKey"),
		value:   []byte("NonExistDeleteValue"),
		flags:   12345,
		exptime: 0,
	}
	deleteCommand(t, con, setprm, false)
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

func TestSetDataOverBytesAndGet(t *testing.T) {
	con := makeConnection(t)
	defer con.Close()
	con.SetReadDeadline(time.Now().Add(10 * time.Second))
	con.SetWriteDeadline(time.Now().Add(10 * time.Second))

	param := &setCommandParam{
		key:     []byte("SetDataOverBytesKey"),
		value:   []byte("SetDataOverBytesValue"),
		flags:   12345,
		exptime: 0,
	}

	outBuf := new(bytes.Buffer)

	// write command
	fmt.Fprintf(outBuf, "set %s %d %d %d\r\n%s\r\n",
		param.key, param.flags, param.exptime, len(param.value)-1, param.value)

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
	if bytes.Equal(expectVer, recvBuf) == false {
		t.Errorf("set command response error. Expect:%x, Actual:%x\n",
			expectVer, recvBuf)
	}

	// rewrite to right response data
	param.value = param.value[:len(param.value)-1]

	// get check
	getCommand(t, con, param)
}

func TestAddDataOverBytesAndGet(t *testing.T) {
	con := makeConnection(t)
	defer con.Close()
	con.SetReadDeadline(time.Now().Add(10 * time.Second))
	con.SetWriteDeadline(time.Now().Add(10 * time.Second))

	param := &setCommandParam{
		key:     []byte("AddDataOverBytesKey"),
		value:   []byte("AddDataOverBytesValue"),
		flags:   12345,
		exptime: 0,
	}

	outBuf := new(bytes.Buffer)

	// write command
	fmt.Fprintf(outBuf, "add %s %d %d %d\r\n%s\r\n",
		param.key, param.flags, param.exptime, len(param.value)-1, param.value)

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
	if bytes.Equal(expectVer, recvBuf) == false {
		t.Errorf("add command response error. Expect:%x, Actual:%x\n",
			expectVer, recvBuf)
	}

	// rewrite to right response data
	param.value = param.value[:len(param.value)-1]

	// get check
	getCommand(t, con, param)
}

func TestReplaceDataOverBytesAndGet(t *testing.T) {
	con := makeConnection(t)
	defer con.Close()
	con.SetReadDeadline(time.Now().Add(10 * time.Second))
	con.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// set
	param := &setCommandParam{
		key:     []byte("ReplaceDataOverBytesKey"),
		value:   []byte("Value"),
		flags:   12345,
		exptime: 0,
	}
	setCommand(t, con, param)

	outBuf := new(bytes.Buffer)
	param.value = []byte("ReplaceDataOverBytesValue")

	// write command
	fmt.Fprintf(outBuf, "replace %s %d %d %d\r\n%s\r\n",
		param.key, param.flags, param.exptime, len(param.value)-1, param.value)

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
	if bytes.Equal(expectVer, recvBuf) == false {
		t.Errorf("add command response error. Expect:%x, Actual:%x\n",
			expectVer, recvBuf)
	}

	// rewrite to right response data
	param.value = param.value[:len(param.value)-1]

	// get check
	getCommand(t, con, param)
}
