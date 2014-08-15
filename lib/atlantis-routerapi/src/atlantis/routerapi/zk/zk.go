package zk


import (
	"log"
	"launchpad.net/gozk"
)

var (
	zkConn	*ZkConn
)



func Init(zkAddr string, config bool) {

	zkConn = ManagedZkConn(zkAddr)
	if config {
		SetupZk()
	}
	
}

func SetupZk(){

	for {
		ev := <-zkConn.ResetCh
	
		if ev == false {
			log.Println("not true")
			continue
		} else {
			log.Println("yay it's true")
			break
		}

	}

	log.Println(zkConn.Conn.Create("/pools", "", 0, zookeeper.WorldACL(zookeeper.PERM_ALL)))
	log.Println(zkConn.Conn.Create("/ports", "", 0, zookeeper.WorldACL(zookeeper.PERM_ALL)))
	log.Println(zkConn.Conn.Create("/rules", "", 0, zookeeper.WorldACL(zookeeper.PERM_ALL)))
	log.Println(zkConn.Conn.Create("/tries", "", 0, zookeeper.WorldACL(zookeeper.PERM_ALL)))

}
