package secret

import (
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCassandraProfile(t *testing.T) {
	oldPath := PATH
	PATH = "./"
	defer func() { PATH = oldPath }()

	testDir := path.Join(".", "cassandra-test")
	os.MkdirAll(testDir, fs.ModePerm)
	secretDir := path.Join(testDir, "secret.json")
	defer os.RemoveAll(testDir)

	testData := `{
  "writer": {
    "endpoints": ["wh1", "wh2"],
    "username": "wu",
    "password": "wp",
    "ca_path": "wcp"
  },
  "reader": {
    "endpoints": ["rh1"],
    "username": "ru",
    "password": "rp",
    "ca_path": "rcp"
  }
}`

	os.WriteFile(secretDir, []byte(testData), fs.ModePerm)

	c := &Cassandra{}
	err := Load("cassandra", "test", c)

	assert.NoError(t, err)
	assert.Equal(t, "test", c.Name())
	assert.NotEmpty(t, c.Path())

	assert.Equal(t, 2, len(c.Writer.Endpoints))
	assert.Equal(t, "wh2", c.Writer.Endpoints[1])
	assert.Equal(t, "wu", c.Writer.Username)
	assert.Equal(t, "wp", c.Writer.Password)
	assert.Equal(t, "wcp", c.Writer.CaPath)

	assert.Equal(t, 1, len(c.Reader.Endpoints))
	assert.Equal(t, "ru", c.Reader.Username)
	assert.Equal(t, "rp", c.Reader.Password)
	assert.Equal(t, "rcp", c.Reader.CaPath)
}
