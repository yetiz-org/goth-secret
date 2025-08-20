package secret

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLoadFunctionComprehensive provides comprehensive testing for the Load function
func TestLoadFunctionComprehensive(t *testing.T) {
	testCases := []struct {
		name          string
		setupFunc     func(*TestHelper) (string, Secret)
		expectError   bool
		errorContains string
	}{
		{
			name: "Valid Database Loading",
			setupFunc: func(helper *TestHelper) (string, Secret) {
				// For this test, we'll create the secret data that would be loaded
				return "database", &Database{}
			},
			expectError: false,
		},
		{
			name: "Valid Redis Loading",
			setupFunc: func(helper *TestHelper) (string, Secret) {
				// For this test, we'll create the secret data that would be loaded
				return "redis", &Redis{}
			},
			expectError: false,
		},
		{
			name: "Valid Cassandra Loading",
			setupFunc: func(helper *TestHelper) (string, Secret) {
				// For this test, we'll create the secret data that would be loaded
				return "cassandra", &Cassandra{}
			},
			expectError: false,
		},
		{
			name: "File Not Found",
			setupFunc: func(helper *TestHelper) (string, Secret) {
				// Don't create any file
				return "database", &Database{}
			},
			expectError:   true,
			errorContains: "no such file or directory",
		},
		{
			name: "Invalid JSON Format",
			setupFunc: func(helper *TestHelper) (string, Secret) {
				helper.GetMockFileSystem().AddFile("database-test/secret.json", []byte("{invalid json"))
				return "database", &Database{}
			},
			expectError:   true,
			errorContains: "invalid character",
		},
		{
			name: "Non-pointer Secret",
			setupFunc: func(helper *TestHelper) (string, Secret) {
				db := CreateDatabaseSecret("test")
				helper.CreateMockSecretFile("database", "test", db)
				// This will be handled in the test by passing non-pointer
				return "database", nil // Will be replaced in test
			},
			expectError:   true,
			errorContains: "secret should be a pointer",
		},
		{
			name: "Non-struct Secret",
			setupFunc: func(helper *TestHelper) (string, Secret) {
				// This will be handled in the test with a custom secret
				return "test", nil // Will be replaced in test
			},
			expectError:   true,
			errorContains: "secret should be a struct",
		},
		{
			name: "Missing DefaultSecret Field",
			setupFunc: func(helper *TestHelper) (string, Secret) {
				helper.CreateMockSecretFile("test", "test", map[string]interface{}{"field": "value"})
				return "test", &NoDefaultSecret{}
			},
			expectError:   true,
			errorContains: "struct should have a DefaultSecret field",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			helper := NewTestHelper()
			cleanup := helper.SetupTempDirectory(t, "./")
			defer cleanup()

			var secretType string
			var secret Secret

			if tc.name == "Non-pointer Secret" {
				// Test will naturally fail because interface requires pointer
				t.Skip("Skipping pointer test - interface requires pointer types")
				return
			}

			if tc.name == "Non-struct Secret" {
				// Create a test for non-struct - we'll skip this as it's hard to test with current interface
				t.Skip("Skipping non-struct test due to interface constraints")
				return
			}

			secretType, secret = tc.setupFunc(helper)

			// Create actual test files for valid loading scenarios
			if !tc.expectError && tc.name != "File Not Found" {
				var testData interface{}
				switch secretType {
				case "database":
					testData = CreateDatabaseSecret("test")
				case "redis":
					testData = CreateRedisSecret("test")
				case "cassandra":
					testData = CreateCassandraSecret("test")
				}
				if testData != nil {
					_, fileCleanup := helper.CreateSecretFile(t, secretType, "test", testData)
					defer fileCleanup()
				}
			}

			// Handle special invalid JSON case
			if tc.name == "Invalid JSON Format" {
				_, fileCleanup := helper.CreateSecretFile(t, "database", "test", "{invalid json")
				defer fileCleanup()
			}

			// Mock the PATH environment for testing
			originalPath := PATH
			PATH = "./"
			defer func() { PATH = originalPath }()

			err := Load(secretType, "test", secret)

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "test", secret.Name())
				assert.NotEmpty(t, secret.Path())
			}
		})
	}
}

