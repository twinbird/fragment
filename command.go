package main

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
)

type Command struct {
	command []byte
	key     []byte
	value   []byte
	exptime int
	flags   int
	bytes   int
	parsed  bool
}

func parseCommand(buf []byte) (*Command, error) {
	idx := bytes.Index(buf, []byte("\r\n"))
	if idx < 0 {
		return nil, fmt.Errorf("invalid command")
	}

	comBuf := buf[:idx]
	comPrm := bytes.Split(comBuf, []byte(" "))

	if len(comPrm) < 1 {
		return nil, fmt.Errorf("invalid command")
	}

	com := &Command{}
	com.command = comPrm[0]

	switch {
	case bytes.Equal(comPrm[0], []byte("set")):
		if len(comPrm) != 5 {
			return nil, fmt.Errorf("invalid command")
		}
		com.key = comPrm[1]
		if v, err := strconv.Atoi(string(comPrm[2])); err != nil {
			return nil, fmt.Errorf("invalid command")
		} else {
			com.flags = v
		}
		if v, err := strconv.Atoi(string(comPrm[3])); err != nil {
			return nil, fmt.Errorf("invalid command")
		} else {
			com.exptime = v
		}
		if v, err := strconv.Atoi(string(comPrm[4])); err != nil {
			return nil, fmt.Errorf("invalid command")
		} else {
			com.bytes = v
		}
		// there is no data yet
		lastIdx := bytes.LastIndex(buf, []byte("\r\n"))
		if lastIdx == idx {
			return com, nil
		}
		if lastIdx < 0 {
			return nil, fmt.Errorf("invalid command")
		}
		com.value = buf[idx+len("\r\n") : lastIdx]
		com.parsed = true
	case bytes.Equal(comPrm[0], []byte("add")):
		if len(comPrm) != 5 {
			return nil, fmt.Errorf("invalid command")
		}
		com.key = comPrm[1]
		if v, err := strconv.Atoi(string(comPrm[2])); err != nil {
			return nil, fmt.Errorf("invalid command")
		} else {
			com.flags = v
		}
		if v, err := strconv.Atoi(string(comPrm[3])); err != nil {
			return nil, fmt.Errorf("invalid command")
		} else {
			com.exptime = v
		}
		if v, err := strconv.Atoi(string(comPrm[4])); err != nil {
			return nil, fmt.Errorf("invalid command")
		} else {
			com.bytes = v
		}
		// there is no data yet
		lastIdx := bytes.LastIndex(buf, []byte("\r\n"))
		if lastIdx == idx {
			return com, nil
		}
		if lastIdx < 0 {
			return nil, fmt.Errorf("invalid command")
		}
		com.value = buf[idx+len("\r\n") : lastIdx]
		com.parsed = true
	case bytes.Equal(comPrm[0], []byte("replace")):
		if len(comPrm) != 5 {
			return nil, fmt.Errorf("invalid command")
		}
		com.key = comPrm[1]
		if v, err := strconv.Atoi(string(comPrm[2])); err != nil {
			return nil, fmt.Errorf("invalid command")
		} else {
			com.flags = v
		}
		if v, err := strconv.Atoi(string(comPrm[3])); err != nil {
			return nil, fmt.Errorf("invalid command")
		} else {
			com.exptime = v
		}
		if v, err := strconv.Atoi(string(comPrm[4])); err != nil {
			return nil, fmt.Errorf("invalid command")
		} else {
			com.bytes = v
		}
		// there is no data yet
		lastIdx := bytes.LastIndex(buf, []byte("\r\n"))
		if lastIdx == idx {
			return com, nil
		}
		if lastIdx < 0 {
			return nil, fmt.Errorf("invalid command")
		}
		com.value = buf[idx+len("\r\n") : lastIdx]
		com.parsed = true
	case bytes.Equal(comPrm[0], []byte("delete")):
		if len(comPrm) != 2 {
			return nil, fmt.Errorf("invalid command")
		}
		com.key = comPrm[1]
		com.parsed = true
	case bytes.Equal(comPrm[0], []byte("get")):
		if len(comPrm) != 2 {
			return nil, fmt.Errorf("invalid command")
		}
		com.key = comPrm[1]
		com.parsed = true
	case bytes.Equal(comPrm[0], []byte("version")):
		com.parsed = true
	default:
		return nil, fmt.Errorf("invalid command: '%s'", string(comPrm[0]))
	}

	return com, nil
}

func handleGet(con net.Conn, com *Command) error {
	v, err := db.get(com.key)
	if err != nil {
		return err
	}

	if v == nil {
		con.Write([]byte("END\r\n"))
	} else {
		con.Write([]byte(v.toString()))
	}

	return nil
}

func handleSet(con net.Conn, com *Command) error {
	sv := &storeValue{}
	sv.key = []byte(com.key)
	sv.flags = com.flags
	sv.exptime = com.exptime
	sv.bytes = com.bytes
	sv.data = com.value

	if err := db.set(com.key, sv); err != nil {
		return err
	}
	con.Write([]byte("STORED\r\n"))

	return nil
}

func handleAdd(con net.Conn, com *Command) error {
	sv := &storeValue{}
	sv.key = []byte(com.key)
	sv.flags = com.flags
	sv.exptime = com.exptime
	sv.bytes = com.bytes
	sv.data = com.value

	if ok, err := db.add(com.key, sv); err != nil {
		return err
	} else if ok == false {
		con.Write([]byte("NOT_STORED\r\n"))
		return nil
	}
	con.Write([]byte("STORED\r\n"))

	return nil
}

func handleReplace(con net.Conn, com *Command) error {
	sv := &storeValue{}
	sv.key = com.key
	sv.flags = com.flags
	sv.exptime = com.exptime
	sv.bytes = com.bytes
	sv.data = com.value

	if ok, err := db.replace(com.key, sv); err != nil {
		return err
	} else if ok == false {
		con.Write([]byte("NOT_STORED\r\n"))
		return nil
	}
	con.Write([]byte("STORED\r\n"))

	return nil
}

func handleDelete(con net.Conn, com *Command) error {
	ok, err := db.delete(com.key)
	if err != nil {
		return err
	}

	if ok == false {
		con.Write([]byte("NOT_FOUND\r\n"))
		return nil
	}
	con.Write([]byte("DELETED\r\n"))

	return nil
}

func handleVersion(con net.Conn, com *Command) error {
	con.Write([]byte("0.0.1\r\n"))
	return nil
}
