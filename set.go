package main

import (
	"net"
	"strconv"
	"strings"
)

const MaxValueSize = 1024

func handleSet(con net.Conn, args []string) error {
	flags, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}
	exptime, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}
	dbytes, err := strconv.Atoi(args[3])
	if err != nil {
		return err
	}

	buf := make([]byte, MaxValueSize)
	dlen, err := con.Read(buf)
	if err != nil {
		return err
	}
	if dlen > dbytes {
		con.Write([]byte("CLIENT_ERROR bad data chunk\n"))
		return nil
	}
	data := string(buf[:dbytes])
	data = strings.TrimRight(data, "\n\r")

	sv := &storeValue{}
	sv.key = []byte(args[0])
	sv.flags = flags
	sv.exptime = exptime
	sv.bytes = dbytes
	sv.data = []byte(data)

	if err = db.set(args[0], sv); err != nil {
		return err
	}
	con.Write([]byte("STORED\n"))

	return nil
}

func handleAdd(con net.Conn, args []string) error {
	flags, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}
	exptime, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}
	dbytes, err := strconv.Atoi(args[3])
	if err != nil {
		return err
	}

	buf := make([]byte, MaxValueSize)
	dlen, err := con.Read(buf)
	if err != nil {
		return err
	}
	if dlen > dbytes {
		con.Write([]byte("CLIENT_ERROR bad data chunk\n"))
		return nil
	}
	data := string(buf[:dbytes])
	data = strings.TrimRight(data, "\n\r")

	sv := &storeValue{}
	sv.key = []byte(args[0])
	sv.flags = flags
	sv.exptime = exptime
	sv.bytes = dbytes
	sv.data = []byte(data)

	if ok, err := db.add(args[0], sv); err != nil {
		return err
	} else if ok == false {
		con.Write([]byte("NOT_STORED\n"))
		return nil
	}
	con.Write([]byte("STORED\n"))

	return nil
}

func handleReplace(con net.Conn, args []string) error {
	flags, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}
	exptime, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}
	dbytes, err := strconv.Atoi(args[3])
	if err != nil {
		return err
	}

	buf := make([]byte, MaxValueSize)
	dlen, err := con.Read(buf)
	if err != nil {
		return err
	}
	if dlen > dbytes {
		con.Write([]byte("CLIENT_ERROR bad data chunk\n"))
		return nil
	}
	data := string(buf[:dbytes])
	data = strings.TrimRight(data, "\n\r")

	sv := &storeValue{}
	sv.key = []byte(args[0])
	sv.flags = flags
	sv.exptime = exptime
	sv.bytes = dbytes
	sv.data = []byte(data)

	if ok, err := db.replace(args[0], sv); err != nil {
		return err
	} else if ok == false {
		con.Write([]byte("NOT_STORED\n"))
		return nil
	}
	con.Write([]byte("STORED\n"))

	return nil
}
