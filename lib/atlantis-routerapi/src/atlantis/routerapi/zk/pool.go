package zk


import (
	"errors"
	routerzk "atlantis/router/zk"
	cfg "atlantis/router/config"
)

func ListPools() ([]string, error){
	
	return routerzk.ListPools(zkConn.Conn)
}

func GetPool(name string) (cfg.Pool, error){
	if name == "" {
		return cfg.Pool{}, errors.New("Please specify a name")
	}
	
	return routerzk.GetPool(zkConn.Conn, name)
}

func SetPool(pool cfg.Pool) error {

	if pool.Name == "" {
                return errors.New("Please specify a name")
        } else if pool.Config.HealthzEvery == "" {
                return errors.New("Please specify a healthz check frequency")
        } else if pool.Config.HealthzTimeout == "" {
                return errors.New("Please specify a healthz timeout")
        } else if pool.Config.RequestTimeout == "" {
                return errors.New("Please specify a request timeout")
        } // no need to check hosts. an empty pool is still a valid pool


	return routerzk.SetPool(zkConn.Conn, pool)
}


func DeletePool(name string) error {
	if name == "" {
		return errors.New("Please specify a name")
	}

	return routerzk.DelPool(zkConn.Conn, name)
}
