// Copyright 2023 N4-Networks.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the main application configuration
type Config struct {
	Service    ServiceConfig    `yaml:"service"`
	Database   DatabaseConfig   `yaml:"database"`
	MessageBus MessageBusConfig `yaml:"messageBus"`
	Protocols  ProtocolsConfig  `yaml:"protocols"`
	Security   SecurityConfig   `yaml:"security"`
	Logging    LoggingConfig    `yaml:"logging"`
}

// ServiceConfig contains service-specific configuration
type ServiceConfig struct {
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version"`
	Environment string            `yaml:"environment"`
	Debug       bool              `yaml:"debug"`
	Metadata    map[string]string `yaml:"metadata,omitempty"`
}

// DatabaseConfig contains database connection settings
type DatabaseConfig struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	URI      string `yaml:"uri,omitempty"`
	Pool     struct {
		MaxConnections int           `yaml:"maxConnections"`
		Timeout        time.Duration `yaml:"timeout"`
	} `yaml:"pool"`
}

// MessageBusConfig contains message bus configuration
type MessageBusConfig struct {
	STOMP StompConfig `yaml:"stomp"`
	MQTT  MqttConfig  `yaml:"mqtt"`
	COAP  CoapConfig  `yaml:"coap"`
}

// StompConfig contains STOMP protocol configuration
type StompConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	TLSPort     int    `yaml:"tlsPort"`
	Mode        string `yaml:"mode"`
	Username    string `yaml:"username,omitempty"`
	Password    string `yaml:"password,omitempty"`
	Queue       string `yaml:"queue"`
	ConnRetry   int    `yaml:"connRetry"`
	EnableTLS   bool   `yaml:"enableTLS"`
	CertFile    string `yaml:"certFile,omitempty"`
	KeyFile     string `yaml:"keyFile,omitempty"`
	CACertFile  string `yaml:"caCertFile,omitempty"`
}

// MqttConfig contains MQTT protocol configuration
type MqttConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Mode      string `yaml:"mode"`
	Username  string `yaml:"username,omitempty"`
	Password  string `yaml:"password,omitempty"`
	Topic     string `yaml:"topic"`
	ClientID  string `yaml:"clientId"`
	EnableTLS bool   `yaml:"enableTLS"`
}

// CoapConfig contains CoAP protocol configuration
type CoapConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DTLSPort int    `yaml:"dtlsPort"`
	Mode     string `yaml:"mode"`
}

// ProtocolsConfig contains protocol-specific settings
type ProtocolsConfig struct {
	HTTP      HTTPConfig      `yaml:"http"`
	GRPC      GRPCConfig      `yaml:"grpc"`
	WebSocket WebSocketConfig `yaml:"websocket"`
	CWMP      CWMPConfig      `yaml:"cwmp"`
}

// HTTPConfig contains HTTP server configuration
type HTTPConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	TLSPort   int    `yaml:"tlsPort"`
	EnableTLS bool   `yaml:"enableTLS"`
	CertFile  string `yaml:"certFile,omitempty"`
	KeyFile   string `yaml:"keyFile,omitempty"`
}

// GRPCConfig contains gRPC server configuration
type GRPCConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	EnableTLS bool   `yaml:"enableTLS"`
	CertFile  string `yaml:"certFile,omitempty"`
	KeyFile   string `yaml:"keyFile,omitempty"`
}

// WebSocketConfig contains WebSocket configuration
type WebSocketConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	TLSPort   int    `yaml:"tlsPort"`
	Path      string `yaml:"path"`
	Mode      string `yaml:"mode"`
	EnableTLS bool   `yaml:"enableTLS"`
}

// CWMPConfig contains CWMP/TR-069 configuration
type CWMPConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	TLSPort  int    `yaml:"tlsPort"`
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// SecurityConfig contains security-related configuration
type SecurityConfig struct {
	Auth  AuthConfig  `yaml:"auth"`
	TLS   TLSConfig   `yaml:"tls"`
	USP   USPConfig   `yaml:"usp"`
	Cache CacheConfig `yaml:"cache"`
}

