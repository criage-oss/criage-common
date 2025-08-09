package archive

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4/v4"
	"github.com/ulikunitz/xz"

	"github.com/criage-oss/criage-common/config"
	"github.com/criage-oss/criage-common/types"
)

// Manager управляет архивами пакетов
type Manager struct {
	config  *config.Config
	version string
	
	// Кодировщики/декодеры
	zstdEncoder *zstd.Encoder
	zstdDecoder *zstd.Decoder
}

// NewManager создает новый менеджер архивов
func NewManager(cfg *config.Config, version string) (*Manager, error) {
	manager := &Manager{
		config:  cfg,
		version: version,
	}

	// Инициализируем zstd encoder/decoder
	var err error
	manager.zstdEncoder, err = zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.EncoderLevelFromZstd(cfg.CompressionLevel)))
	if err != nil {
		return nil, fmt.Errorf("failed to create zstd encoder: %w", err)
	}

	manager.zstdDecoder, err = zstd.NewReader(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create zstd decoder: %w", err)
	}

	return manager, nil
}

// Close освобождает ресурсы менеджера
func (m *Manager) Close() error {
	if m.zstdEncoder != nil {
		m.zstdEncoder.Close()
	}
	if m.zstdDecoder != nil {
		m.zstdDecoder.Close()
	}
	return nil
}

// DetectFormat определяет формат архива по расширению файла
func (m *Manager) DetectFormat(filename string) types.ArchiveFormat {
	filename = strings.ToLower(filename)
	
	switch {
	case strings.HasSuffix(filename, ".tar.zst"):
		return types.FormatTarZst
	case strings.HasSuffix(filename, ".tar.lz4"):
		return types.FormatTarLZ4
	case strings.HasSuffix(filename, ".tar.xz"):
		return types.FormatTarXZ
	case strings.HasSuffix(filename, ".tar.gz"):
		return types.FormatTarGZ
	case strings.HasSuffix(filename, ".zip"):
		return types.FormatZip
	case strings.HasSuffix(filename, ".criage"):
		// .criage файлы могут быть в любом формате, пытаемся определить
		return m.detectCriageFormat(filename)
	default:
		return types.FormatTarZst // по умолчанию
	}
}

// detectCriageFormat пытается определить формат .criage файла
func (m *Manager) detectCriageFormat(filename string) types.ArchiveFormat {
	file, err := os.Open(filename)
	if err != nil {
		return types.FormatTarZst
	}
	defer file.Close()

	// Читаем первые несколько байт для определения формата
	header := make([]byte, 16)
	_, err = file.Read(header)
	if err != nil {
		return types.FormatTarZst
	}

	// Проверяем магические байты
	switch {
	case len(header) >= 4 && header[0] == 0x28 && header[1] == 0xB5 && header[2] == 0x2F && header[3] == 0xFD:
		return types.FormatTarZst // Zstandard
	case len(header) >= 4 && header[0] == 0x04 && header[1] == 0x22 && header[2] == 0x4D && header[3] == 0x18:
		return types.FormatTarLZ4 // LZ4
	case len(header) >= 6 && header[0] == 0xFD && header[1] == 0x37 && header[2] == 0x7A && header[3] == 0x58 && header[4] == 0x5A && header[5] == 0x00:
		return types.FormatTarXZ // XZ
	case len(header) >= 2 && header[0] == 0x1F && header[1] == 0x8B:
		return types.FormatTarGZ // Gzip
	case len(header) >= 4 && header[0] == 0x50 && header[1] == 0x4B && (header[2] == 0x03 || header[2] == 0x05):
		return types.FormatZip // ZIP
	default:
		return types.FormatTarZst
	}
}

// ExtractArchive извлекает архив в указанную директорию
func (m *Manager) ExtractArchive(archivePath, destDir string, format types.ArchiveFormat) error {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	switch format {
	case types.FormatZip:
		return m.extractZip(archivePath, destDir)
	default:
		return m.extractTar(archivePath, destDir, format)
	}
}

// extractTar извлекает tar архив с различными алгоритмами сжатия
func (m *Manager) extractTar(archivePath, destDir string, format types.ArchiveFormat) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Создаем декомпрессор в зависимости от формата
	var reader io.Reader
	switch format {
	case types.FormatTarZst:
		reader, err = m.zstdDecoder.DecodeAll(nil, nil)
		if err != nil {
			return err
		}
		reader = file
	case types.FormatTarLZ4:
		reader = lz4.NewReader(file)
	case types.FormatTarXZ:
		reader, err = xz.NewReader(file)
		if err != nil {
			return err
		}
	case types.FormatTarGZ:
		reader, err = gzip.NewReader(file)
		if err != nil {
			return err
		}
	default:
		reader = file
	}

	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(destDir, header.Name)

		// Проверяем безопасность пути
		if !strings.HasPrefix(target, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid path: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := m.extractFile(tarReader, target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		}
	}

	return nil
}

// extractZip извлекает ZIP архив
func (m *Manager) extractZip(archivePath, destDir string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		target := filepath.Join(destDir, file.Name)

		// Проверяем безопасность пути
		if !strings.HasPrefix(target, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid path: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(target, file.FileInfo().Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}

		if err := m.extractFile(fileReader, target, file.FileInfo().Mode()); err != nil {
			fileReader.Close()
			return err
		}
		fileReader.Close()
	}

	return nil
}

// extractFile извлекает отдельный файл
func (m *Manager) extractFile(src io.Reader, destPath string, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	dest, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)
	return err
}

