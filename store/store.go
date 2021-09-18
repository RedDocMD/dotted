package store

import (
	"crypto/sha1"
	"fmt"
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

func LoadStore(config *config.Config) (*Store, error) {
	pathFilePath := Fs.Join(config.StoreLocation, "paths")
	pathFileBytes, err := Afs.ReadFile(pathFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load store")
	}
	paths := strings.Split(string(pathFileBytes), "\n")
	pathsDone := make(map[string]struct{})
	var dotFiles []*file.DotFile
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
