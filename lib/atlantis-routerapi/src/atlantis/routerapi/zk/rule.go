package zk


import (
	"errors"
	routerzk "atlantis/router/zk"
	cfg "atlantis/router/config"
)

func ListRules() ([]string, error){
	
	return routerzk.ListRules(zkConn.Conn)
}

func GetRule(name string) (cfg.Rule, error){
	if name == "" {
		return cfg.Rule{}, errors.New("Please specify a name")
	}
	
	return routerzk.GetRule(zkConn.Conn, name)
}

func SetRule(rule cfg.Rule) error {

	if rule.Name == "" {
		return errors.New("Please specify a name")
	} else if rule.Type == "" {
		return errors.New("Please specify a type")
	} else if rule.Value == "" {
		return errors.New("Please specify a value")
	} else if rule.Next == "" {
		return errors.New("Please specify a next value")
	} else if rule.Pool == "" {
		return errors.New("Please specify a pool")
	}

	return routerzk.SetRule(zkConn.Conn, rule)
}


func DeleteRule(name string) error {
	if name == "" {
		return errors.New("Please specify a name")
	}

	return routerzk.DelRule(zkConn.Conn, name)
}
