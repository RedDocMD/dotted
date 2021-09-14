package file

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type Sha = [sha1.Size]byte

type HistoryNode struct {
	content  *string
	parent   *HistoryNode // RI: (parent != nil) ^ (content != nil) == 1
	patches  []diffmatchpatch.Patch
	checksum Sha
	children []*HistoryNode
	uuid     uuid.UUID
}

// NewHistory creates a new history tree and returns
// the root node.
func NewHistory(contents string) *HistoryNode {
	sum := sha1.Sum([]byte(contents))
	uuid := uuid.New()
	return &HistoryNode{
		content:  &contents,
		parent:   nil,
		patches:  nil,
		checksum: sum,
		children: nil,
		uuid:     uuid,
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
	uuid := uuid.New()
	newNode := &HistoryNode{
		content:  nil,
		parent:   history,
		patches:  patches,
		checksum: sum,
		children: nil,
		uuid:     uuid,
	}
	history.children = append(history.children, newNode)
	return newNode
}

func (history *HistoryNode) pathFromRoot() []*HistoryNode {
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
	path := history.pathFromRoot()
	baseContent := *path[0].content
	var patches []diffmatchpatch.Patch
	for _, node := range path[1:] {
		patches = append(patches, node.patches...)
	}
	dmp := diffmatchpatch.New()
	currentContent, _ := dmp.PatchApply(patches, baseContent)
	return currentContent
}

type jsonHistoryNode struct {
	Parent   string
	Patches  string
	Checksum string
	Children []string
	Uuid     string
}

func newJsonHistoryNode(node *HistoryNode) jsonHistoryNode {
	dmp := diffmatchpatch.New()
	patches := dmp.PatchToText(node.patches)
	checksum := fmt.Sprintf("%x", node.checksum)
	children := make([]string, len(node.children))
	for i, child := range node.children {
		children[i] = child.uuid.String()
	}
	var parentUuid string
	if node.parent == nil {
		parentUuid = ""
	} else {
		parentUuid = node.parent.uuid.String()
	}
	return jsonHistoryNode{
		Parent:   parentUuid,
		Patches:  patches,
		Checksum: checksum,
		Children: children,
		Uuid:     node.uuid.String(),
	}
}

// All nodes in the sub-tree rooted at node
func (node *HistoryNode) toJsonNodes() []jsonHistoryNode {
	jsonNodes := []jsonHistoryNode{newJsonHistoryNode(node)}
	for _, child := range node.children {
		childNodes := child.toJsonNodes()
		jsonNodes = append(jsonNodes, childNodes...)
	}
	return jsonNodes
}

func (node *HistoryNode) ToJSON() []byte {
	nodes := node.toJsonNodes()
	bytes, _ := json.Marshal(nodes)
	return bytes
}
