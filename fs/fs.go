package fs

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/afero"
)

type Fs interface {
	afero.Fs

	UserHomeDir() string
	Join(components ...string) string
	IsAbs(path string) bool
	Abs(path string) string
}

// Filesystem while working on OS
var OsFs = NewWrappedOsFs()
var OsAfs = afero.Afero{Fs: OsFs}

// Filesystem while testing
var MockFs = NewWrappedMockFs()
var MockAfs = afero.Afero{Fs: MockFs}

type WrappedOsFs struct{ afero.OsFs }
type WrappedMockFs struct{ afero.MemMapFs }

func NewWrappedOsFs() Fs {
	return &WrappedOsFs{}
}

func NewWrappedMockFs() Fs {
	return &WrappedMockFs{}
}

func (fs WrappedOsFs) UserHomeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to find user home directory")
		os.Exit(1)
	}
	return dir
}

func (fs WrappedMockFs) UserHomeDir() string {
	return "/home/dknite"
}

func (fs WrappedOsFs) Join(components ...string) string {
	return filepath.Join(components...)
}

func (fs WrappedMockFs) Join(components ...string) string {
	return path.Join(components...)
}

func (fs WrappedOsFs) IsAbs(path string) bool {
	return filepath.IsAbs(path)
}

func (fs WrappedMockFs) IsAbs(pathstr string) bool {
	return path.IsAbs(pathstr)
}

func (fs WrappedOsFs) Abs(path string) string {
	dir, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to make path absolute")
		os.Exit(1)
	}
	return dir
}

func (fs WrappedMockFs) Abs(path string) string {
	if !fs.IsAbs(path) {
		return fs.Join(fs.UserHomeDir(), path)
	}
	return path
}
