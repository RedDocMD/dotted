package file

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"runtime"
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

func TestDotFileNameHash(t *testing.T) {
	assert := assert.New(t)
	firstPath, _ := filepath.Abs(filepath.Join("testdata", "first.txt"))
	secondPath, _ := filepath.Abs(filepath.Join("testdata", "second.txt"))

	firstFile, err := NewDotFile(firstPath, "first", true)
	assert.Equal(err, nil)
	if runtime.GOOS == "windows" {
		assert.Equal("d7b2f8b446d2a3caa72f37bb90d4f1821e6340df", firstFile.NameHash())
	} else {
		assert.Equal("65c4309953d88144d1f3ee0694d068413a9edf11", firstFile.NameHash())
	}

	secondFile, err := NewDotFile(secondPath, "second", false)
	assert.Equal(err, nil)
	if runtime.GOOS == "windows" {
		assert.Equal("af1959c46e62c5944ab4441437c39fdbde3cd636", secondFile.NameHash())
	} else {
		assert.Equal("c539e8f55898550d5dc52225e0433e8b512e7032", secondFile.NameHash())
	}
}