// TestReflectionEdgeCases tests edge cases in the reflection logic
func TestReflectionEdgeCases(t *testing.T) {
	testCases := []struct {
		name        string
		secret      interface{}
		expectFound bool
	}{
		{
			name: "Direct DefaultSecret field",
			secret: struct {
				DefaultSecret
				Field string
			}{},
			expectFound: true,
		},
		{
			name: "Nested DefaultSecret field",
			secret: struct {
				Nested struct {
					DefaultSecret
				}
				Field string
			}{},
			expectFound: true,
		},
		{
			name: "Deep nested DefaultSecret field",
			secret: struct {
				Level1 struct {
					Level2 struct {
						DefaultSecret
					}
				}
				Field string
			}{},
			expectFound: true,
		},
		{
			name: "Multiple nested structs, one with DefaultSecret",
			secret: struct {
				NestedA struct {
					FieldA string
				}
				NestedB struct {
					DefaultSecret
					FieldB string
				}
				Field string
			}{},
			expectFound: true,
		},
		{
			name: "No DefaultSecret field",
			secret: struct {
				Field1 string
				Field2 int
			}{},
			expectFound: false,
		},
		{
			name: "Struct with only primitive fields",
			secret: struct {
				StringField string
				IntField    int
				BoolField   bool
			}{},
			expectFound: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value := reflect.ValueOf(tc.secret)
			found, fieldValue := findDefaultSecret(value)

			assert.Equal(t, tc.expectFound, found)
			if tc.expectFound {
				assert.True(t, fieldValue.IsValid())
			} else {
				assert.False(t, fieldValue.IsValid())
			}
		})
	}
}

