package secret

import (
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedisProfile(t *testing.T) {
	// 設置測試環境
	oldPath := PATH
	PATH = "./"
	defer func() { PATH = oldPath }()

	// 創建測試目錄和文件
	testDir := path.Join(".", "redis-test")
	os.MkdirAll(testDir, fs.ModePerm)
	secretDir := path.Join(testDir, "secret.json")
	defer os.RemoveAll(testDir)

	// 創建測試數據
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

	// 創建 Redis 實例並載入數據
	redis := &Redis{}
	err := Load("redis", "test", redis)

	// 驗證
	assert.NoError(t, err)
	assert.Equal(t, "test", redis.Name())
	assert.NotEmpty(t, redis.Path())

	// 驗證 Master 欄位
	assert.Equal(t, "localhost", redis.Master.Host)
	assert.Equal(t, uint(6379), redis.Master.Port)

	// 驗證 Slave 欄位
	assert.Equal(t, "slave.localhost", redis.Slave.Host)
	assert.Equal(t, uint(6380), redis.Slave.Port)
}
