package zk

import (
	"launchpad.net/gozk"
)



var (
	zkConn	*ZkConn
)



func Init(zkAddr string) {

	zkConn = ManagedZkConn(zkAddr)
}
