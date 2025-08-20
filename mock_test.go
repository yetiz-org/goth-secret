package secret

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMockFileSystem(t *testing.T) {
	mfs := NewMockFileSystem()
	
	// Test adding files
	content := []byte("test content")
	mfs.AddFile("test.txt", content)
	
	// Test reading existing file
	readContent, err := mfs.ReadFile("test.txt")
	assert.NoError(t, err)
	assert.Equal(t, content, readContent)
	
	// Test reading non-existent file
	_, err = mfs.ReadFile("nonexistent.txt")
	assert.Error(t, err)
	assert.Equal(t, os.ErrNotExist, err)
	
	// Test adding file error
	testErr := fmt.Errorf("mock error")
	mfs.AddFileError("error.txt", testErr)
	
	_, err = mfs.ReadFile("error.txt")
	assert.Error(t, err)
	assert.Equal(t, testErr, err)
	
	// Test Stat functionality
	stat, err := mfs.Stat("test.txt")
	assert.NoError(t, err)
	assert.Equal(t, "test.txt", stat.Name())
	assert.Equal(t, int64(len(content)), stat.Size())
	assert.False(t, stat.IsDir())
	
	// Test Stat on non-existent file
	_, err = mfs.Stat("nonexistent.txt")
	assert.Error(t, err)
	assert.Equal(t, os.ErrNotExist, err)
}

func TestMockFileSystemAddSecretFile(t *testing.T) {
	mfs := NewMockFileSystem()
	
	testSecret := map[string]interface{}{
		"writer": map[string]interface{}{
			"adapter": "mysql",
			"params": map[string]interface{}{
				"host": "localhost",
				"port": float64(3306), // JSON unmarshaling converts numbers to float64
			},
		},
	}
	
	err := mfs.AddSecretFile("database", "test", testSecret)
	assert.NoError(t, err)
	
	// Verify the file was created with correct content
	content, err := mfs.ReadFile("database-test/secret.json")
	assert.NoError(t, err)
	
	var loadedSecret map[string]interface{}
	err = json.Unmarshal(content, &loadedSecret)
	assert.NoError(t, err)
	assert.Equal(t, testSecret, loadedSecret)
}

func TestMockEnvironment(t *testing.T) {
	env := NewMockEnvironment()
	
	// Test setting and getting environment variables
	env.SetVar("TEST_VAR", "test_value")
	value := env.Getenv("TEST_VAR")
	assert.Equal(t, "test_value", value)
	
	// Test getting non-existent variable
	value = env.Getenv("NONEXISTENT_VAR")
	assert.Equal(t, "", value)
	
	// Test GOTH_SECRET_PATH specifically
	env.SetVar("GOTH_SECRET_PATH", "/test/path")
	value = env.Getenv("GOTH_SECRET_PATH")
	assert.Equal(t, "/test/path", value)
}

func TestMockSecretLoader(t *testing.T) {
	loader := NewMockSecretLoader()
	
	// Test default behavior (no error)
	db := &Database{}
	err := loader.Load("database", "test", db)
	assert.NoError(t, err)
	
	// Verify call was recorded
	calls := loader.GetLoadCalls()
	assert.Len(t, calls, 1)
	assert.Equal(t, "database", calls[0].Type)
	assert.Equal(t, "test", calls[0].Name)
	assert.Equal(t, db, calls[0].Secret)
	
	// Test custom load function
	customErr := fmt.Errorf("custom load error")
	loader.SetLoadFunc(func(typ, name string, secret Secret) error {
		return customErr
	})
	
	redis := &Redis{}
	err = loader.Load("redis", "test", redis)
	assert.Error(t, err)
	assert.Equal(t, customErr, err)
	
	// Verify both calls were recorded
	calls = loader.GetLoadCalls()
	assert.Len(t, calls, 2)
	
	// Test clearing calls
	loader.ClearLoadCalls()
	calls = loader.GetLoadCalls()
	assert.Len(t, calls, 0)
}

func TestRealFileSystemInterfaces(t *testing.T) {
	// Test that real implementations satisfy the interfaces
	var fs FileSystemInterface = &RealFileSystem{}
	var env EnvironmentInterface = &RealEnvironment{}
	
	// Basic interface checks
	assert.NotNil(t, fs)
	assert.NotNil(t, env)
	
	// Test environment functionality
	// Note: We can't test actual file operations in unit tests
	// as they depend on the file system state
	testValue := env.Getenv("PATH") // PATH should exist on most systems
	// We just verify it doesn't panic and returns a string
	assert.IsType(t, "", testValue)
}

func TestMockFileInfoInterface(t *testing.T) {
	info := &mockFileInfo{
		name:    "test.txt",
		size:    100,
		mode:    0644,
		modTime: time.Now(),
		isDir:   false,
	}
	
	assert.Equal(t, "test.txt", info.Name())
	assert.Equal(t, int64(100), info.Size())
	assert.Equal(t, os.FileMode(0644), info.Mode())
	assert.False(t, info.IsDir())
	assert.Nil(t, info.Sys())
	assert.NotZero(t, info.ModTime()) // Should have a valid time
}

