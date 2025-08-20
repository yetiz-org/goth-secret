package secret

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHelper provides utilities for testing secret loading functionality
type TestHelper struct {
	originalPath string
	tempDir      string
	mockFS       *MockFileSystem
	mockEnv      *MockEnvironment
}

// NewTestHelper creates a new test helper instance
func NewTestHelper() *TestHelper {
	return &TestHelper{
		mockFS:  NewMockFileSystem(),
		mockEnv: NewMockEnvironment(),
	}
}

// SetupTempDirectory creates a temporary directory for testing and sets the secret path
func (th *TestHelper) SetupTempDirectory(t *testing.T, tempPath string) func() {
	th.originalPath = PATH
	PATH = tempPath
	
	return func() {
		PATH = th.originalPath
		if th.tempDir != "" {
			os.RemoveAll(th.tempDir)
		}
	}
}

// CreateSecretFile creates a secret file with the given content in the temporary directory
func (th *TestHelper) CreateSecretFile(t *testing.T, typ, name string, content interface{}) (string, func()) {
	testDir := path.Join(".", fmt.Sprintf("%s-%s", typ, name))
	err := os.MkdirAll(testDir, fs.ModePerm)
	assert.NoError(t, err)
	
	secretPath := path.Join(testDir, "secret.json")
	
	var jsonData []byte
	if str, ok := content.(string); ok {
		jsonData = []byte(str)
	} else {
		jsonData, err = json.Marshal(content)
		assert.NoError(t, err)
	}
	
	err = os.WriteFile(secretPath, jsonData, fs.ModePerm)
	assert.NoError(t, err)
	
	th.tempDir = testDir
	
	return secretPath, func() {
		os.RemoveAll(testDir)
	}
}

// CreateMockSecretFile adds a secret file to the mock file system
func (th *TestHelper) CreateMockSecretFile(typ, name string, content interface{}) error {
	return th.mockFS.AddSecretFile(typ, name, content)
}

// SetMockPath sets the path in the mock environment
func (th *TestHelper) SetMockPath(path string) {
	th.mockEnv.SetVar("GOTH_SECRET_PATH", path)
}

// GetMockFileSystem returns the mock file system for advanced usage
func (th *TestHelper) GetMockFileSystem() *MockFileSystem {
	return th.mockFS
}

// GetMockEnvironment returns the mock environment for advanced usage
func (th *TestHelper) GetMockEnvironment() *MockEnvironment {
	return th.mockEnv
}

// AssertDatabaseEquals asserts that two Database structs are equal
func AssertDatabaseEquals(t *testing.T, expected, actual *Database) {
	assert.Equal(t, expected.Name(), actual.Name())
	assert.Equal(t, expected.Writer.Adapter, actual.Writer.Adapter)
	assert.Equal(t, expected.Writer.Params.Charset, actual.Writer.Params.Charset)
	assert.Equal(t, expected.Writer.Params.Host, actual.Writer.Params.Host)
	assert.Equal(t, expected.Writer.Params.Port, actual.Writer.Params.Port)
	assert.Equal(t, expected.Writer.Params.DBName, actual.Writer.Params.DBName)
	assert.Equal(t, expected.Writer.Params.Username, actual.Writer.Params.Username)
	assert.Equal(t, expected.Writer.Params.Password, actual.Writer.Params.Password)
	
	assert.Equal(t, expected.Reader.Adapter, actual.Reader.Adapter)
	assert.Equal(t, expected.Reader.Params.Charset, actual.Reader.Params.Charset)
	assert.Equal(t, expected.Reader.Params.Host, actual.Reader.Params.Host)
	assert.Equal(t, expected.Reader.Params.Port, actual.Reader.Params.Port)
	assert.Equal(t, expected.Reader.Params.DBName, actual.Reader.Params.DBName)
	assert.Equal(t, expected.Reader.Params.Username, actual.Reader.Params.Username)
	assert.Equal(t, expected.Reader.Params.Password, actual.Reader.Params.Password)
}

// AssertRedisEquals asserts that two Redis structs are equal
func AssertRedisEquals(t *testing.T, expected, actual *Redis) {
	assert.Equal(t, expected.Name(), actual.Name())
	assert.Equal(t, expected.Master.Host, actual.Master.Host)
	assert.Equal(t, expected.Master.Port, actual.Master.Port)
	assert.Equal(t, expected.Slave.Host, actual.Slave.Host)
	assert.Equal(t, expected.Slave.Port, actual.Slave.Port)
}

// AssertCassandraEquals asserts that two Cassandra structs are equal
func AssertCassandraEquals(t *testing.T, expected, actual *Cassandra) {
	assert.Equal(t, expected.Name(), actual.Name())
	assert.Equal(t, expected.Writer.Endpoints, actual.Writer.Endpoints)
	assert.Equal(t, expected.Writer.Keyspace, actual.Writer.Keyspace)
	assert.Equal(t, expected.Writer.Username, actual.Writer.Username)
	assert.Equal(t, expected.Writer.Password, actual.Writer.Password)
	assert.Equal(t, expected.Writer.CaPath, actual.Writer.CaPath)
	
	assert.Equal(t, expected.Reader.Endpoints, actual.Reader.Endpoints)
	assert.Equal(t, expected.Reader.Keyspace, actual.Reader.Keyspace)
	assert.Equal(t, expected.Reader.Username, actual.Reader.Username)
	assert.Equal(t, expected.Reader.Password, actual.Reader.Password)
	assert.Equal(t, expected.Reader.CaPath, actual.Reader.CaPath)
}

