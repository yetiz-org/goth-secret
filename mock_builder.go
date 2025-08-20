package secret

// DatabaseMockBuilder provides a fluent interface for creating Database mocks
type DatabaseMockBuilder struct {
	database *Database
}

// NewDatabaseMock creates a new Database mock builder
func NewDatabaseMock() *DatabaseMockBuilder {
	return &DatabaseMockBuilder{
		database: &Database{},
	}
}

// WithName sets the name for the Database mock
func (b *DatabaseMockBuilder) WithName(name string) *DatabaseMockBuilder {
	b.database.DefaultSecret._Name = name
	return b
}

// WithPath sets the path for the Database mock
func (b *DatabaseMockBuilder) WithPath(path string) *DatabaseMockBuilder {
	b.database.DefaultSecret._Path = path
	return b
}

// WithWriter configures the writer database connection
func (b *DatabaseMockBuilder) WithWriter(adapter, host string, port uint, dbname, username, password string) *DatabaseMockBuilder {
	b.database.Writer.Adapter = adapter
	b.database.Writer.Params.Host = host
	b.database.Writer.Params.Port = port
	b.database.Writer.Params.DBName = dbname
	b.database.Writer.Params.Username = username
	b.database.Writer.Params.Password = password
	return b
}

// WithWriterCharset sets the charset for the writer database connection
func (b *DatabaseMockBuilder) WithWriterCharset(charset string) *DatabaseMockBuilder {
	b.database.Writer.Params.Charset = charset
	return b
}

// WithReader configures the reader database connection
func (b *DatabaseMockBuilder) WithReader(adapter, host string, port uint, dbname, username, password string) *DatabaseMockBuilder {
	b.database.Reader.Adapter = adapter
	b.database.Reader.Params.Host = host
	b.database.Reader.Params.Port = port
	b.database.Reader.Params.DBName = dbname
	b.database.Reader.Params.Username = username
	b.database.Reader.Params.Password = password
	return b
}

// WithReaderCharset sets the charset for the reader database connection
func (b *DatabaseMockBuilder) WithReaderCharset(charset string) *DatabaseMockBuilder {
	b.database.Reader.Params.Charset = charset
	return b
}

// WithDefaultMySQL sets up default MySQL configuration for both writer and reader
func (b *DatabaseMockBuilder) WithDefaultMySQL() *DatabaseMockBuilder {
	return b.WithWriter("mysql", "localhost", 3306, "test_db", "test_user", "test_password").
		WithWriterCharset("utf8mb4").
		WithReader("mysql", "readonly.localhost", 3306, "test_db_readonly", "readonly_user", "readonly_password").
		WithReaderCharset("utf8mb4")
}

// Build returns the configured Database mock
func (b *DatabaseMockBuilder) Build() *Database {
	return b.database
}

// RedisMockBuilder provides a fluent interface for creating Redis mocks
type RedisMockBuilder struct {
	redis *Redis
}

// NewRedisMock creates a new Redis mock builder
func NewRedisMock() *RedisMockBuilder {
	return &RedisMockBuilder{
		redis: &Redis{},
	}
}

// WithName sets the name for the Redis mock
func (b *RedisMockBuilder) WithName(name string) *RedisMockBuilder {
	b.redis.DefaultSecret._Name = name
	return b
}

// WithPath sets the path for the Redis mock
func (b *RedisMockBuilder) WithPath(path string) *RedisMockBuilder {
	b.redis.DefaultSecret._Path = path
	return b
}

// WithMaster configures the master Redis connection
func (b *RedisMockBuilder) WithMaster(host string, port uint) *RedisMockBuilder {
	b.redis.Master.Host = host
	b.redis.Master.Port = port
	return b
}

// WithSlave configures the slave Redis connection
func (b *RedisMockBuilder) WithSlave(host string, port uint) *RedisMockBuilder {
	b.redis.Slave.Host = host
	b.redis.Slave.Port = port
	return b
}

// WithDefault sets up default Redis configuration
func (b *RedisMockBuilder) WithDefault() *RedisMockBuilder {
	return b.WithMaster("localhost", 6379).
		WithSlave("slave.localhost", 6380)
}

// Build returns the configured Redis mock
func (b *RedisMockBuilder) Build() *Redis {
	return b.redis
}

// CassandraMockBuilder provides a fluent interface for creating Cassandra mocks
type CassandraMockBuilder struct {
	cassandra *Cassandra
}

// NewCassandraMock creates a new Cassandra mock builder
func NewCassandraMock() *CassandraMockBuilder {
	return &CassandraMockBuilder{
		cassandra: &Cassandra{},
	}
}

// WithName sets the name for the Cassandra mock
func (b *CassandraMockBuilder) WithName(name string) *CassandraMockBuilder {
	b.cassandra.DefaultSecret._Name = name
	return b
}