func TestMockIntegration(t *testing.T) {
	// Test integration between different mock components
	helper := NewTestHelper()
	
	// Setup mock environment
	helper.SetMockPath("/mock/path")
	
	// Create mock secret file
	testData := map[string]interface{}{
		"master": map[string]interface{}{
			"host": "mock.redis.com",
			"port": float64(6379), // JSON unmarshaling converts numbers to float64
		},
		"slave": map[string]interface{}{
			"host": "mock.slave.com",
			"port": float64(6380), // JSON unmarshaling converts numbers to float64
		},
	}
	
	err := helper.CreateMockSecretFile("redis", "integration", testData)
	assert.NoError(t, err)
	
	// Verify file was created in mock file system
	content, err := helper.GetMockFileSystem().ReadFile("redis-integration/secret.json")
	assert.NoError(t, err)
	
	var loaded map[string]interface{}
	err = json.Unmarshal(content, &loaded)
	assert.NoError(t, err)
	assert.Equal(t, testData, loaded)
	
	// Verify environment was set
	path := helper.GetMockEnvironment().Getenv("GOTH_SECRET_PATH")
	assert.Equal(t, "/mock/path", path)
}

func TestMockErrorScenarios(t *testing.T) {
	testCases := []struct {
		name          string
		setupFunc     func(*MockFileSystem)
		filename      string
		expectedError string
	}{
		{
			name: "File not found",
			setupFunc: func(mfs *MockFileSystem) {
				// Don't add any files
			},
			filename:      "nonexistent.json",
			expectedError: "file does not exist",
		},
		{
			name: "Read permission error",
			setupFunc: func(mfs *MockFileSystem) {
				mfs.AddFileError("permission.json", fmt.Errorf("permission denied"))
			},
			filename:      "permission.json",
			expectedError: "permission denied",
		},
		{
			name: "Network error simulation",
			setupFunc: func(mfs *MockFileSystem) {
				mfs.AddFileError("network.json", fmt.Errorf("network is unreachable"))
			},
			filename:      "network.json",
			expectedError: "network is unreachable",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mfs := NewMockFileSystem()
			tc.setupFunc(mfs)
			
			_, err := mfs.ReadFile(tc.filename)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedError)
			
			_, err = mfs.Stat(tc.filename)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedError)
		})
	}
}

func TestMockSecretFileFormats(t *testing.T) {
	testCases := []struct {
		name    string
		content interface{}
		isValid bool
	}{
		{
			name: "Valid Database JSON",
			content: map[string]interface{}{
				"writer": map[string]interface{}{
					"adapter": "mysql",
					"params": map[string]interface{}{
						"host": "localhost",
						"port": 3306,
					},
				},
			},
			isValid: true,
		},
		{
			name: "Valid Redis JSON",
			content: map[string]interface{}{
				"master": map[string]interface{}{
					"host": "localhost",
					"port": 6379,
				},
			},
			isValid: true,
		},
		{
			name: "Empty JSON object",
			content: map[string]interface{}{},
			isValid: true,
		},
		{
			name: "Complex nested structure",
			content: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"level3": []string{"item1", "item2"},
					},
				},
			},
			isValid: true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mfs := NewMockFileSystem()
			
			err := mfs.AddSecretFile("test", "case", tc.content)
			if tc.isValid {
				assert.NoError(t, err)
				
				// Verify file was created and can be read
				content, err := mfs.ReadFile("test-case/secret.json")
				assert.NoError(t, err)
				
				// Verify JSON is valid
				var loaded interface{}
				err = json.Unmarshal(content, &loaded)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestMockBuilderIntegration(t *testing.T) {
	// Test that mock builders work with mock file system
	helper := NewTestHelper()
	
	// Create mocks using builders
	db := NewDatabaseMock().WithName("test").WithDefaultMySQL().Build()
	redis := NewRedisMock().WithName("test").WithDefault().Build()
	cassandra := NewCassandraMock().WithName("test").WithDefault().Build()
	
	// Add them to mock file system
	err := helper.CreateMockSecretFile("database", "test", db)
	assert.NoError(t, err)
	
	err = helper.CreateMockSecretFile("redis", "test", redis)
	assert.NoError(t, err)
	
	err = helper.CreateMockSecretFile("cassandra", "test", cassandra)
	assert.NoError(t, err)
	
	// Verify they can be read back
	dbContent, err := helper.GetMockFileSystem().ReadFile("database-test/secret.json")
	assert.NoError(t, err)
	assert.NotEmpty(t, dbContent)
	
	redisContent, err := helper.GetMockFileSystem().ReadFile("redis-test/secret.json")
	assert.NoError(t, err)
	assert.NotEmpty(t, redisContent)
	
	cassandraContent, err := helper.GetMockFileSystem().ReadFile("cassandra-test/secret.json")
	assert.NoError(t, err)
	assert.NotEmpty(t, cassandraContent)
}
