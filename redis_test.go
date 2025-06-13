package secret

import (
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedisProfile(t *testing.T) {
	oldPath := PATH
	PATH = "./"
	defer func() { PATH = oldPath }()

	testDir := path.Join(".", "redis-test")
	os.MkdirAll(testDir, fs.ModePerm)
	secretDir := path.Join(testDir, "secret.json")
	defer os.RemoveAll(testDir)

	testData := `{
	"master": {
		"host": "localhost",
		"port": 6379
	},
	"slave": {
		"host": "slave.localhost",
		"port": 6380
	}
}`

	os.WriteFile(secretDir, []byte(testData), fs.ModePerm)

	redis := &Redis{}
	err := Load("redis", "test", redis)

	assert.NoError(t, err)
	assert.Equal(t, "test", redis.Name())
	assert.NotEmpty(t, redis.Path())

	assert.Equal(t, "localhost", redis.Master.Host)
	assert.Equal(t, uint(6379), redis.Master.Port)

	assert.Equal(t, "slave.localhost", redis.Slave.Host)
	assert.Equal(t, uint(6380), redis.Slave.Port)
}
