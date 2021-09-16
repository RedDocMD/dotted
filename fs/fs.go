package fs

import "github.com/spf13/afero"

// Filesystem while working on OS
var OsFs = afero.NewOsFs()
var OSAfs = afero.Afero{Fs: OsFs}

// Filesystem while working testing
var MockFs = afero.NewMemMapFs()
var MockAfs = afero.Afero{Fs: MockFs}
