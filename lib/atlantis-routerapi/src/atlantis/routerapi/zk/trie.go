package zk

import (
	cfg "atlantis/router/config"
	routerzk "atlantis/router/zk"
	"errors"
)

func ListTries() ([]string, error) {

	return routerzk.ListTries(zkConn.Conn)
}

func GetTrie(name string) (cfg.Trie, error) {
	if name == "" {
		return cfg.Trie{}, errors.New("Please specify a name")
	}

	return routerzk.GetTrie(zkConn.Conn, name)
}

func SetTrie(trie cfg.Trie) error {

	if trie.Name == "" {
		return errors.New("Please specify a name")
	} else if len(trie.Rules) <= 0 {
		return errors.New("Please specify a rule")
	}

	return routerzk.SetTrie(zkConn.Conn, trie)
}

func DeleteTrie(name string) error {
	if name == "" {
		return errors.New("Please specify a name")
	}

	return routerzk.DelTrie(zkConn.Conn, name)
}
