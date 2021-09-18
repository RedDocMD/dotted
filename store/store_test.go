package store

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/RedDocMD/dotted/config"
	"github.com/RedDocMD/dotted/file"
	"github.com/RedDocMD/dotted/fs"
	"github.com/stretchr/testify/suite"
)

type StoreSuite struct {
	suite.Suite
}

func (suite *StoreSuite) SetupSuite() {
	Fs = fs.MockFs
	file.Fs = fs.MockFs
	config.Fs = fs.MockFs
	Afs = fs.MockAfs
	file.Afs = fs.MockAfs
	config.Afs = fs.MockAfs
}

func (suite *StoreSuite) TearDownSuite() {
	Fs = fs.OsFs
	file.Fs = fs.OsFs
	config.Fs = fs.OsFs
	Afs = fs.OsAfs
	file.Afs = fs.OsAfs
	config.Afs = fs.OsAfs
}

func (suite *StoreSuite) SetupTest() {
	var history, metadata string
	Afs.Mkdir("store", 0755)
	paths := `.config/alacritty/alacritty.yml
.tmux.conf`
	Afs.WriteFile("store/paths", []byte(paths), 0644)
	Afs.Mkdir("store/14b4f00abd93c6222516ff054e4a9f66295d03fa", 0755)
	buf, err := os.ReadFile(filepath.Join("testdata", "alacritty.yml"))
	if err != nil {
		suite.T().Fatal(err)
	}
	Afs.WriteFile("store/14b4f00abd93c6222516ff054e4a9f66295d03fa/content", buf, 0644)
	metadata = metadataJson("alacritty", "887cd650-21c0-4d1f-8e3f-c76425f550b2", true)
	Afs.WriteFile("store/14b4f00abd93c6222516ff054e4a9f66295d03fa/metadata", []byte(metadata), 0644)
	history = `
[
	{
		"Parent": "",
		"Patches": "",
		"Checksum": "555298afac7ed1ffbf44e9a6bc7afc09e4049ec8",
		"Children": ["887cd650-21c0-4d1f-8e3f-c76425f550b2"],
		"Uuid": "d032a2c2-d846-4f68-b055-5964a210d194"
	},
	{
		"Parent": "d032a2c2-d846-4f68-b055-5964a210d194",
		"Patches": "@@ -860,16 +860,31 @@\n OR: %221%22%0A\n+  YOLO: wuddup%0A\n %0A#window\n",
		"Checksum": "cd4c1c9693b8fc014ddf30c1f6bf261cbba4e777",
		"Children": [],
		"Uuid": "887cd650-21c0-4d1f-8e3f-c76425f550b2"
	}
]
	`
	Afs.WriteFile("store/14b4f00abd93c6222516ff054e4a9f66295d03fa/history", []byte(history), 0644)

	Afs.Mkdir("store/97aa776c8b768a52732c7978fd5f0af5ce5a1135", 0755)
	metadata = metadataJson("tmux", "", false)
	Afs.WriteFile("store/97aa776c8b768a52732c7978fd5f0af5ce5a1135/metadata", []byte(metadata), 0644)
	buf, err = os.ReadFile(filepath.Join("testdata", "tmux.conf"))
	if err != nil {
		suite.T().Fatal(err)
	}
	Afs.WriteFile("store/97aa776c8b768a52732c7978fd5f0af5ce5a1135/content", buf, 0644)
	Afs.WriteFile(Fs.Abs(".tmux.conf"), buf, 0644)
	buf, err = os.ReadFile(filepath.Join("testdata", "alacritty2.yml"))
	if err != nil {
		suite.T().Fatal(err)
	}
	Afs.MkdirAll(Fs.Abs(".config/alacritty"), 0755)
	Afs.WriteFile(Fs.Abs(".config/alacritty/alacritty.yml"), buf, 0644)
}

func (suite *StoreSuite) TearDownTest() {
	Afs.RemoveAll("/")
}

func metadataJson(mnemonic, currentHistory string, hasHistory bool) string {
	return fmt.Sprintf("{\"Mnemonic\":\"%s\", \"CurrentHistory\": \"%s\", \"HasHistory\": %t}", mnemonic, currentHistory, hasHistory)
}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, &StoreSuite{})
}

func (suite *StoreSuite) TestLoadStoreMatchConfig() {
	config := &config.Config{
		Name: "Linux",
		WithHistory: []config.FileEntry{
			{
				Path:     ".config/alacritty/alacritty.yml",
				Mnemonic: "alacritty",
			},
		},
		WithoutHistory: []config.FileEntry{
			{
				Path:     ".tmux.conf",
				Mnemonic: "tmux",
			},
		},
		StoreLocation: "store",
	}
	store, err := LoadStore(config)
	suite.Nil(err)

	suite.Equal(store.path, "store")
	suite.Equal(store.name, "Linux")
	suite.Len(store.files, 2)
	suite.True(containsFilePath(store.files, ".config/alacritty/alacritty.yml"))
	suite.True(containsFilePath(store.files, ".tmux.conf"))

	var exists bool
	exists, _ = Afs.DirExists("store/97aa776c8b768a52732c7978fd5f0af5ce5a1135")
	suite.True(exists)
	exists, _ = Afs.DirExists("store/14b4f00abd93c6222516ff054e4a9f66295d03fa")
	suite.True(exists)
}

func containsFilePath(files []*file.DotFile, path string) bool {
	for _, file := range files {
		if file.Path() == Fs.Abs(path) {
			return true
		}
	}
	return false
}

func (suite *StoreSuite) TestLoadStoreFileNotInStore() {
	config := &config.Config{
		Name: "Linux",
		WithHistory: []config.FileEntry{
			{
				Path:     ".config/alacritty/alacritty.yml",
				Mnemonic: "alacritty",
			},
		},
		WithoutHistory: []config.FileEntry{
			{
				Path:     ".tmux.conf",
				Mnemonic: "tmux",
			},
		},
		StoreLocation: "store",
	}
	Afs.RemoveAll("store/14b4f00abd93c6222516ff054e4a9f66295d03fa")

	store, err := LoadStore(config)
	suite.Nil(err)

	suite.Equal(store.path, "store")
	suite.Equal(store.name, "Linux")
	suite.Len(store.files, 2)
	suite.True(containsFilePath(store.files, ".config/alacritty/alacritty.yml"))
	suite.True(containsFilePath(store.files, ".tmux.conf"))

	var exists bool
	exists, _ = Afs.DirExists("store/97aa776c8b768a52732c7978fd5f0af5ce5a1135")
	suite.True(exists)
}

func (suite *StoreSuite) TestLoadStoreFileNotInConfig() {
	config := &config.Config{
		Name: "Linux",
		WithHistory: []config.FileEntry{
			{
				Path:     ".config/alacritty/alacritty.yml",
				Mnemonic: "alacritty",
			},
		},
		WithoutHistory: []config.FileEntry{},
		StoreLocation:  "store",
	}

	store, err := LoadStore(config)
	suite.Nil(err)

	suite.Equal(store.path, "store")
	suite.Equal(store.name, "Linux")
	suite.Len(store.files, 1)
	suite.True(containsFilePath(store.files, ".config/alacritty/alacritty.yml"))

	var exists bool
	exists, _ = Afs.DirExists("store/14b4f00abd93c6222516ff054e4a9f66295d03fa")
	suite.True(exists)
	exists, _ = Afs.DirExists("store/97aa776c8b768a52732c7978fd5f0af5ce5a1135")
	suite.False(exists)
}