// ExtractMetadataFromArchive извлекает метаданные из архива
func (m *Manager) ExtractMetadataFromArchive(archivePath string, format types.ArchiveFormat) (*types.PackageMetadata, error) {
	// Создаем временную директорию
	tempDir, err := os.MkdirTemp("", "criage-metadata-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	// Извлекаем только метаданные
	if err := m.extractMetadataOnly(archivePath, tempDir, format); err != nil {
		return nil, err
	}

	// Ищем файл метаданных
	metadataFile := filepath.Join(tempDir, ".criage-metadata.json")
	if _, err := os.Stat(metadataFile); os.IsNotExist(err) {
		// Пытаемся найти манифест пакета
		manifestFile := filepath.Join(tempDir, "criage.yaml")
		if _, err := os.Stat(manifestFile); os.IsNotExist(err) {
			return nil, fmt.Errorf("no metadata found in archive")
		}
		return m.createMetadataFromManifest(manifestFile)
	}

	// Загружаем метаданные
	data, err := os.ReadFile(metadataFile)
	if err != nil {
		return nil, err
	}

	var metadata types.PackageMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// extractMetadataOnly извлекает только файлы метаданных
func (m *Manager) extractMetadataOnly(archivePath, destDir string, format types.ArchiveFormat) error {
	// Упрощенная версия - извлекаем всё, но только для поиска метаданных
	return m.ExtractArchive(archivePath, destDir, format)
}

// createMetadataFromManifest создает метаданные из манифеста пакета
func (m *Manager) createMetadataFromManifest(manifestPath string) (*types.PackageMetadata, error) {
	// Здесь должна быть логика чтения YAML манифеста
	// Для упрощения возвращаем базовые метаданные
	return &types.PackageMetadata{
		CreatedBy: "criage",
		Version:   m.version,
	}, nil
}

// CreateArchiveWithMetadata создает архив с встроенными метаданными
func (m *Manager) CreateArchiveWithMetadata(sourceDir, outputPath string, format types.ArchiveFormat, includeFiles, excludeFiles []string, metadata *types.PackageMetadata) error {
	// Создаем временный файл метаданных
	tempDir, err := os.MkdirTemp("", "criage-build-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	metadataFile := filepath.Join(tempDir, ".criage-metadata.json")
	metadataData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(metadataFile, metadataData, 0644); err != nil {
		return err
	}

	// Создаем архив
	return m.createArchive(sourceDir, outputPath, format, includeFiles, excludeFiles, metadataFile)
}

// createArchive создает архив
func (m *Manager) createArchive(sourceDir, outputPath string, format types.ArchiveFormat, includeFiles, excludeFiles []string, metadataFile string) error {
	switch format {
	case types.FormatZip:
		return m.createZipArchive(sourceDir, outputPath, includeFiles, excludeFiles, metadataFile)
	default:
		return m.createTarArchive(sourceDir, outputPath, format, includeFiles, excludeFiles, metadataFile)
	}
}

// createTarArchive создает tar архив с сжатием
func (m *Manager) createTarArchive(sourceDir, outputPath string, format types.ArchiveFormat, includeFiles, excludeFiles []string, metadataFile string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Создаем компрессор
	var writer io.Writer
	switch format {
	case types.FormatTarZst:
		writer = m.zstdEncoder.EncodeAll(nil, nil)
		writer = file
	case types.FormatTarLZ4:
		writer = lz4.NewWriter(file)
	case types.FormatTarXZ:
		writer, err = xz.NewWriter(file)
		if err != nil {
			return err
		}
	case types.FormatTarGZ:
		writer = gzip.NewWriter(file)
	default:
		writer = file
	}

	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	// Добавляем метаданные
	if err := m.addFileToTar(tarWriter, metadataFile, ".criage-metadata.json"); err != nil {
		return err
	}

	// Добавляем файлы из источника
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Пропускаем корневую директорию
		if path == sourceDir {
			return nil
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// Проверяем фильтры
		if m.shouldExclude(relPath, includeFiles, excludeFiles) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		return m.addFileToTar(tarWriter, path, relPath)
	})
}

// createZipArchive создает ZIP архив
func (m *Manager) createZipArchive(sourceDir, outputPath string, includeFiles, excludeFiles []string, metadataFile string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	// Добавляем метаданные
	if err := m.addFileToZip(zipWriter, metadataFile, ".criage-metadata.json"); err != nil {
		return err
	}

	// Добавляем файлы из источника
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == sourceDir {
			return nil
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		if m.shouldExclude(relPath, includeFiles, excludeFiles) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		return m.addFileToZip(zipWriter, path, relPath)
	})
}

// addFileToTar добавляет файл в tar архив
func (m *Manager) addFileToTar(tarWriter *tar.Writer, sourcePath, archivePath string) error {
	info, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	header := &tar.Header{
		Name:    archivePath,
		Mode:    int64(info.Mode()),
		Size:    info.Size(),
		ModTime: info.ModTime(),
	}

	if info.IsDir() {
		header.Typeflag = tar.TypeDir
		return tarWriter.WriteHeader(header)
	}

	header.Typeflag = tar.TypeReg
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	file, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(tarWriter, file)
	return err
}

// addFileToZip добавляет файл в ZIP архив
func (m *Manager) addFileToZip(zipWriter *zip.Writer, sourcePath, archivePath string) error {
	info, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		_, err := zipWriter.Create(archivePath + "/")
		return err
	}

	file, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer, err := zipWriter.Create(archivePath)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}

// shouldExclude проверяет, должен ли файл быть исключен
func (m *Manager) shouldExclude(path string, includeFiles, excludeFiles []string) bool {
	// Проверяем исключения
	for _, pattern := range excludeFiles {
		if matched, _ := filepath.Match(pattern, path); matched {
			return true
		}
	}

	// Если есть список включений, проверяем его
	if len(includeFiles) > 0 {
		for _, pattern := range includeFiles {
			if matched, _ := filepath.Match(pattern, path); matched {
				return false
			}
		}
		return true
	}

	return false
}
