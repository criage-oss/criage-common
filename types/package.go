package types

import (
	"time"
)

// PackageManifest содержит основную информацию о пакете
type PackageManifest struct {
	// Основная информация
	Name        string `json:"name" yaml:"name"`
	Version     string `json:"version" yaml:"version"`
	Description string `json:"description" yaml:"description"`
	Author      string `json:"author" yaml:"author"`
	License     string `json:"license" yaml:"license"`
	Homepage    string `json:"homepage,omitempty" yaml:"homepage,omitempty"`
	Repository  string `json:"repository,omitempty" yaml:"repository,omitempty"`

	// Ключевые слова и категории
	Keywords []string `json:"keywords,omitempty" yaml:"keywords,omitempty"`

	// Зависимости
	Dependencies map[string]string `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	DevDeps      map[string]string `json:"devDependencies,omitempty" yaml:"devDependencies,omitempty"`

	// Скрипты жизненного цикла
	Scripts map[string]string `json:"scripts,omitempty" yaml:"scripts,omitempty"`

	// Файлы для включения/исключения
	Files   []string `json:"files,omitempty" yaml:"files,omitempty"`
	Exclude []string `json:"exclude,omitempty" yaml:"exclude,omitempty"`

	// Поддерживаемые платформы
	Arch []string `json:"arch,omitempty" yaml:"arch,omitempty"`
	OS   []string `json:"os,omitempty" yaml:"os,omitempty"`

	// Требования
	MinVersion string `json:"minVersion,omitempty" yaml:"minVersion,omitempty"`

	// Хуки жизненного цикла
	Hooks *PackageHooks `json:"hooks,omitempty" yaml:"hooks,omitempty"`

	// Дополнительные метаданные
	Metadata map[string]any `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// PackageHooks определяет скрипты жизненного цикла пакета
type PackageHooks struct {
	PreInstall  []string `json:"preInstall,omitempty" yaml:"preInstall,omitempty"`
	PostInstall []string `json:"postInstall,omitempty" yaml:"postInstall,omitempty"`
	PreRemove   []string `json:"preRemove,omitempty" yaml:"preRemove,omitempty"`
	PostRemove  []string `json:"postRemove,omitempty" yaml:"postRemove,omitempty"`
	PreUpdate   []string `json:"preUpdate,omitempty" yaml:"preUpdate,omitempty"`
	PostUpdate  []string `json:"postUpdate,omitempty" yaml:"postUpdate,omitempty"`
}

// BuildManifest содержит информацию о сборке пакета
type BuildManifest struct {
	Name         string            `json:"name" yaml:"name"`
	Version      string            `json:"version" yaml:"version"`
	BuildScript  string            `json:"buildScript,omitempty" yaml:"buildScript,omitempty"`
	OutputDir    string            `json:"outputDir" yaml:"outputDir"`
	IncludeFiles []string          `json:"includeFiles" yaml:"includeFiles"`
	ExcludeFiles []string          `json:"excludeFiles,omitempty" yaml:"excludeFiles,omitempty"`
	Compression  CompressionConfig `json:"compression" yaml:"compression"`
	Targets      []BuildTarget     `json:"targets" yaml:"targets"`
	Environment  map[string]string `json:"environment,omitempty" yaml:"environment,omitempty"`
}

// CompressionConfig настройки сжатия
type CompressionConfig struct {
	Format string `json:"format" yaml:"format"`
	Level  int    `json:"level" yaml:"level"`
}

// BuildTarget целевая платформа для сборки
type BuildTarget struct {
	OS   string `json:"os" yaml:"os"`
	Arch string `json:"arch" yaml:"arch"`
}

// PackageMetadata метаданные пакета со встроенной информацией о сборке
type PackageMetadata struct {
	PackageManifest *PackageManifest `json:"package,omitempty"`
	BuildManifest   *BuildManifest   `json:"build,omitempty"`
	CompressionType string           `json:"compressionType"`
	CreatedAt       time.Time        `json:"createdAt"`
	CreatedBy       string           `json:"createdBy"`
	Version         string           `json:"version"`
}

// PackageInfo информация об установленном пакете
type PackageInfo struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Author       string            `json:"author"`
	InstallDate  time.Time         `json:"installDate"`
	InstallPath  string            `json:"installPath"`
	Global       bool              `json:"global"`
	Dependencies map[string]string `json:"dependencies"`
	Size         int64             `json:"size"`
	Files        []string          `json:"files"`
	Scripts      map[string]string `json:"scripts"`
}

// ArchiveFormat формат архива
type ArchiveFormat string

const (
	FormatTarZst ArchiveFormat = "tar.zst"
	FormatTarLZ4 ArchiveFormat = "tar.lz4"
	FormatTarXZ  ArchiveFormat = "tar.xz"
	FormatTarGZ  ArchiveFormat = "tar.gz"
	FormatZip    ArchiveFormat = "zip"
)

// CompressionLevel уровни сжатия
const (
	CompressionFastest = 1
	CompressionNormal  = 5
	CompressionBest    = 9
)
