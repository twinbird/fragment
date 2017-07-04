package main

import (
	"net"
)

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
