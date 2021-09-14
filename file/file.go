package file

import (
	"bytes"
	"io"
	"os"

	"github.com/pkg/errors"
)

type DotFile struct {
	path           string
	mnemonic       string
	historyRoot    *HistoryNode
	currentHistory *HistoryNode
}

func NewDotFile(path, mnemonic string) (*DotFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dot file")
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dot file")
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
