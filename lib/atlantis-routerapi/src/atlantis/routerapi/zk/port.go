package zk

import (
	cfg "atlantis/router/config"
	routerzk "atlantis/router/zk"
	"errors"
	"strconv"
)

func ListPorts() ([]uint16, error) {

	return routerzk.ListPorts(zkConn.Conn)
}

func GetPort(name string) (cfg.Port, error) {
	if name == "" {
		return cfg.Port{}, errors.New("Please specify a port")
	}

	pUint, err := strconv.ParseUint(name, 10, 16)

	if err != nil {
		return cfg.Port{}, err
	}

	return routerzk.GetPort(zkConn.Conn, uint16(pUint))
}

func SetPort(port cfg.Port) error {

	if port.Port == 0 {
		return errors.New("Please specify a port")
	} else if port.Trie == "" {
		return errors.New("Please specify a trie")
	}

	return routerzk.SetPort(zkConn.Conn, port)
}

func DeletePort(name string) error {
	if name == "" {
		return errors.New("Please specify a port")
	}

	pUint, err := strconv.ParseUint(name, 10, 16)

	if err != nil {
		return err
	}

	return routerzk.DelPort(zkConn.Conn, uint16(pUint))
}
