package file

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/pkg/errors"
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
		patches:  []diffmatchpatch.Patch{},
		checksum: sum,
		children: []*HistoryNode{},
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
		children: []*HistoryNode{},
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
	if sha1.Sum([]byte(currentContent)) != history.checksum {
		fmt.Fprintf(os.Stderr, "checksum of file at history %s doesn't match", history.uuid)
		os.Exit(1)
	}
	return currentContent
}

func (node *HistoryNode) NodeWithUUID(uuid string) *HistoryNode {
	if node.uuid.String() == uuid {
		return node
	}
	for _, child := range node.children {
		subNode := child.NodeWithUUID(uuid)
		if subNode != nil {
			return subNode
		}
	}
	return nil
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

func decodeJsonHistoryNode(node jsonHistoryNode, parent *HistoryNode, content *string) (*HistoryNode, error) {
	dmp := diffmatchpatch.New()
	patches, err := dmp.PatchFromText(node.Patches)
	if err != nil {
		return nil, err
	}
	uuid, err := uuid.Parse(node.Uuid)
	if err != nil {
		return nil, err
	}
	checksum, err := parseChecksum(node.Checksum)
	if err != nil {
		return nil, err
	}
	newNode := &HistoryNode{
		content:  content,
		parent:   parent,
		children: []*HistoryNode{},
		patches:  patches,
		checksum: checksum,
		uuid:     uuid,
	}
	return newNode, nil
}

func hexDigitToDecimal(digit byte) (uint8, error) {
	if digit >= '0' && digit <= '9' {
		return digit - '0', nil
	} else if digit >= 'a' && digit <= 'z' {
		return 10 + digit - 'a', nil
	} else if digit >= 'A' && digit <= 'Z' {
		return 10 + digit - 'A', nil
	} else {
		return 0, fmt.Errorf("invalid hex digit")
	}
}

func parseChecksum(str string) (Sha, error) {
	var sum Sha
	if len(str) != 40 {
		return sum, fmt.Errorf("failed to parse SHA1 sum: invalid checksum length: %d, expected 40", len(str))
	}
	for i := 0; i < len(str); i += 2 {
		first, err := hexDigitToDecimal(str[i])
		if err != nil {
			return sum, errors.Wrap(err, "failed to parse SHA1 sum")
		}
		second, err := hexDigitToDecimal(str[i+1])
		if err != nil {
			return sum, errors.Wrap(err, "failed to parse SHA1 sum")
		}
		value := first*16 + second
		sum[i/2] = value
	}
	return sum, nil
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
	bytes, err := json.Marshal(nodes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to convert history node to JSON: %v\n", node)
		os.Exit(1)
	}
	return bytes
}

func FromJSON(data []byte, content string) (*HistoryNode, error) {
	var jsonNodes []jsonHistoryNode
	err := json.Unmarshal(data, &jsonNodes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode history")
	}
	jsonNodesMap := make(map[string]jsonHistoryNode)
	var rootJsonNode jsonHistoryNode
	for _, node := range jsonNodes {
		jsonNodesMap[node.Uuid] = node
		if node.Parent == "" {
			rootJsonNode = node
		}
	}
	rootNode, err := decodeJsonHistoryNode(rootJsonNode, nil, &content)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode node")
	}
	stack := []*HistoryNode{rootNode}
	for len(stack) != 0 {
		ptr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		jsonNode := jsonNodesMap[ptr.uuid.String()]
		for _, childUuid := range jsonNode.Children {
			childJsonNode := jsonNodesMap[childUuid]
			childNode, err := decodeJsonHistoryNode(childJsonNode, ptr, nil)
			if err != nil {
				return nil, errors.Wrap(err, "failed to decode node")
			}
			ptr.children = append(ptr.children, childNode)
			stack = append(stack, childNode)
		}
	}
	return rootNode, nil
}
