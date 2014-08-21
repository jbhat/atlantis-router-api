package zk

import (
	"launchpad.net/gozk"
	"log"
)

var (
	zkConn   *ZkConn
	LastZkOk = true
)

func SetZkConn(c *zookeeper.Conn, evc <-chan zookeeper.Event, servers string) {

	zkConn = &ZkConn{
		ResetCh: make(chan bool),
		servers: servers,
		killCh:  make(chan bool),
		Conn:    c,
		eventCh: evc,
	}

	go zkConn.monitorEventCh()
	//TODO:figure out if/when i need to do this
	//zkConn.ResetCh <- true
}

//This is used only when SetZkConn is used to allow
//the ZkConn to be closed by whatever created it and still
//allow the killCh to be signaled
func KillConnection() {
	zkConn.killCh <- true
}

func Init(zkAddr string, config bool) {

	zkConn = ManagedZkConn(zkAddr)
	if config {
		SetupZk()
	}

}

func SetupZk() {

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

//To check health of ZK we check what the event status was during the
//last check and the current check. If they are both false then
//return critical, otherwise continue. The reason for this is the
//ZK has built in functionality to fix a broken connection. So, the
//ZK could be reconnecting from an expired session when healthz was called.
//So we say it has to be down for a whole interval of healthz checks.
func IsZkConnOk() bool {

	if zkConn == nil {
		return false
	}

	//if any of the necessary paths are not created
	if _, err := zkConn.Conn.Exists("/pools"); err != nil {
		return false
	}

	if _, err := zkConn.Conn.Exists("/rules"); err != nil {
		return false
	}

	if _, err := zkConn.Conn.Exists("/tries"); err != nil {
		return false
	}

	if _, err := zkConn.Conn.Exists("/ports"); err != nil {
		return false
	}

	//attempt to create/delete a simple test path to be sure ZK is working
	_, err := zkConn.Conn.Create("/test", "", 0, zookeeper.WorldACL(zookeeper.PERM_ALL))
	if err != nil {
		return false
	}

	err = zkConn.Conn.Delete("/test", -1)
	if err != nil {
		return false
	}

	//check current status and compare to last
	newOk := zkConn.IsZkOk()
	//if both false
	if !LastZkOk && !newOk {
		return false
	}

	LastZkOk = newOk

	return true
}
