package file

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDotFile(t *testing.T) {
	assert := assert.New(t)
	firstPath := filepath.Join("testdata", "first.txt")
	firstPathAbs, _ := filepath.Abs(firstPath)

	_, err := NewDotFile(firstPath, "", false)
	assert.NotEqual(err, nil)
	dotFileWithHistory, err := NewDotFile(firstPathAbs, "test", true)
	assert.Equal(err, nil)
	assert.NotEqual(dotFileWithHistory, nil)
	dotFileWithoutHistory, err := NewDotFile(firstPathAbs, "test", false)
	assert.Equal(err, nil)
	assert.NotEqual(dotFileWithoutHistory, nil)
}

func TestCommitDotFile(t *testing.T) {
	assert := assert.New(t)
	firstPath, _ := filepath.Abs(filepath.Join("testdata", "first.txt"))
	secondPath, _ := filepath.Abs(filepath.Join("testdata", "second.txt"))

	dotFileWithHistory, _ := NewDotFile(firstPath, "test", true)
	dotFileWithoutHistory, _ := NewDotFile(firstPath, "test", false)

	firstContents, err := os.ReadFile(firstPath)
	if err != nil {
		t.Fatal(err)
	}
	secondContents, err := os.ReadFile(secondPath)
	if err != nil {
		t.Fatal(err)
	}
	firstBuf := bytes.NewBuffer(firstContents)
	secondBuf := bytes.NewBuffer(secondContents)
	copyToFile(t, firstPath, secondBuf)

	changed, err := dotFileWithHistory.AddCommit()
	assert.Equal(err, nil)
	assert.True(changed)

	changed, err = dotFileWithoutHistory.AddCommit()
	assert.Error(err, "failed to create commit: file without history")
	assert.False(changed)

	copyToFile(t, firstPath, firstBuf)
}

func copyToFile(t *testing.T, path string, content *bytes.Buffer) {
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	_, err = io.Copy(file, content)
	if err != nil {
		t.Fatal(err)
	}
}
