package testutils

import (
	gozk "launchpad.net/gozk"
	"log"
	"os"
)

type ZkTestServer struct {
	Server         *gozk.Server
	Zk             *ZkConn
	ZkEventChan    <-chan gozk.Event
	TestServerPort int
	TestServerDir  string
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func NewZkTestServer(port int) *ZkTestServer {
	return &ZkTestServer{TestServerPort: port, TestServerDir: "/tmp/zktest"}
}

func (z *ZkTestServer) Init() error {
	os.RemoveAll(z.TestServerDir)
	var err error
	zkPath := "/usr/local/Cellar/zookeeper/3.4.6/libexec"
	exists, err := exists(zkPath)
	if err != nil {
		return err
	}
	if !exists {
		zkPath = "/usr/share/java"
	}
	z.Server, err = gozk.CreateServer(z.TestServerPort, z.TestServerDir, zkPath)
	if err != nil {
		return err
	}
	err = z.Server.Start()
	if err != nil {
		return err
	}
	addr, err := z.Server.Addr()
	if err != nil {
		return err
	}
	z.Zk, z.ZkEventChan, err = GetZk(addr)
	if err != nil {
		return err
	}
	WaitOnConnect(z.ZkEventChan)

	//build the initial necessary paths
	log.Println(z.Zk.Conn.Create("/pools", "", 0, gozk.WorldACL(gozk.PERM_ALL)))
	log.Println(z.Zk.Conn.Create("/ports", "", 0, gozk.WorldACL(gozk.PERM_ALL)))
	log.Println(z.Zk.Conn.Create("/rules", "", 0, gozk.WorldACL(gozk.PERM_ALL)))
	log.Println(z.Zk.Conn.Create("/tries", "", 0, gozk.WorldACL(gozk.PERM_ALL)))

	return nil
}

func (z *ZkTestServer) Destroy() error {
	z.Zk.Conn.Close()
	err := z.Server.Destroy()
	if err != nil {
		return err
	}
	os.RemoveAll(z.TestServerDir)
	return nil
}
