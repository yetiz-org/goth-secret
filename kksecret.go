package kksecret

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var PATH = "/the-secret/"

type Database struct {
	Name   string
	Path   string
	Writer DatabaseMeta `json:"writer"`
	Reader DatabaseMeta `json:"reader"`
}

type Redis struct {
	Name   string
	Path   string
	Master RedisMeta `json:"master"`
	Slave  RedisMeta `json:"slave"`
}

type DatabaseMeta struct {
	Adapter string
	Params  struct {
		Charset  string `json:"charset"`
		Host     string `json:"host"`
		DBName   string `json:"dbname"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"params"`
}

type RedisMeta struct {
	Host string `json:"host"`
	Port uint   `json:"port"`
}

func DatabaseProfile(dbname string) *Database {
	path := fmt.Sprintf("%sdb-%s/db.json", PATH, dbname)
	if _, e := os.Stat(path); os.IsNotExist(e) {
		return nil
	}

	if bytes, err := ioutil.ReadFile(path); err == nil {
		database := &Database{}
		if err := json.Unmarshal(bytes, database); err != nil {
			return nil
		}

		database.Name = dbname
		database.Path = path
		return database
	}

	return nil
}

func RedisProfile(redisName string) *Redis {
	path := fmt.Sprintf("%sredis-%s/redis.json", PATH, redisName)
	if _, e := os.Stat(path); os.IsNotExist(e) {
		return nil
	}

	if bytes, err := ioutil.ReadFile(path); err == nil {
		var redis = &Redis{}
		if err := json.Unmarshal(bytes, redis); err != nil {
			return nil
		}

		redis.Name = redisName
		redis.Path = path
		return redis
	}

	return nil
}
