package file

import (
	"bytes"
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
}

func NewDotFile(path, mnemonic string) (*DotFile, error) {
	if !filepath.IsAbs(path) {
		return nil, fmt.Errorf("failed to create dot file: %s is not absolute path", path)
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create dot file")
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create dot file")
	}
	history := NewHistory(buf.String())
	dotFile := &DotFile{
		path:           path,
		mnemonic:       mnemonic,
		historyRoot:    history,
		currentHistory: history,
	}
	return dotFile, nil
}

func (file *DotFile) AddCommit() (bool, error) {
	osFile, err := os.Open(file.path)
	if err != nil {
		return false, errors.Wrap(err, "failed to dot file")
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, osFile)
	if err != nil {
		return false, errors.Wrap(err, "failed to dot file")
	}
	node := file.currentHistory.AddCommit(buf.String())
	if node != nil {
		return false, nil
	} else {
		file.currentHistory = node
		return true, nil
	}
}
