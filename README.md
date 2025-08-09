# Criage Common

Общий модуль для экосистемы пакетного менеджера Criage, содержащий типы данных, конфигурацию и утилиты, используемые всеми компонентами.

## 🏗️ Архитектура

Этот модуль предоставляет общие компоненты для:

- **criage-client** - основной CLI пакетного менеджера
- **criage-server** - HTTP сервер репозитория пакетов
- **criage-mcp** - MCP сервер для интеграции с AI

## 📦 Компоненты

### Types (`types/`)

Общие типы данных:

```go
import "github.com/criage-oss/criage-common/types"

// Основные типы
type PackageManifest struct { ... }
type PackageMetadata struct { ... }
type Repository struct { ... }
type PackageEntry struct { ... }
```

### Configuration (`config/`)

Управление конфигурацией:

```go
import "github.com/criage-oss/criage-common/config"

// Конфигурация по умолчанию
cfg := config.DefaultConfig()

// Конфигурация сервера
serverCfg := config.DefaultServerConfig()

// Конфигурация MCP
mcpCfg := config.DefaultMCPConfig()
```

### Archive (`archive/`)

Работа с архивами пакетов:

```go
import "github.com/criage-oss/criage-common/archive"

// Создание менеджера архивов
manager, err := archive.NewManager(cfg, "1.0.0")
defer manager.Close()

// Извлечение архива
err = manager.ExtractArchive("package.tar.zst", "./output", types.FormatTarZst)

// Определение формата
format := manager.DetectFormat("package.criage")

// Извлечение метаданных
metadata, err := manager.ExtractMetadataFromArchive("package.tar.zst", format)
```

## 🚀 Использование

### Добавление зависимости

```bash
go mod init your-project
go get github.com/criage-oss/criage-common
```

### Пример использования

```go
package main

import (
    "fmt"
    "github.com/criage-oss/criage-common/config"
    "github.com/criage-oss/criage-common/types"
    "github.com/criage-oss/criage-common/archive"
)

func main() {
    // Создаем конфигурацию
    cfg := config.DefaultConfig()
    
    // Создаем менеджер архивов
    archiveManager, err := archive.NewManager(cfg, "1.0.0")
    if err != nil {
        panic(err)
    }
    defer archiveManager.Close()
    
    // Определяем формат архива
    format := archiveManager.DetectFormat("example.tar.zst")
    fmt.Printf("Detected format: %s\n", format)
    
    // Извлекаем метаданные
    metadata, err := archiveManager.ExtractMetadataFromArchive("example.tar.zst", format)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    if metadata.PackageManifest != nil {
        fmt.Printf("Package: %s v%s\n", 
            metadata.PackageManifest.Name, 
            metadata.PackageManifest.Version)
    }
}
```

## 🔧 Поддерживаемые форматы архивов

- **tar.zst** - Tar с Zstandard сжатием (рекомендуется)
- **tar.lz4** - Tar с LZ4 сжатием (быстрый)
- **tar.xz** - Tar с XZ сжатием (компактный)
- **tar.gz** - Tar с Gzip сжатием (совместимость)
- **zip** - ZIP архивы (Windows)
- **criage** - Универсальное расширение с автоопределением формата

## 📊 Структуры данных

### PackageManifest

Описание пакета (`criage.yaml`):

```yaml
name: my-package
version: 1.0.0
description: Example package
author: Developer Name
license: MIT
dependencies:
  some-lib: "^1.2.0"
scripts:
  install: ./install.sh
```

### Repository

Конфигурация репозитория:

```yaml
repositories:
  - name: official
    url: https://packages.criage.ru
    priority: 100
    enabled: true
```

### PackageEntry

Запись в индексе репозитория с версиями и файлами для разных платформ.

## 🏷️ Версионирование

Модуль использует [Semantic Versioning](https://semver.org/):

- `1.0.0` - стабильная версия
- `1.x.x` - обратно совместимые изменения
- `2.0.0` - breaking changes

## 🤝 Интеграция

### С criage-client

```go
import "github.com/criage-oss/criage-common/config"
import "github.com/criage-oss/criage-common/archive"

cfg := config.DefaultConfig()
archMgr, _ := archive.NewManager(cfg, clientVersion)
```

### С criage-server

```go
import "github.com/criage-oss/criage-common/config"
import "github.com/criage-oss/criage-common/types"

serverCfg := config.DefaultServerConfig()
var index types.RepositoryIndex
```

### С criage-mcp

```go
import "github.com/criage-oss/criage-common/config"
import "github.com/criage-oss/criage-common/types"

mcpCfg := config.DefaultMCPConfig()
var response types.ApiResponse
```

## 🔗 Связанные проекты

- [criage-client](https://github.com/criage-oss/criage-client) - CLI пакетного менеджера
- [criage-server](https://github.com/criage-oss/criage-server) - HTTP сервер репозитория
- [criage-mcp](https://github.com/criage-oss/criage-mcp) - MCP сервер для AI интеграции

## 📄 Лицензия

MIT License. См. [LICENSE](LICENSE) для подробностей.

## 🤝 Участие в разработке

1. Fork проекта
2. Создайте feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit изменения (`git commit -m 'Add some AmazingFeature'`)
4. Push в branch (`git push origin feature/AmazingFeature`)
5. Откройте Pull Request

## 📞 Поддержка

- [GitHub Issues](https://github.com/criage-oss/criage-common/issues)
- [Документация](https://criage.ru)
- [Сообщество](https://github.com/criage-oss)
