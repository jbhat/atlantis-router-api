package zk

import (
	cfg "atlantis/router/config"
	routerzk "atlantis/router/zk"
	"errors"
)

func GetHosts(poolName string) (map[string]cfg.Host, error) {

	if poolName == "" {
		return nil, errors.New("Please specify a pool name to get the hosts from")
	}

	return routerzk.GetHosts(zkConn.Conn, poolName)

}

func AddHosts(poolName string, hosts map[string]cfg.Host) error {

	if poolName == "" {
		return errors.New("Please specify a pool name to add the hosts to")
	} else if len(hosts) == 0 {
		return errors.New("Please specify at least one host to add to the pool")
	}

	return routerzk.AddHosts(zkConn.Conn, poolName, hosts)
}

func DeleteHosts(poolName string, hosts []string) error {

	if poolName == "" {
		return errors.New("Please specify a pool name to delete hosts from")
	} else if len(hosts) == 0 {
		return errors.New("Please specifiy at least one host to delete from the pool")
	}

	return routerzk.DelHosts(zkConn.Conn, poolName, hosts)
}
