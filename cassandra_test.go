package kksecret

import (
	"io/fs"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCassandraProfile(t *testing.T) {
	data := "{\n  \"writer\": {\n    \"hosts\": [\"wh1\", \"wh2\"],\n    \"username\": \"wu\",\n    \"password\": \"wp\",\n    \"ca_path\": \"wcp\"\n  },\n  \"reader\": {\n    \"hosts\": [\"rh1\"],\n    \"username\": \"ru\",\n    \"password\": \"rp\",\n    \"ca_path\": \"rcp\"\n  }\n}"
	os.Mkdir("cassandra-test", fs.ModePerm)
	ioutil.WriteFile("cassandra-test/cassandra.json", []byte(data), fs.ModePerm)
	PATH = "./"
	csd := CassandraProfile("test")
	os.RemoveAll("cassandra-test")
	assert.NotNil(t, csd)
	assert.Equal(t, "test", csd.Name)
	assert.Equal(t, 2, len(csd.Writer.Hosts))
	assert.Equal(t, "wh2", csd.Writer.Hosts[1])
	assert.Equal(t, "wu", csd.Writer.Username)
	assert.Equal(t, "wp", csd.Writer.Password)
	assert.Equal(t, "wcp", csd.Writer.CaPath)
	assert.Equal(t, 1, len(csd.Reader.Hosts))
	assert.Equal(t, "ru", csd.Reader.Username)
	assert.Equal(t, "rp", csd.Reader.Password)
	assert.Equal(t, "rcp", csd.Reader.CaPath)
}
