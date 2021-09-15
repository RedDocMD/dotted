package file

import (
	"bytes"
	"encoding/json"
	"errors"
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
	homedir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(homedir, ".config")
	err = os.Mkdir(configPath, os.ModeDir|os.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		t.Fatal(err)
	}
	filePath := filepath.Join(configPath, "dotted.yaml")
	osFile, err := os.Create(filePath)
	if err != nil {
		t.Fatal(err)
	}
	osFile.Close()

	file, err := NewDotFile(filePath, "first", true)
	assert.Equal(err, nil)
	if runtime.GOOS == "windows" {
		assert.Equal("195f56a15cad7a5576ad5fff1491db609aacd529", file.NameHash())
	} else {
		assert.Equal("1cc58199db412f2610d547f76fefc9f8b90aae8d", file.NameHash())
	}

	err = os.Remove(filePath)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDotFileToJSON(t *testing.T) {
	assert := assert.New(t)
	firstPath, _ := filepath.Abs(filepath.Join("testdata", "first.txt"))
	dotFile, _ := NewDotFile(firstPath, "first", true)
	dotFileJson := dotFile.MetadataToJSON()
	var values map[string]interface{}
	err := json.Unmarshal(dotFileJson, &values)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(values["Mnemonic"], "first")
	assert.Equal(values["HasHistory"], true)
}
