package kksecret

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Cassandra struct {
	Name   string
	Path   string
	Writer CassandraMeta `json:"writer"`
	Reader CassandraMeta `json:"reader"`
}

type CassandraMeta struct {
	Hosts    []string `json:"hosts"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	CaPath   string   `json:"ca_path"`
}

func CassandraProfile(cassandraName string) *Cassandra {
	path := fmt.Sprintf("%scassandra-%s/cassandra.json", PATH, cassandraName)
	if _, e := os.Stat(path); os.IsNotExist(e) {
		return nil
	}

	if bytes, err := ioutil.ReadFile(path); err == nil {
		var cassandra = &Cassandra{}
		if err := json.Unmarshal(bytes, &cassandra); err != nil {
			return nil
		}

		cassandra.Name = cassandraName
		cassandra.Path = path
		return cassandra
	}

	return nil
}
