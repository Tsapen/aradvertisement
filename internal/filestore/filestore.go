package filestore

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// Storage struct contains info about useful paths.
type Storage struct {
	storagePath string
}

const (
	storageDirName = "gltf_storage"
	arPage         = "AR.html"
)

// FileMetadata contains metadata for object found.
type FileMetadata struct {
	Username string
	Name     string
}

// NewFM returns new file metadata.
func NewFM(username string, name int) FileMetadata {
	return FileMetadata{
		Username: username,
		Name:     strconv.Itoa(name),
	}
}

// PrepareStorage prepares storage.
func PrepareStorage(path string) (*Storage, error) {
	var storagePath = filepath.Join(path, storageDirName)
	if err := os.Mkdir(storagePath, 0777); err != nil && !os.IsExist(err) {
		return nil, err
	}

	if err := os.Chdir(storagePath); err != nil {
		return nil, err
	}

	var s = Storage{
		storagePath: storagePath,
	}

	return &s, nil
}

// CreateUserDir creates user glTF files  directory.
func (s *Storage) CreateUserDir(dir string) error {
	var dirPath = filepath.Join(s.storagePath, dir)
	if err := os.Mkdir(dirPath, 0777); err != nil && !os.IsExist(err) {
		return err
	}

	return nil
}

// DeleteUserDir deletes user glTF files  directory.
func (s *Storage) DeleteUserDir(dir string) error {
	return os.RemoveAll(filepath.Join(s.storagePath, dir))
}

// WriteFile writes file to path.
func (s *Storage) WriteFile(fm FileMetadata, b []byte) error {
	var filePath = filepath.Join(s.storagePath, fm.Username, fm.Name)
	var f, err = os.Create(filePath)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	return err
}

// ReadFile reads file.
func (s *Storage) ReadFile(fm FileMetadata) ([]byte, error) {
	var filePath = filepath.Join(s.storagePath, fm.Username, fm.Name)
	return ioutil.ReadFile(filePath)
}

// DeleteFile deletes file.
func (s *Storage) DeleteFile(fm FileMetadata) error {
	var filePath = filepath.Join(s.storagePath, fm.Username, fm.Name)
	return os.Remove(filePath)
}

// TODO delete

// GetARPage returns ar page.
func (s *Storage) GetARPage() ([]byte, error) {
	return s.getPage(arPage)
}

func (s *Storage) getPage(page string) ([]byte, error) {
	var path = filepath.Join("..", "ara-personal", "src", "views", "html", page)

	var f, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return f, nil
}
