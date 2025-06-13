package secret

import (
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseProfile(t *testing.T) {
	oldPath := PATH
	PATH = "./"
	defer func() { PATH = oldPath }()

	testDir := path.Join(".", "database-test")
	os.MkdirAll(testDir, fs.ModePerm)
	secretDir := path.Join(testDir, "secret.json")
	defer os.RemoveAll(testDir)

	testData := `{
	"writer": {
		"adapter": "mysql",
		"params": {
			"charset": "utf8mb4",
			"host": "localhost",
			"port": 3306,
			"dbname": "test_db",
			"username": "test_user",
			"password": "test_password"
		}
	},
	"reader": {
		"adapter": "mysql",
		"params": {
			"charset": "utf8mb4",
			"host": "readonly.localhost",
			"port": 3306,
			"dbname": "test_db_readonly",
			"username": "readonly_user",
			"password": "readonly_password"
		}
	}
}`

	os.WriteFile(secretDir, []byte(testData), fs.ModePerm)

	db := &Database{}
	err := Load("database", "test", db)

	assert.NoError(t, err)
	assert.Equal(t, "test", db.Name())
	assert.NotEmpty(t, db.Path())

	assert.Equal(t, "mysql", db.Writer.Adapter)
	assert.Equal(t, "utf8mb4", db.Writer.Params.Charset)
	assert.Equal(t, "localhost", db.Writer.Params.Host)
	assert.Equal(t, uint(3306), db.Writer.Params.Port)
	assert.Equal(t, "test_db", db.Writer.Params.DBName)
	assert.Equal(t, "test_user", db.Writer.Params.Username)
	assert.Equal(t, "test_password", db.Writer.Params.Password)

	assert.Equal(t, "mysql", db.Reader.Adapter)
	assert.Equal(t, "utf8mb4", db.Reader.Params.Charset)
	assert.Equal(t, "readonly.localhost", db.Reader.Params.Host)
	assert.Equal(t, uint(3306), db.Reader.Params.Port)
	assert.Equal(t, "test_db_readonly", db.Reader.Params.DBName)
	assert.Equal(t, "readonly_user", db.Reader.Params.Username)
	assert.Equal(t, "readonly_password", db.Reader.Params.Password)
}
