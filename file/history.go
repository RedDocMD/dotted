package file

import (
	"crypto/sha1"
	"fmt"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type Sha = [sha1.Size]byte

type History struct {
	first         string
	last          string
	patches       []string
	checksums     []Sha
	firstChecksum Sha
}

func NewHistory(contents string) History {
	sum := sha1.Sum([]byte(contents))
	return History{
		first:         contents,
		last:          contents,
		patches:       nil,
		checksums:     nil,
		firstChecksum: sum,
	}
}

// AddCommit adds a commit if necessary and returns whether
// it added a commit.
func (history *History) AddCommit(contents string) bool {
	if contents == history.last {
		return false
	}
	sum := sha1.Sum([]byte(contents))
	history.checksums = append(history.checksums, sum)
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(history.last, contents, false)
	patch := dmp.PatchMake(diffs)
	history.patches = append(history.patches, dmp.PatchToText(patch))
	history.last = contents
	return true
}

func (history *History) GetContentAtChecksum(sum Sha) (string, error) {
	if sum == history.firstChecksum {
		return history.first, nil
	}
	for i, historySum := range history.checksums {
		if historySum == sum {
			patchStrings := history.patches[:i+1]
			dmp := diffmatchpatch.New()
			var patches []diffmatchpatch.Patch
			for _, str := range patchStrings {
				newPatches, _ := dmp.PatchFromText(str)
				patches = append(patches, newPatches...)
			}
			newContent, _ := dmp.PatchApply(patches, history.first)
			return newContent, nil
		}
	}
	return "", fmt.Errorf("failed to find %x in history", sum)
}
