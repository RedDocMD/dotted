package file

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/RedDocMD/dotted/fs"
	"github.com/pkg/errors"
)

var Fs = fs.OsFs
var Afs = fs.OsAfs

type DotFile struct {
	path           string
	mnemonic       string
	historyRoot    *HistoryNode
	currentHistory *HistoryNode
	hasHistory     bool
	content        *string // RI: hasHistory ^ (content != nil) == 1
}

func (file *DotFile) Mnemonic() string {
	return file.mnemonic
}

func (file *DotFile) HasHistory() bool {
	return file.hasHistory
}

func (file *DotFile) RemoveHistory() {
	if !file.hasHistory {
		fmt.Fprintf(os.Stderr, "%s does not have a history, cannot remove it.\n", file.path)
		os.Exit(1)
	}
	currentContent := file.currentHistory.Content()
	file.content = &currentContent
	file.hasHistory = false
	file.currentHistory = nil
	file.historyRoot = nil
}

func (file *DotFile) InitHistory() {
	if file.hasHistory {
		fmt.Fprintf(os.Stderr, "%s already has a history, cannot init it.\n", file.path)
		os.Exit(1)
	}
	historyRoot := NewHistory(*file.content, currentTime())
	file.hasHistory = true
	file.historyRoot = historyRoot
	file.currentHistory = historyRoot
	file.content = nil
}

func currentTime() time.Time {
	timeNow := time.Now()
	timeNowString := timeNow.Format(time.UnixDate)
	timeNowUnix, _ := time.Parse(time.UnixDate, timeNowString)
	return timeNowUnix
}

func (file *DotFile) Path() string {
	return file.path
}

func NewDotFile(path, mnemonic string, hasHistory bool) (*DotFile, error) {
	if !Fs.IsAbs(path) {
		return nil, fmt.Errorf("failed to create dot file: %s is not absolute path", path)
	}
	buf, err := Afs.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create dot file")
	}
	content := string(buf)
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
	history := NewHistory(content, currentTime())
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

func (file *DotFile) RelativePath() string {
	homedir := Fs.UserHomeDir()
	path := file.path[len(homedir)+1:]
	return path
}

func (file *DotFile) RelativePathHash() string {
	path := file.RelativePath()
	sum := sha1.Sum([]byte(path))
	return fmt.Sprintf("%x", sum)
}

func (file *DotFile) AddCommit() (bool, error) {
	if !file.hasHistory {
		return false, fmt.Errorf("failed to create commit: file without history")
	}
	buf, err := Afs.ReadFile(file.path)
	if err != nil {
		return false, errors.Wrap(err, "failed to create commit")
	}
	node := file.currentHistory.AddCommit(string(buf), currentTime())
	if node == nil {
		return false, nil
	} else {
		file.currentHistory = node
		return true, nil
	}
}

func (file *DotFile) UpdateContent() (bool, error) {
	if file.hasHistory {
		return false, fmt.Errorf("failed to update content: file has history")
	}
	buf, err := Afs.ReadFile(file.path)
	if err != nil {
		return false, errors.Wrap(err, "failed to update content")
	}
	content := string(buf)
	changed := content != *file.content
	file.content = &content
	return changed, nil
}

type jsonDotFileMetadata struct {
	Mnemonic       string
	HasHistory     bool
	CurrentHistory string // UUID of node
}

func (file *DotFile) MetadataToJSON() []byte {
	var currentHistory string
	if file.currentHistory != nil {
		currentHistory = file.currentHistory.uuid.String()
	}
	jsonFile := jsonDotFileMetadata{
		Mnemonic:       file.mnemonic,
		HasHistory:     file.hasHistory,
		CurrentHistory: currentHistory,
	}
	bytes, err := json.Marshal(jsonFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to convert dot file to JSON: %v\n", file)
		os.Exit(1)
	}
	return bytes
}

func (file *DotFile) SaveToDisk(basePath string) error {
	if file.hasHistory {
		historyFilePath := Fs.Join(basePath, "history")
		historyFile, err := Afs.Create(historyFilePath)
		if err != nil {
			return errors.Wrap(err, "failed to save dot file to disk")
		}
		defer historyFile.Close()
		historyData := file.historyRoot.ToJSON()
		historyDataBuf := bytes.NewBuffer(historyData)
		_, err = io.Copy(historyFile, historyDataBuf)
		if err != nil {
			return errors.Wrap(err, "failed to save dot file to disk")
		}
	}

	var content string
	if file.hasHistory {
		content = *file.historyRoot.content
	} else {
		content = *file.content
	}
	contentFilePath := Fs.Join(basePath, "content")
	contentFile, err := Afs.Create(contentFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to save dot file to disk")
	}
	defer contentFile.Close()
	_, err = contentFile.WriteString(content)
	if err != nil {
		return errors.Wrap(err, "failed to save dot file to disk")
	}

	metadataFilePath := Fs.Join(basePath, "metadata")
	metadataFile, err := Afs.Create(metadataFilePath)
	if err != nil {
		return errors.Wrap(err, "failed to save dot file to disk")
	}
	defer metadataFile.Close()
	metadata := file.MetadataToJSON()
	metadataBuf := bytes.NewBuffer(metadata)
	_, err = io.Copy(metadataFile, metadataBuf)
	if err != nil {
		return errors.Wrap(err, "failed to save dot file to disk")
	}
	return nil
}

var BasePathNotFound = errors.New("base path directory not found")

func LoadDotFileFromDisk(basePath, dotFilePath string) (*DotFile, error) {
	if !Fs.IsAbs(dotFilePath) {
		return nil, fmt.Errorf(fmt.Sprintf("failed to read dot file from disk: %s is not absolute path", dotFilePath))
	}
	if exists, err := Afs.DirExists(basePath); err == nil && !exists {
		return nil, BasePathNotFound
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "failed to check for existence of %s\n", basePath)
		os.Exit(1)
	}
	metadataFilePath := Fs.Join(basePath, "metadata")
	metadataBytes, err := Afs.ReadFile(metadataFilePath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to read dot file from disk: %s", basePath))
	}
	var metadata jsonDotFileMetadata
	err = json.Unmarshal(metadataBytes, &metadata)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to read dot file from disk: %s", basePath))
	}
	contentFilePath := Fs.Join(basePath, "content")
	contentBytes, err := Afs.ReadFile(contentFilePath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to read dot file from disk: %s", basePath))
	}
	content := string(contentBytes)
	var historyRoot, currentHistory *HistoryNode
	var dotFileContent *string
	if metadata.HasHistory {
		historyFilePath := Fs.Join(basePath, "history")
		historyFileBytes, err := Afs.ReadFile(historyFilePath)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to read dot file from disk: %s", basePath))
		}
		historyRoot, err = FromJSON(historyFileBytes, content)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to read dot file from disk: %s", basePath))
		}
		currentHistory = historyRoot.NodeWithUUID(metadata.CurrentHistory)
		if currentHistory == nil {
			return nil, fmt.Errorf(fmt.Sprintf("failed to read dot file from disk, %s not found as current history", metadata.CurrentHistory))
		}
	} else {
		dotFileContent = &content
	}
	dotFile := &DotFile{
		path:           dotFilePath,
		mnemonic:       metadata.Mnemonic,
		historyRoot:    historyRoot,
		currentHistory: currentHistory,
		hasHistory:     metadata.HasHistory,
		content:        dotFileContent,
	}
	return dotFile, nil
}
