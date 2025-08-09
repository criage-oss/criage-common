package config

import (
	"os"
	"path/filepath"
	"runtime"
	
	"github.com/criage-oss/criage-common/types"
)

// Config основная конфигурация Criage
type Config struct {
	// Пути
	InstallPath   string `json:"installPath" yaml:"installPath"`
	CachePath     string `json:"cachePath" yaml:"cachePath"`
	TempPath      string `json:"tempPath" yaml:"tempPath"`
	ConfigPath    string `json:"configPath" yaml:"configPath"`
	
	// Сетевые настройки
	Timeout         int                `json:"timeout" yaml:"timeout"`
	MaxConnections  int                `json:"maxConnections" yaml:"maxConnections"`
	UserAgent       string             `json:"userAgent" yaml:"userAgent"`
	
	// Репозитории
	Repositories    []types.Repository `json:"repositories" yaml:"repositories"`
	
	// Сжатие
	CompressionLevel int    `json:"compressionLevel" yaml:"compressionLevel"`
	PreferredFormat  string `json:"preferredFormat" yaml:"preferredFormat"`
	
	// Другие настройки
	Parallel        bool   `json:"parallel" yaml:"parallel"`
	MaxParallel     int    `json:"maxParallel" yaml:"maxParallel"`
	Language        string `json:"language" yaml:"language"`
	Debug           bool   `json:"debug" yaml:"debug"`
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	criageDir := filepath.Join(homeDir, ".criage")
	
	return &Config{
		InstallPath:      filepath.Join(criageDir, "packages"),
		CachePath:        filepath.Join(criageDir, "cache"),
		TempPath:         filepath.Join(criageDir, "tmp"),
		ConfigPath:       filepath.Join(criageDir, "config.yaml"),
		Timeout:          30,
		MaxConnections:   10,
		UserAgent:        "Criage/1.0.0",
		Repositories:     defaultRepositories(),
		CompressionLevel: types.CompressionNormal,
		PreferredFormat:  "tar.zst",
		Parallel:         true,
		MaxParallel:      runtime.NumCPU(),
		Language:         "en",
		Debug:            false,
	}
}

// defaultRepositories возвращает репозитории по умолчанию
func defaultRepositories() []types.Repository {
	return []types.Repository{
		{
			Name:     "official",
			URL:      "https://packages.criage.ru",
			Priority: 100,
			Enabled:  true,
		},
	}
}

// ServerConfig конфигурация для сервера репозитория
type ServerConfig struct {
	// Сервер
	Host string `json:"host" yaml:"host"`
	Port int    `json:"port" yaml:"port"`
	
	// Пути
	StoragePath string `json:"storagePath" yaml:"storagePath"`
	IndexPath   string `json:"indexPath" yaml:"indexPath"`
	
	// Безопасность
	AuthEnabled bool   `json:"authEnabled" yaml:"authEnabled"`
	AuthToken   string `json:"authToken,omitempty" yaml:"authToken,omitempty"`
	
	// Ограничения
	MaxFileSize     int64    `json:"maxFileSize" yaml:"maxFileSize"`
	AllowedFormats  []string `json:"allowedFormats" yaml:"allowedFormats"`
	RateLimit       int      `json:"rateLimit" yaml:"rateLimit"`
	
	// Логирование
	LogLevel  string `json:"logLevel" yaml:"logLevel"`
	LogFile   string `json:"logFile,omitempty" yaml:"logFile,omitempty"`
	
	// CORS
	CORSEnabled bool     `json:"corsEnabled" yaml:"corsEnabled"`
	CORSOrigins []string `json:"corsOrigins,omitempty" yaml:"corsOrigins,omitempty"`
}

// DefaultServerConfig возвращает конфигурацию сервера по умолчанию
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Host:        "0.0.0.0",
		Port:        8080,
		StoragePath: "./packages",
		IndexPath:   "./index.json",
		AuthEnabled: false,
		MaxFileSize: 100 * 1024 * 1024, // 100MB
		AllowedFormats: []string{
			"tar.zst", "tar.lz4", "tar.xz", 
			"tar.gz", "zip", "criage",
		},
		RateLimit:   60, // requests per minute
		LogLevel:    "info",
		CORSEnabled: true,
		CORSOrigins: []string{"*"},
	}
}

// MCPConfig конфигурация для MCP сервера
type MCPConfig struct {
	// Пути к исполняемым файлам
	CriageClientPath string `json:"criageClientPath" yaml:"criageClientPath"`
	ConfigPath       string `json:"configPath" yaml:"configPath"`
	
	// Ограничения
	MaxConcurrency int `json:"maxConcurrency" yaml:"maxConcurrency"`
	Timeout        int `json:"timeout" yaml:"timeout"`
	
	// Логирование
	LogLevel string `json:"logLevel" yaml:"logLevel"`
	LogFile  string `json:"logFile,omitempty" yaml:"logFile,omitempty"`
	
	// Безопасность
	AllowedOperations []string `json:"allowedOperations" yaml:"allowedOperations"`
	RestrictedPaths   []string `json:"restrictedPaths,omitempty" yaml:"restrictedPaths,omitempty"`
}

// DefaultMCPConfig возвращает конфигурацию MCP сервера по умолчанию
func DefaultMCPConfig() *MCPConfig {
	return &MCPConfig{
		CriageClientPath: "criage",
		ConfigPath:       "",
		MaxConcurrency:   3,
		Timeout:          120,
		LogLevel:         "info",
		AllowedOperations: []string{
			"install", "uninstall", "search", 
			"list", "info", "update",
		},
		RestrictedPaths: []string{
			"/etc", "/usr", "/bin", "/sbin",
		},
	}
}