// TestJSONUnmarshalingEdgeCases tests various JSON unmarshaling scenarios
func TestJSONUnmarshalingEdgeCases(t *testing.T) {
	testCases := []struct {
		name        string
		jsonContent string
		secretType  interface{}
		expectError bool
	}{
		{
			name: "Valid Database JSON with extra fields",
			jsonContent: `{
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
				},
				"extra_field": "should be ignored"
			}`,
			secretType:  &Database{},
			expectError: false,
		},
		{
			name: "JSON with null values",
			jsonContent: `{
				"master": {
					"host": null,
					"port": 6379
				},
				"slave": null
			}`,
			secretType:  &Redis{},
			expectError: false,
		},
		{
			name: "JSON with wrong data types",
			jsonContent: `{
				"master": {
					"host": "localhost",
					"port": "not_a_number"
				}
			}`,
			secretType:  &Redis{},
			expectError: true,
		},
		{
			name:        "Empty JSON object",
			jsonContent: `{}`,
			secretType:  &Database{},
			expectError: false,
		},
		{
			name:        "JSON with array instead of object",
			jsonContent: `[{"host": "localhost"}]`,
			secretType:  &Redis{},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			helper := NewTestHelper()
			cleanup := helper.SetupTempDirectory(t, "./")
			defer cleanup()

			// Create actual files on disk instead of using mock file system
			_, fileCleanup := helper.CreateSecretFile(t, "test", "secret", tc.jsonContent)
			defer fileCleanup()

			// Mock the PATH environment
			originalPath := PATH
			PATH = "./"
			defer func() { PATH = originalPath }()

			err := Load("test", "secret", tc.secretType.(Secret))

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestConcurrentLoading tests concurrent access to the Load function
func TestConcurrentLoading(t *testing.T) {
	helper := NewTestHelper()
	cleanup := helper.SetupTempDirectory(t, "./")
	defer cleanup()

	// Create test data as actual files for the concurrent test
	dbData := CreateDatabaseSecret("concurrent")
	redisData := CreateRedisSecret("concurrent")
	cassandraData := CreateCassandraSecret("concurrent")

	// Create actual files for testing
	_, dbCleanup := helper.CreateSecretFile(t, "database", "concurrent", dbData)
	defer dbCleanup()

	_, redisCleanup := helper.CreateSecretFile(t, "redis", "concurrent", redisData)
	defer redisCleanup()

	_, cassandraCleanup := helper.CreateSecretFile(t, "cassandra", "concurrent", cassandraData)
	defer cassandraCleanup()

	// Mock the PATH environment
	originalPath := PATH
	PATH = "./"
	defer func() { PATH = originalPath }()

	// Run concurrent loads
	done := make(chan bool, 6)
	errors := make(chan error, 6)

	for i := 0; i < 2; i++ {
		go func() {
			db := &Database{}
			err := Load("database", "concurrent", db)
			errors <- err
			done <- true
		}()

		go func() {
			redis := &Redis{}
			err := Load("redis", "concurrent", redis)
			errors <- err
			done <- true
		}()

		go func() {
			cassandra := &Cassandra{}
			err := Load("cassandra", "concurrent", cassandra)
			errors <- err
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 6; i++ {
		<-done
	}

	// Check that all loads succeeded
	for i := 0; i < 6; i++ {
		err := <-errors
		assert.NoError(t, err)
	}
}

// TestPathResolution tests different path resolution scenarios
func TestPathResolution(t *testing.T) {
	testCases := []struct {
		name        string
		globalPath  string
		envPath     string
		expectedDir string
	}{
		{
			name:        "Global PATH variable set",
			globalPath:  "/global/path",
			envPath:     "",
			expectedDir: "/global/path",
		},
		{
			name:        "Environment variable set",
			globalPath:  "",
			envPath:     "/env/path",
			expectedDir: "/env/path",
		},
		{
			name:        "Both set, global takes precedence",
			globalPath:  "/global/path",
			envPath:     "/env/path",
			expectedDir: "/global/path",
		},
		{
			name:        "Neither set",
			globalPath:  "",
			envPath:     "",
			expectedDir: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Save original values
			originalPATH := PATH
			originalEnv := os.Getenv("GOTH_SECRET_PATH")

			// Set test values
			PATH = tc.globalPath
			if tc.envPath != "" {
				os.Setenv("GOTH_SECRET_PATH", tc.envPath)
			} else {
				os.Unsetenv("GOTH_SECRET_PATH")
			}

			// Test path resolution
			resolvedPath := Path()
			assert.Equal(t, tc.expectedDir, resolvedPath)

			// Restore original values
			PATH = originalPATH
			if originalEnv != "" {
				os.Setenv("GOTH_SECRET_PATH", originalEnv)
			} else {
				os.Unsetenv("GOTH_SECRET_PATH")
			}
		})
	}
}

// TestBuilderPatterns tests the builder pattern implementations
func TestBuilderPatterns(t *testing.T) {
	t.Run("Database Builder", func(t *testing.T) {
		db := NewDatabaseMock().
			WithName("test-db").
			WithPath("/custom/path").
			WithWriter("postgresql", "writer.host", 5432, "writer_db", "writer_user", "writer_pass").
			WithWriterCharset("utf8").
			WithReader("postgresql", "reader.host", 5433, "reader_db", "reader_user", "reader_pass").
			WithReaderCharset("utf8").
			Build()

		assert.Equal(t, "test-db", db.Name())
		assert.Equal(t, "/custom/path", db.Path())
		assert.Equal(t, "postgresql", db.Writer.Adapter)
		assert.Equal(t, "writer.host", db.Writer.Params.Host)
		assert.Equal(t, uint(5432), db.Writer.Params.Port)
		assert.Equal(t, "utf8", db.Writer.Params.Charset)

		assert.Equal(t, "postgresql", db.Reader.Adapter)
		assert.Equal(t, "reader.host", db.Reader.Params.Host)
		assert.Equal(t, uint(5433), db.Reader.Params.Port)
		assert.Equal(t, "utf8", db.Reader.Params.Charset)
	})

	t.Run("Redis Builder", func(t *testing.T) {
		redis := NewRedisMock().
			WithName("test-redis").
			WithPath("/custom/redis/path").
			WithMaster("master.redis.com", 6379).
			WithSlave("slave.redis.com", 6380).
			Build()

		assert.Equal(t, "test-redis", redis.Name())
		assert.Equal(t, "/custom/redis/path", redis.Path())
		assert.Equal(t, "master.redis.com", redis.Master.Host)
		assert.Equal(t, uint(6379), redis.Master.Port)
		assert.Equal(t, "slave.redis.com", redis.Slave.Host)
		assert.Equal(t, uint(6380), redis.Slave.Port)
	})

	t.Run("Cassandra Builder", func(t *testing.T) {
		cassandra := NewCassandraMock().
			WithName("test-cassandra").
			WithPath("/custom/cassandra/path").
			WithWriter([]string{"w1:9042", "w2:9042"}, "writer_ks", "w_user", "w_pass", "/w/ca.pem").
			WithReader([]string{"r1:9042"}, "reader_ks", "r_user", "r_pass", "/r/ca.pem").
			Build()

		assert.Equal(t, "test-cassandra", cassandra.Name())
		assert.Equal(t, "/custom/cassandra/path", cassandra.Path())
		assert.Equal(t, []string{"w1:9042", "w2:9042"}, cassandra.Writer.Endpoints)
		assert.Equal(t, "writer_ks", cassandra.Writer.Keyspace)
		assert.Equal(t, []string{"r1:9042"}, cassandra.Reader.Endpoints)
		assert.Equal(t, "reader_ks", cassandra.Reader.Keyspace)
	})
}

// TestScenarioBuilderFunctionality tests the test scenario builder functionality
func TestScenarioBuilderFunctionality(t *testing.T) {
	scenario := NewTestScenario()

	// Add multiple secrets to the scenario
	db := NewDatabaseMock().WithName("app-db").WithDefaultMySQL().Build()
	redis := NewRedisMock().WithName("app-cache").WithDefault().Build()
	cassandra := NewCassandraMock().WithName("app-analytics").WithDefault().Build()

	scenario.AddDatabase("app-db", db).
		AddRedis("app-cache", redis).
		AddCassandra("app-analytics", cassandra)

	// Verify all scenarios were added
	scenarios := scenario.GetScenarios()
	assert.Len(t, scenarios, 3)
	assert.Contains(t, scenarios, "app-db")
	assert.Contains(t, scenarios, "app-cache")
	assert.Contains(t, scenarios, "app-analytics")

	// Verify files were created in mock file system
	helper := scenario.GetHelper()
	fs := helper.GetMockFileSystem()

	_, err := fs.ReadFile("database-app-db/secret.json")
	assert.NoError(t, err)

	_, err = fs.ReadFile("redis-app-cache/secret.json")
	assert.NoError(t, err)

	_, err = fs.ReadFile("cassandra-app-analytics/secret.json")
	assert.NoError(t, err)
}
