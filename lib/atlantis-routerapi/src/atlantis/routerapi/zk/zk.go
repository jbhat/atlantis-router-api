package zk

import (
	"launchpad.net/gozk"
)



var (
	ZkObj	*ZkConn
	zkConn	*zookeeper.Conn
)



func Init(zkAddr string) {

	ZkObj = ManagedZkConn(zkAddr)
	zkConn = ZkObj.Conn
}
