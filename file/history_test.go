package file

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPatchRegeneration(t *testing.T) {
	const str1 = "This is the first line"
	const str2 = `This is the first line
This is the second line`
	const str3 = `This is the first line
This is the modified second line`
	const str4 = "This is the modified second line"

	history1 := NewHistory(str1)
	history2 := history1.AddCommit(str2)
	history3 := history2.AddCommit(str3)
	history4 := history3.AddCommit(str4)

	assert := assert.New(t)
	newStr4 := history4.Content()
	assert.Equal(str4, newStr4)
	newStr3 := history3.Content()
	assert.Equal(str3, newStr3)
	newStr2 := history2.Content()
	assert.Equal(str2, newStr2)
}

func TestTreeToJSON(t *testing.T) {
	assert := assert.New(t)
	tree := makeTree()
	jsonBytes := tree.ToJSON()

	var items []map[string]interface{}
	err := json.Unmarshal(jsonBytes, &items)
	assert.Equal(err, nil)
	assert.Equal(len(items), 7)

	checkString(t, items[0]["Checksum"], "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d")
	checkString(t, items[1]["Checksum"], "88fdd585121a4ccb3d1540527aee53a77c77abb8")
	checkString(t, items[2]["Checksum"], "0f1defd5135596709273b3a1a07e466ea2bf4fff")
	checkString(t, items[3]["Checksum"], "53d001f65e513a8c9560a0a40b1b823ece93204c")
	checkString(t, items[4]["Checksum"], "8f0bc65da355c6cb184de9d17bfe1baaeefbd443")
	checkString(t, items[5]["Checksum"], "afee8cb6e87492cc3d6bf07e1c387cc8845b4177")
	checkString(t, items[6]["Checksum"], "49921ef888da14b586d7498719fdbe504ba65385")
}

func checkString(t *testing.T, src interface{}, target string) {
	switch strSrc := src.(type) {
	case string:
		assert.Equal(t, target, strSrc)
	default:
		t.Fatal(fmt.Sprintf("expected string %v", src))
	}
}

func makeTree() *HistoryNode {
	root := NewHistory("hello")
	root.AddCommit("hello1")
	a := root.AddCommit("hello2")
	a.AddCommit("hello3")
	b := a.AddCommit("hello4")
	b.AddCommit("hello5")
	a.AddCommit("hello6")
	return root
}
