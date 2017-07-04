package main

import (
	"net"
	"strconv"
	"strings"
)

func handleSet(con net.Conn, args []string) error {
	flags, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}
	exptime, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}
	bytes, err := strconv.Atoi(args[3])
	if err != nil {
		return err
	}

	buf := make([]byte, bytes)
	_, err = con.Read(buf)
	if err != nil {
		return err
	}
	data := string(buf[:bytes])
	data = strings.TrimRight(data, "\n\r")

	sv := &storeValue{}
	sv.key = []byte(args[0])
	sv.flags = flags
	sv.exptime = exptime
	sv.bytes = bytes
	sv.data = []byte(data)

	db.set(args[0], sv)
	con.Write([]byte("STORED\n"))

	return nil
}
