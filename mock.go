package secret

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// FileSystemInterface abstracts file system operations for testing
type FileSystemInterface interface {
	ReadFile(filename string) ([]byte, error)
	Stat(name string) (os.FileInfo, error)
}

// RealFileSystem provides the actual file system implementation
type RealFileSystem struct{}

func (fs *RealFileSystem) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (fs *RealFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// MockFileSystem provides an in-memory file system for testing
type MockFileSystem struct {
	files  map[string][]byte
	errors map[string]error
	stats  map[string]os.FileInfo
}

// NewMockFileSystem creates a new mock file system
func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		files:  make(map[string][]byte),
		errors: make(map[string]error),
		stats:  make(map[string]os.FileInfo),
	}
}

// AddFile adds a file to the mock file system
func (mfs *MockFileSystem) AddFile(path string, content []byte) {
	mfs.files[path] = content
	mfs.stats[path] = &mockFileInfo{
		name:    filepath.Base(path),
		size:    int64(len(content)),
		mode:    0644,
		modTime: time.Now(),
		isDir:   false,
	}
}

// AddFileError sets an error to be returned when accessing a specific file
func (mfs *MockFileSystem) AddFileError(path string, err error) {
	mfs.errors[path] = err
}

// AddSecretFile adds a secret file with JSON content
func (mfs *MockFileSystem) AddSecretFile(typ, name string, content interface{}) error {
	jsonData, err := json.Marshal(content)
	if err != nil {
		return err
	}
	
	secretPath := fmt.Sprintf("%s-%s/secret.json", typ, name)
	mfs.AddFile(secretPath, jsonData)
	return nil
}

func (mfs *MockFileSystem) ReadFile(filename string) ([]byte, error) {
	if err, exists := mfs.errors[filename]; exists {
		return nil, err
	}
	
	if content, exists := mfs.files[filename]; exists {
		return content, nil
	}
	
	return nil, os.ErrNotExist
}

func (mfs *MockFileSystem) Stat(name string) (os.FileInfo, error) {
	if err, exists := mfs.errors[name]; exists {
		return nil, err
	}
	
	if stat, exists := mfs.stats[name]; exists {
		return stat, nil
	}
	
	return nil, os.ErrNotExist
}

// mockFileInfo implements os.FileInfo for testing
type mockFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (fi *mockFileInfo) Name() string       { return fi.name }
func (fi *mockFileInfo) Size() int64        { return fi.size }
func (fi *mockFileInfo) Mode() os.FileMode  { return fi.mode }
func (fi *mockFileInfo) ModTime() time.Time { return fi.modTime }
func (fi *mockFileInfo) IsDir() bool        { return fi.isDir }
func (fi *mockFileInfo) Sys() interface{}   { return nil }

// EnvironmentInterface abstracts environment variable access
type EnvironmentInterface interface {
	Getenv(key string) string
}

// RealEnvironment provides actual environment variable access
type RealEnvironment struct{}

func (env *RealEnvironment) Getenv(key string) string {
	return os.Getenv(key)
}

// MockEnvironment provides controllable environment variables for testing
type MockEnvironment struct {
	vars map[string]string
}

// NewMockEnvironment creates a new mock environment
func NewMockEnvironment() *MockEnvironment {
	return &MockEnvironment{
		vars: make(map[string]string),
	}
}

// SetVar sets an environment variable in the mock
func (env *MockEnvironment) SetVar(key, value string) {
	env.vars[key] = value
}

func (env *MockEnvironment) Getenv(key string) string {
	if value, exists := env.vars[key]; exists {
		return value
	}
	return ""
}

// SecretLoaderInterface abstracts the secret loading functionality
type SecretLoaderInterface interface {
	Load(typ string, name string, secret Secret) error
}

// MockSecretLoader provides controllable secret loading for testing
type MockSecretLoader struct {
	loadFunc func(typ string, name string, secret Secret) error
	loadCalls []LoadCall
}

// LoadCall records parameters of Load function calls for verification
type LoadCall struct {
	Type   string
	Name   string
	Secret Secret
}

// NewMockSecretLoader creates a new mock secret loader
func NewMockSecretLoader() *MockSecretLoader {
	return &MockSecretLoader{
		loadCalls: make([]LoadCall, 0),
	}
}

// SetLoadFunc sets the function to be called when Load is invoked
func (msl *MockSecretLoader) SetLoadFunc(f func(typ string, name string, secret Secret) error) {
	msl.loadFunc = f
}

// Load implements SecretLoaderInterface
func (msl *MockSecretLoader) Load(typ string, name string, secret Secret) error {
	msl.loadCalls = append(msl.loadCalls, LoadCall{
		Type:   typ,
		Name:   name,
		Secret: secret,
	})
	
	if msl.loadFunc != nil {
		return msl.loadFunc(typ, name, secret)
	}
	return nil
}

// GetLoadCalls returns all recorded Load function calls
func (msl *MockSecretLoader) GetLoadCalls() []LoadCall {
	return msl.loadCalls
}

// ClearLoadCalls clears the recorded Load function calls
func (msl *MockSecretLoader) ClearLoadCalls() {
	msl.loadCalls = msl.loadCalls[:0]
}
