package main

import (
	"net"
	"time"
)

func handleGet(con net.Conn, args []string) error {
	key := args[0]
	v, err := db.get(key)
	if err != nil {
		return err
	}

	con.SetWriteDeadline(time.Now().Add(10 * time.Second))
	con.Write([]byte(v.toString()))

	return nil
}