// WithPath sets the path for the Cassandra mock
func (b *CassandraMockBuilder) WithPath(path string) *CassandraMockBuilder {
	b.cassandra.DefaultSecret._Path = path
	return b
}

// WithWriter configures the writer Cassandra connection
func (b *CassandraMockBuilder) WithWriter(endpoints []string, keyspace, username, password, caPath string) *CassandraMockBuilder {
	b.cassandra.Writer.Endpoints = endpoints
	b.cassandra.Writer.Keyspace = keyspace
	b.cassandra.Writer.Username = username
	b.cassandra.Writer.Password = password
	b.cassandra.Writer.CaPath = caPath
	return b
}

// WithReader configures the reader Cassandra connection
func (b *CassandraMockBuilder) WithReader(endpoints []string, keyspace, username, password, caPath string) *CassandraMockBuilder {
	b.cassandra.Reader.Endpoints = endpoints
	b.cassandra.Reader.Keyspace = keyspace
	b.cassandra.Reader.Username = username
	b.cassandra.Reader.Password = password
	b.cassandra.Reader.CaPath = caPath
	return b
}

// WithDefault sets up default Cassandra configuration
func (b *CassandraMockBuilder) WithDefault() *CassandraMockBuilder {
	return b.WithWriter(
		[]string{"localhost:9042", "backup.localhost:9042"},
		"test_keyspace",
		"writer_user",
		"writer_password",
		"/path/to/writer/ca.pem",
	).WithReader(
		[]string{"readonly.localhost:9042"},
		"test_keyspace",
		"reader_user",
		"reader_password",
		"/path/to/reader/ca.pem",
	)
}

// Build returns the configured Cassandra mock
func (b *CassandraMockBuilder) Build() *Cassandra {
	return b.cassandra
}

// MockSecretBuilder provides a unified interface for creating any type of secret mock
type MockSecretBuilder struct {
	secretType string
	name       string
	path       string
}

// NewMockSecret creates a new mock secret builder
func NewMockSecret(secretType, name string) *MockSecretBuilder {
	return &MockSecretBuilder{
		secretType: secretType,
		name:       name,
		path:       "",
	}
}

// WithPath sets a custom path for the mock secret
func (b *MockSecretBuilder) WithPath(path string) *MockSecretBuilder {
	b.path = path
	return b
}

// BuildDatabase creates a Database mock with the specified configuration
func (b *MockSecretBuilder) BuildDatabase() *Database {
	builder := NewDatabaseMock().WithName(b.name).WithDefaultMySQL()
	if b.path != "" {
		builder = builder.WithPath(b.path)
	}
	return builder.Build()
}

// BuildRedis creates a Redis mock with the specified configuration
func (b *MockSecretBuilder) BuildRedis() *Redis {
	builder := NewRedisMock().WithName(b.name).WithDefault()
	if b.path != "" {
		builder = builder.WithPath(b.path)
	}
	return builder.Build()
}

// BuildCassandra creates a Cassandra mock with the specified configuration
func (b *MockSecretBuilder) BuildCassandra() *Cassandra {
	builder := NewCassandraMock().WithName(b.name).WithDefault()
	if b.path != "" {
		builder = builder.WithPath(b.path)
	}
	return builder.Build()
}

// TestScenarioBuilder helps create complex test scenarios with multiple secrets
type TestScenarioBuilder struct {
	helper    *TestHelper
	scenarios map[string]Secret
}

// NewTestScenario creates a new test scenario builder
func NewTestScenario() *TestScenarioBuilder {
	return &TestScenarioBuilder{
		helper:    NewTestHelper(),
		scenarios: make(map[string]Secret),
	}
}

// AddDatabase adds a database secret to the test scenario
func (b *TestScenarioBuilder) AddDatabase(name string, config *Database) *TestScenarioBuilder {
	b.scenarios[name] = config
	b.helper.CreateMockSecretFile("database", name, config)
	return b
}

// AddRedis adds a redis secret to the test scenario
func (b *TestScenarioBuilder) AddRedis(name string, config *Redis) *TestScenarioBuilder {
	b.scenarios[name] = config
	b.helper.CreateMockSecretFile("redis", name, config)
	return b
}

// AddCassandra adds a cassandra secret to the test scenario
func (b *TestScenarioBuilder) AddCassandra(name string, config *Cassandra) *TestScenarioBuilder {
	b.scenarios[name] = config
	b.helper.CreateMockSecretFile("cassandra", name, config)
	return b
}

// GetHelper returns the test helper for advanced usage
func (b *TestScenarioBuilder) GetHelper() *TestHelper {
	return b.helper
}

// GetScenarios returns all configured scenarios
func (b *TestScenarioBuilder) GetScenarios() map[string]Secret {
	return b.scenarios
}
