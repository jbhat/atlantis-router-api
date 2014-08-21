package client

import (
	cfg "atlantis/router/config"
	"encoding/json"
	"errors"
	"fmt"
)

func (c *AddPoolCommand) BuildPool() (cfg.Pool, error) {
	if c.Name == "" {
		return cfg.Pool{}, errors.New("Please specify a name for your pool")
	} else if c.HealthCheckEvery == "" {
		return cfg.Pool{}, errors.New("Please specify a HealthzCheckEvery for your pool")
	} else if c.HealthzTimeout == "" {
		return cfg.Pool{}, errors.New("Please specify a HealthzTimeout for your pool")
	} else if c.RequestTimeout == "" {
		return cfg.Pool{}, errors.New("Please specify a RequestTimeout for your pool")
	} else if c.Status == "" {
		return cfg.Pool{}, errors.New("Please specify a Status for your pool")
	} else if len(c.Hosts) == 0 {
		return cfg.Pool{}, errors.New("Please specify at least one host for your pool")
	}

	hMap := make(map[string]cfg.Host, len(c.Hosts))
	for key, value := range c.Hosts {
		fmt.Printf("myval: %v : %v\n", key, value)
		hMap[value] = cfg.Host{value}
	}
	return cfg.Pool{c.Name,
		false,
		hMap,
		cfg.PoolConfig{c.HealthCheckEvery, c.HealthzTimeout,
			c.RequestTimeout, c.Status}}, nil
}

func (c *AddPoolCommand) GetPoolJson() (string, error) {

	pool, err := c.BuildPool()
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(pool)
	if err != nil {
		return "", err
	}

	return string(b), nil

}

func (c *AddRuleCommand) BuildRule() (cfg.Rule, error) {
	if c.Name == "" {
		return cfg.Rule{}, errors.New("Please specify a name for your rule")
	} else if c.Type == "" {
		return cfg.Rule{}, errors.New("Please specify a type for your rule")
	} else if c.Value == "" {
		return cfg.Rule{}, errors.New("Please specify a value for your rule")
	} else if c.Pool == "" {
		return cfg.Rule{}, errors.New("Please specify a pool for your rule")
	}
	return cfg.Rule{c.Name, c.Type, c.Value, c.Next, c.Pool, false}, nil
}

func (c *AddRuleCommand) GetRuleJson() (string, error) {

	rule, err := c.BuildRule()
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(rule)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (c *AddTrieCommand) BuildTrie() (cfg.Trie, error) {
	if c.Name == "" {
		return cfg.Trie{}, errors.New("Please specify a name for your trie")
	} else if len(c.Rules) == 0 {
		return cfg.Trie{}, errors.New("Please specify at least one rule for your trie")
	}
	return cfg.Trie{c.Name, c.Rules, false}, nil
}

func (c *AddTrieCommand) GetTrieJson() (string, error) {

	trie, err := c.BuildTrie()
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(trie)
	if err != nil {
		return "", err
	}

	return string(b), nil

}

func (c *AddPortCommand) BuildPort() (cfg.Port, error) {
	if c.Port <= 0 {
		return cfg.Port{}, errors.New("Please specify a valid port number")
	} else if c.Trie == "" {
		return cfg.Port{}, errors.New("Please specify a root trie for this port")
	}
	return cfg.Port{c.Port, c.Trie, false}, nil
}

func (c *AddPortCommand) GetPortJson() (string, error) {

	port, err := c.BuildPort()
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(port)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func ExpandAndPrintData(uri, jName, jListData string) error {

	var m map[string][]interface{}
	err := json.Unmarshal([]byte(jListData), &m)
	if err != nil {
		return err
	}

	vList := m[jName]
	for _, value := range vList {
		instatusCode, indata, err := BuildAndSendRequest("GET", fmt.Sprintf("%s%v", uri, value), "")
		if err != nil {
			return err
		}

		Output(instatusCode, indata)
	}
	return nil
}
