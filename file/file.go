package file

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type DotFile struct {
	path           string
	mnemonic       string
	historyRoot    *HistoryNode
	currentHistory *HistoryNode
	hasHistory     bool
	content        *string // RI: hasHistory ^ (content != nil) == 1
}

func NewDotFile(path, mnemonic string, hasHistory bool) (*DotFile, error) {
	if !filepath.IsAbs(path) {
		return nil, fmt.Errorf("failed to create dot file: %s is not absolute path", path)
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create dot file")
	}
	defer file.Close()
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create dot file")
	}
	content := buf.String()
	if !hasHistory {
		dotFile := &DotFile{
			path:           path,
			mnemonic:       mnemonic,
			historyRoot:    nil,
			currentHistory: nil,
			hasHistory:     hasHistory,
			content:        &content,
		}
		return dotFile, nil
	}
	history := NewHistory(content)
	dotFile := &DotFile{
		path:           path,
		mnemonic:       mnemonic,
		historyRoot:    history,
		currentHistory: history,
		hasHistory:     hasHistory,
		content:        nil,
	}
	return dotFile, nil
}

func (file *DotFile) NameHash() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to retrieve user home directory")
		os.Exit(1)
	}
	path := file.path[len(homedir)+1:]
	sum := sha1.Sum([]byte(path))
	return fmt.Sprintf("%x", sum)
}

func (file *DotFile) AddCommit() (bool, error) {
	if !file.hasHistory {
		return false, fmt.Errorf("failed to create commit: file without history")
	}
	osFile, err := os.Open(file.path)
	if err != nil {
		return false, errors.Wrap(err, "failed to create commit")
	}
	defer osFile.Close()
	var buf bytes.Buffer
	_, err = io.Copy(&buf, osFile)
	if err != nil {
		return false, errors.Wrap(err, "failed to create commit")
	}
	node := file.currentHistory.AddCommit(buf.String())
	if node == nil {
		return false, nil
	} else {
		file.currentHistory = node
		return true, nil
	}
}

type jsonDotFileMetadata struct {
	Mnemonic       string
	HasHistory     bool
	CurrentHistory string // UUID of node
}

func (file *DotFile) MetadataToJSON() []byte {
	jsonFile := jsonDotFileMetadata{
		Mnemonic:       file.mnemonic,
		HasHistory:     file.hasHistory,
		CurrentHistory: file.currentHistory.uuid.String(),
	}
	bytes, err := json.Marshal(jsonFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to convert dot file to JSON: %v\n", file)
		os.Exit(1)
	}
	return bytes
}
