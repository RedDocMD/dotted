package file

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/RedDocMD/dotted/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DotFileTestSuite struct {
	suite.Suite
	configPath        string
	firstPath         string
	secondPath        string
	firstRelativePath string
	storePath         string
}

const globalFirstFileContent = `This is the first line`
const globalSecondFileContent = `This is the first line
This is the second line`

func (suite *DotFileTestSuite) SetupSuite() {
	Fs = fs.MockFs
	Afs = fs.MockAfs
}

func (suite *DotFileTestSuite) SetupTest() {
	pwd, err := os.Getwd()
	if err != nil {
		suite.T().Fatal(err)
	}
	homedir, err := os.UserHomeDir()
	if err != nil {
		suite.T().Fatal(err)
	}
	basedir := filepath.Join(pwd, "testdata")
	Fs.MkdirAll(basedir, 0755)
	Fs.Mkdir("testdata", 0755)
	suite.firstPath = filepath.Join(basedir, "first.txt")
	Afs.WriteFile(suite.firstPath, []byte(globalFirstFileContent), 0644)
	suite.firstRelativePath = filepath.Join("testdata", "first.txt")
	Afs.WriteFile(suite.firstRelativePath, []byte(globalFirstFileContent), 0644)
	suite.secondPath = filepath.Join(basedir, "second.txt")
	Afs.WriteFile(suite.secondPath, []byte(globalSecondFileContent), 0644)
	suite.configPath = filepath.Join(homedir, ".config", "dotted.yaml")
	Afs.Create(suite.configPath)
	suite.storePath = filepath.Join("testdir", "store")
	Afs.Mkdir(suite.storePath, 0644)
}

func (suite *DotFileTestSuite) TearDownTest() {
	Fs.RemoveAll("/")
}

func (suite *DotFileTestSuite) TearDownSuite() {
	Fs = fs.OsFs
	Afs = fs.OsAfs
}

func TestDotFileTestSuite(t *testing.T) {
	suite.Run(t, new(DotFileTestSuite))
}

func (suite *DotFileTestSuite) TestCreateDotFile() {
	assert := assert.New(suite.T())

	_, err := NewDotFile(suite.firstRelativePath, "", false)
	assert.NotEqual(err, nil)
	dotFileWithHistory, err := NewDotFile(suite.firstPath, "test", true)
	assert.Equal(err, nil)
	assert.NotEqual(dotFileWithHistory, nil)
	dotFileWithoutHistory, err := NewDotFile(suite.firstPath, "test", false)
	assert.Equal(err, nil)
	assert.NotEqual(dotFileWithoutHistory, nil)
}

func (suite *DotFileTestSuite) TestCommitDotFile() {
	assert := assert.New(suite.T())

	dotFileWithHistory, _ := NewDotFile(suite.firstPath, "test", true)
	dotFileWithoutHistory, _ := NewDotFile(suite.firstPath, "test", false)

	err := Afs.Remove(suite.firstPath)
	if err != nil {
		suite.T().Fatal(err)
	}
	err = Afs.Rename(suite.secondPath, suite.firstPath)
	if err != nil {
		suite.T().Fatal(err)
	}

	changed, err := dotFileWithHistory.AddCommit()
	assert.Equal(err, nil)
	assert.True(changed)

	changed, err = dotFileWithoutHistory.AddCommit()
	assert.Error(err, "failed to create commit: file without history")
	assert.False(changed)
}

func (suite *DotFileTestSuite) TestDotFileNameHash() {
	assert := assert.New(suite.T())
	file, err := NewDotFile(suite.configPath, "config", true)
	assert.Equal(err, nil)
	// if runtime.GOOS == "windows" {
	// 	assert.Equal("195f56a15cad7a5576ad5fff1491db609aacd529", file.NameHash())
	// } else {
	// 	assert.Equal("1cc58199db412f2610d547f76fefc9f8b90aae8d", file.NameHash())
	// }
	assert.Equal("1cc58199db412f2610d547f76fefc9f8b90aae8d", file.NameHash())
}

func (suite *DotFileTestSuite) TestDotFileMetadataToJSON() {
	assert := assert.New(suite.T())
	dotFile, _ := NewDotFile(suite.firstPath, "first", true)
	dotFileJson := dotFile.MetadataToJSON()
	var values map[string]interface{}
	err := json.Unmarshal(dotFileJson, &values)
	if err != nil {
		suite.T().Fatal(err)
	}
	assert.Equal(values["Mnemonic"], "first")
	assert.Equal(values["HasHistory"], true)
}

func (suite *DotFileTestSuite) TestDotFileStoreAndLoad() {
	assert := assert.New(suite.T())
	dotFile, _ := NewDotFile(suite.firstPath, "first", true)
	err := dotFile.SaveToDisk(suite.storePath)
	assert.Equal(err, nil)
	restoredDotFile, err := LoadDotFileFromDisk(suite.storePath, suite.firstPath)
	assert.Equal(err, nil)
	assert.Equal(dotFile, restoredDotFile)
}
