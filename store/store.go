package store

import (
	"crypto/sha1"
	"fmt"
	"os"
	"strings"

	"github.com/RedDocMD/dotted/config"
	"github.com/RedDocMD/dotted/file"
	"github.com/RedDocMD/dotted/fs"
	"github.com/pkg/errors"
)

var Fs = fs.OsFs
var Afs = fs.OsAfs

type Store struct {
	files []*file.DotFile
	path  string
	name  string
}

func LoadFromDisk(config *config.Config) (*Store, error) {
	pathsDone := make(map[string]struct{})
	var dotFiles []*file.DotFile

	pathFilePath := Fs.Join(config.StoreLocation, "paths")
	pathFileBytes, err := Afs.ReadFile(pathFilePath)
	if err == nil || os.IsNotExist(err) {
		paths := strings.Split(string(pathFileBytes), "\n")
		if len(paths[len(paths)-1]) == 0 {
			paths = paths[:len(paths)-1]
		}
		for _, path := range paths {
			basePath := Fs.Join(config.StoreLocation, storePath(path))
			dotFile, err := file.LoadDotFileFromDisk(basePath, Fs.Abs(path))
			if err != nil && !errors.Is(err, file.BasePathNotFound) {
				return nil, errors.Wrap(err, "failed to load store")
			}
			var fileInStore, fileInConfig, fileHasHistory bool
			fileInStore = err == nil
			if containsPath(path, config.WithHistory) {
				fileInConfig = true
				fileHasHistory = true
			} else if containsPath(path, config.WithoutHistory) {
				fileInConfig = true
				fileHasHistory = false
			}
			if fileInConfig && fileInStore {
				if dotFile.HasHistory() && !fileHasHistory {
					dotFile.RemoveHistory()
				} else if !dotFile.HasHistory() && fileHasHistory {
					dotFile.InitHistory()
				}
				dotFiles = append(dotFiles, dotFile)
				pathsDone[path] = struct{}{}
			} else if !fileInConfig && fileInStore {
				err = Afs.RemoveAll(basePath)
				if err != nil {
					return nil, errors.Wrap(err, "failed to load store")
				}
			} else if !fileInConfig && !fileInStore {
				return nil, errors.New(fmt.Sprintf("failed to load store: store in inconsistent state: directory for %s listed but not found", path))
			}
		}
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to load store")
	}
	for _, entry := range config.WithHistory {
		path := entry.Path
		if _, ok := pathsDone[path]; !ok {
			dotFile, err := file.NewDotFile(Fs.Abs(path), entry.Mnemonic, true)
			if err != nil {
				return nil, errors.Wrap(err, "failed to load store")
			}
			dotFiles = append(dotFiles, dotFile)
		}
	}
	for _, entry := range config.WithoutHistory {
		path := entry.Path
		if _, ok := pathsDone[path]; !ok {
			dotFile, err := file.NewDotFile(Fs.Abs(path), entry.Mnemonic, false)
			if err != nil {
				return nil, errors.Wrap(err, "failed to load store")
			}
			dotFiles = append(dotFiles, dotFile)
		}
	}
	store := &Store{
		files: dotFiles,
		path:  config.StoreLocation,
		name:  config.Name,
	}
	return store, nil
}

func containsPath(path string, entries []config.FileEntry) bool {
	for _, entry := range entries {
		if entry.Path == path {
			return true
		}
	}
	return false
}

func storePath(path string) string {
	sum := sha1.Sum([]byte(path))
	return fmt.Sprintf("%x", sum)
}

func (store *Store) SaveToDisk() error {
	err := makeDirIfNotExist(store.path)
	if err != nil {
		return errors.Wrap(err, "failed to save store to disk")
	}
	var pathFileContents string
	for _, file := range store.files {
		pathFileContents += file.RelativePath() + "\n"
	}
	err = Afs.WriteFile(Fs.Join(store.path, "paths"), []byte(pathFileContents), 0644)
	if err != nil {
		return errors.Wrap(err, "failed to save store to disk")
	}
	for _, file := range store.files {
		fileDir := Fs.Join(store.path, file.RelativePathHash())
		err = makeDirIfNotExist(fileDir)
		if err != nil {
			return errors.Wrap(err, "failed to save store to disk")
		}
		err = file.SaveToDisk(fileDir)
		if err != nil {
			return errors.Wrap(err, "failed to save store to disk")
		}
	}
	return nil
}

func makeDirIfNotExist(dirPath string) error {
	pathExists, err := Afs.Exists(dirPath)
	if err != nil {
		return err
	}
	if pathExists {
		dirExists, err := Afs.DirExists(dirPath)
		if err != nil {
			return nil
		}
		if !dirExists {
			return fmt.Errorf("%s already exists and is not a directory", dirPath)
		}
		return nil
	}
	err = Afs.MkdirAll(dirPath, 0755)
	return err
}