// AuthConfig contains authentication configuration
type AuthConfig struct {
	Type     string `yaml:"type"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Token    string `yaml:"token,omitempty"`
}

// TLSConfig contains TLS configuration
type TLSConfig struct {
	Enabled    bool   `yaml:"enabled"`
	CertFile   string `yaml:"certFile,omitempty"`
	KeyFile    string `yaml:"keyFile,omitempty"`
	CACertFile string `yaml:"caCertFile,omitempty"`
}

// USPConfig contains USP protocol configuration
type USPConfig struct {
	ControllerEndpointID string `yaml:"controllerEndpointId"`
	ProtocolVersion      string `yaml:"protocolVersion"`
	VersionCheck         bool   `yaml:"versionCheck"`
	AgentID              string `yaml:"agentId,omitempty"`
}

// CacheConfig contains cache configuration
type CacheConfig struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password,omitempty"`
	Database int    `yaml:"database"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`
	File       string `yaml:"file,omitempty"`
	MaxSize    int    `yaml:"maxSize,omitempty"`
	MaxBackups int    `yaml:"maxBackups,omitempty"`
	MaxAge     int    `yaml:"maxAge,omitempty"`
	Compress   bool   `yaml:"compress,omitempty"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(configPath string) (*Config, error) {
	// If no config path provided, try to find it
	if configPath == "" {
		configPath = findConfigFile()
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %v", configPath, err)
	}

	// Expand environment variables in config
	expanded := os.ExpandEnv(string(data))

	var config Config
	if err := yaml.Unmarshal([]byte(expanded), &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %v", configPath, err)
	}

	return &config, nil
}

// findConfigFile tries to locate the configuration file
func findConfigFile() string {
	// Priority order for config file locations
	locations := []string{
		"./config.yaml",
		"./configs/config.yaml",
		"./openusp.yaml",
		"./configs/openusp.yaml",
		"/etc/openusp/config.yaml",
		"/usr/local/etc/openusp/config.yaml",
	}

	for _, location := range locations {
		if _, err := os.Stat(location); err == nil {
			return location
		}
	}

	// Return default if none found
	return "./config.yaml"
}

// ValidateConfig validates the configuration
func (c *Config) ValidateConfig() error {
	if c.Service.Name == "" {
		return fmt.Errorf("service name is required")
	}

	if c.Database.Type == "" {
		return fmt.Errorf("database type is required")
	}

	if c.Database.Host == "" && c.Database.URI == "" {
		return fmt.Errorf("database host or URI is required")
	}

	return nil
}

// GetDatabaseURI returns the database connection URI
func (c *Config) GetDatabaseURI() string {
	if c.Database.URI != "" {
		return c.Database.URI
	}

	if c.Database.Username != "" && c.Database.Password != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
			c.Database.Username,
			c.Database.Password,
			c.Database.Host,
			c.Database.Port,
			c.Database.Name,
		)
	}

	return fmt.Sprintf("mongodb://%s:%d/%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)
}

// GetStompAddress returns the STOMP server address
func (c *Config) GetStompAddress() string {
	return fmt.Sprintf("%s:%d", c.MessageBus.STOMP.Host, c.MessageBus.STOMP.Port)
}

// GetMqttAddress returns the MQTT server address
func (c *Config) GetMqttAddress() string {
	return fmt.Sprintf("%s:%d", c.MessageBus.MQTT.Host, c.MessageBus.MQTT.Port)
}

// GetHTTPAddress returns the HTTP server address
func (c *Config) GetHTTPAddress() string {
	return fmt.Sprintf("%s:%d", c.Protocols.HTTP.Host, c.Protocols.HTTP.Port)
}

// GetGRPCAddress returns the gRPC server address
func (c *Config) GetGRPCAddress() string {
	return fmt.Sprintf("%s:%d", c.Protocols.GRPC.Host, c.Protocols.GRPC.Port)
}

// GetCacheAddress returns the cache server address
func (c *Config) GetCacheAddress() string {
	return fmt.Sprintf("%s:%d", c.Security.Cache.Host, c.Security.Cache.Port)
}