// CreateDatabaseSecret creates a Database secret with default values for testing
func CreateDatabaseSecret(name string) *Database {
	db := &Database{}
	db.DefaultSecret._Name = name
	db.DefaultSecret._Path = fmt.Sprintf("database-%s/secret.json", name)
	
	// Set default writer configuration
	db.Writer.Adapter = "mysql"
	db.Writer.Params.Charset = "utf8mb4"
	db.Writer.Params.Host = "localhost"
	db.Writer.Params.Port = 3306
	db.Writer.Params.DBName = "test_db"
	db.Writer.Params.Username = "test_user"
	db.Writer.Params.Password = "test_password"
	
	// Set default reader configuration
	db.Reader.Adapter = "mysql"
	db.Reader.Params.Charset = "utf8mb4"
	db.Reader.Params.Host = "readonly.localhost"
	db.Reader.Params.Port = 3306
	db.Reader.Params.DBName = "test_db_readonly"
	db.Reader.Params.Username = "readonly_user"
	db.Reader.Params.Password = "readonly_password"
	
	return db
}

// CreateRedisSecret creates a Redis secret with default values for testing
func CreateRedisSecret(name string) *Redis {
	redis := &Redis{}
	redis.DefaultSecret._Name = name
	redis.DefaultSecret._Path = fmt.Sprintf("redis-%s/secret.json", name)
	
	// Set default master configuration
	redis.Master.Host = "localhost"
	redis.Master.Port = 6379
	
	// Set default slave configuration
	redis.Slave.Host = "slave.localhost"
	redis.Slave.Port = 6380
	
	return redis
}

// CreateCassandraSecret creates a Cassandra secret with default values for testing
func CreateCassandraSecret(name string) *Cassandra {
	cassandra := &Cassandra{}
	cassandra.DefaultSecret._Name = name
	cassandra.DefaultSecret._Path = fmt.Sprintf("cassandra-%s/secret.json", name)
	
	// Set default writer configuration
	cassandra.Writer.Endpoints = []string{"localhost:9042", "backup.localhost:9042"}
	cassandra.Writer.Keyspace = "test_keyspace"
	cassandra.Writer.Username = "writer_user"
	cassandra.Writer.Password = "writer_password"
	cassandra.Writer.CaPath = "/path/to/writer/ca.pem"
	
	// Set default reader configuration
	cassandra.Reader.Endpoints = []string{"readonly.localhost:9042"}
	cassandra.Reader.Keyspace = "test_keyspace"
	cassandra.Reader.Username = "reader_user"
	cassandra.Reader.Password = "reader_password"
	cassandra.Reader.CaPath = "/path/to/reader/ca.pem"
	
	return cassandra
}

// ErrorTestCase represents a test case for error scenarios
type ErrorTestCase struct {
	Name          string
	SetupFunc     func(*TestHelper) error
	SecretType    string
	SecretName    string
	Secret        Secret
	ExpectedError string
}

// RunErrorTests runs a set of error test cases
func RunErrorTests(t *testing.T, testCases []ErrorTestCase) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			helper := NewTestHelper()
			cleanup := helper.SetupTempDirectory(t, "./")
			defer cleanup()
			
			if tc.SetupFunc != nil {
				err := tc.SetupFunc(helper)
				assert.NoError(t, err)
			}
			
			err := Load(tc.SecretType, tc.SecretName, tc.Secret)
			assert.Error(t, err)
			if tc.ExpectedError != "" {
				assert.Contains(t, err.Error(), tc.ExpectedError)
			}
		})
	}
}

// BenchmarkHelper provides utilities for performance testing
type BenchmarkHelper struct {
	testData map[string]interface{}
}

// NewBenchmarkHelper creates a new benchmark helper
func NewBenchmarkHelper() *BenchmarkHelper {
	return &BenchmarkHelper{
		testData: make(map[string]interface{}),
	}
}

// SetupBenchmarkData prepares test data for benchmarks
func (bh *BenchmarkHelper) SetupBenchmarkData() {
	bh.testData["database"] = map[string]interface{}{
		"writer": map[string]interface{}{
			"adapter": "mysql",
			"params": map[string]interface{}{
				"charset":  "utf8mb4",
				"host":     "localhost",
				"port":     3306,
				"dbname":   "test_db",
				"username": "test_user",
				"password": "test_password",
			},
		},
		"reader": map[string]interface{}{
			"adapter": "mysql",
			"params": map[string]interface{}{
				"charset":  "utf8mb4",
				"host":     "readonly.localhost",
				"port":     3306,
				"dbname":   "test_db_readonly",
				"username": "readonly_user",
				"password": "readonly_password",
			},
		},
	}
	
	bh.testData["redis"] = map[string]interface{}{
		"master": map[string]interface{}{
			"host": "localhost",
			"port": 6379,
		},
		"slave": map[string]interface{}{
			"host": "slave.localhost",
			"port": 6380,
		},
	}
	
	bh.testData["cassandra"] = map[string]interface{}{
		"writer": map[string]interface{}{
			"endpoints": []string{"localhost:9042", "backup.localhost:9042"},
			"keyspace":  "test_keyspace",
			"username":  "writer_user",
			"password":  "writer_password",
			"ca_path":   "/path/to/writer/ca.pem",
		},
		"reader": map[string]interface{}{
			"endpoints": []string{"readonly.localhost:9042"},
			"keyspace":  "test_keyspace",
			"username":  "reader_user",
			"password":  "reader_password",
			"ca_path":   "/path/to/reader/ca.pem",
		},
	}
}

// GetTestData returns test data for the specified type
func (bh *BenchmarkHelper) GetTestData(dataType string) interface{} {
	return bh.testData[dataType]
}
