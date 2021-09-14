package file

import (
	"crypto/sha1"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type Sha = [sha1.Size]byte

type HistoryNode struct {
	content  *string
	parent   *HistoryNode // RI: (parent != nil) ^ (content != nil) == 1
	patches  []diffmatchpatch.Patch
	checksum Sha
	children []*HistoryNode
}

// NewHistory creates a new history tree and returns
// the root node.
func NewHistory(contents string) *HistoryNode {
	sum := sha1.Sum([]byte(contents))
	return &HistoryNode{
		content:  &contents,
		parent:   nil,
		patches:  nil,
		checksum: sum,
		children: nil,
	}
}

// AddCommit adds a commit if necessary and returns
// the created node or nil if nothing was created.
func (history *HistoryNode) AddCommit(contents string) *HistoryNode {
	sum := sha1.Sum([]byte(contents))
	if sum == history.checksum {
		return nil
	}
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(history.Content(), contents, false)
	patches := dmp.PatchMake(diffs)
	newNode := &HistoryNode{
		content:  nil,
		parent:   history,
		patches:  patches,
		checksum: sum,
		children: nil,
	}
	history.children = append(history.children, newNode)
	return newNode
}

func (history *HistoryNode) PathFromRoot() []*HistoryNode {
	nodes := []*HistoryNode{history}
	ptr := history.parent
	for ptr != nil {
		nodes = append(nodes, ptr)
		ptr = ptr.parent
	}
	for i := 0; i < len(nodes)/2; i++ {
		tmp := nodes[i]
		nodes[i] = nodes[len(nodes)-i-1]
		nodes[len(nodes)-i-1] = tmp
	}
	return nodes
}

// Content returns the content corresponding to this node
func (history *HistoryNode) Content() string {
	if history.parent == nil {
		return *history.content
	}
	path := history.PathFromRoot()
	baseContent := *path[0].content
	var patches []diffmatchpatch.Patch
	for _, node := range path[1:] {
		patches = append(patches, node.patches...)
	}
	dmp := diffmatchpatch.New()
	currentContent, _ := dmp.PatchApply(patches, baseContent)
	return currentContent
}
