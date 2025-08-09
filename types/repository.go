package types

import "time"

// Repository информация о репозитории пакетов
type Repository struct {
	Name     string `json:"name" yaml:"name"`
	URL      string `json:"url" yaml:"url"`
	Priority int    `json:"priority" yaml:"priority"`
	Enabled  bool   `json:"enabled" yaml:"enabled"`
	
	// Авторизация
	AuthToken string `json:"authToken,omitempty" yaml:"authToken,omitempty"`
	Username  string `json:"username,omitempty" yaml:"username,omitempty"`
	Password  string `json:"password,omitempty" yaml:"password,omitempty"`
}

// PackageEntry запись о пакете в репозитории
type PackageEntry struct {
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Author        string         `json:"author"`
	License       string         `json:"license"`
	Homepage      string         `json:"homepage,omitempty"`
	Repository    string         `json:"repository,omitempty"`
	Keywords      []string       `json:"keywords,omitempty"`
	LatestVersion string         `json:"latestVersion"`
	Versions      []VersionEntry `json:"versions"`
	Downloads     int64          `json:"downloads"`
	Updated       time.Time      `json:"updated"`
}

// VersionEntry информация о версии пакета
type VersionEntry struct {
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Dependencies map[string]string `json:"dependencies,omitempty"`
	DevDeps      map[string]string `json:"devDependencies,omitempty"`
	Files        []FileEntry       `json:"files"`
	Size         int64             `json:"size"`
	Checksum     string            `json:"checksum"`
	Uploaded     time.Time         `json:"uploaded"`
	Downloads    int64             `json:"downloads"`
}

// FileEntry информация о файле пакета
type FileEntry struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Format   string `json:"format"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Checksum string `json:"checksum"`
}

// SearchResult результат поиска
type SearchResult struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Downloads   int64     `json:"downloads"`
	Updated     time.Time `json:"updated"`
	Score       float64   `json:"score"`
}

// RepositoryIndex индекс репозитория
type RepositoryIndex struct {
	LastUpdated   time.Time               `json:"lastUpdated"`
	TotalPackages int                     `json:"totalPackages"`
	Packages      map[string]*PackageEntry `json:"packages"`
	Statistics    *Statistics             `json:"statistics"`
}

// Statistics статистика репозитория
type Statistics struct {
	TotalDownloads    int64            `json:"totalDownloads"`
	PackagesByLicense map[string]int   `json:"packagesByLicense"`
	PackagesByAuthor  map[string]int   `json:"packagesByAuthor"`
	PopularPackages   []string         `json:"popularPackages"`
}

// ApiResponse стандартный ответ API
type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// PackageListResponse ответ со списком пакетов
type PackageListResponse struct {
	Packages []PackageEntry `json:"packages"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PerPage  int            `json:"perPage"`
}

// UploadRequest запрос на загрузку пакета
type UploadRequest struct {
	PackageName string `json:"packageName"`
	Version     string `json:"version"`
	Filename    string `json:"filename"`
	Size        int64  `json:"size"`
	Checksum    string `json:"checksum"`
}

// UploadResponse ответ на загрузку пакета
type UploadResponse struct {
	PackageName string `json:"packageName"`
	Version     string `json:"version"`
	Filename    string `json:"filename"`
	UploadURL   string `json:"uploadUrl,omitempty"`
	Status      string `json:"status"`
}
