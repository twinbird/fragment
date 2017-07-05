package main

import (
	"log"
	"net"
	"strconv"
	"strings"
)

const MaxValueSize = 1000000

func handleGet(con net.Conn, args []string) error {
	key := args[0]
	v, err := db.get(key)
	if err != nil {
		return err
	}

	if v == nil {
		con.Write([]byte("END"))
	} else {
		con.Write([]byte(v.toString()))
	}

	return nil
}

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
	log.Println("before Read")
	dlen, err := con.Read(buf)
	log.Println("after Read")
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

func handleDelete(con net.Conn, args []string) error {
	key := args[0]
	ok, err := db.delete(key)
	if err != nil {
		return err
	}

	if ok == false {
		con.Write([]byte("NOT_DELETED"))
		return nil
	}
	con.Write([]byte("DELETED"))

	return nil
}