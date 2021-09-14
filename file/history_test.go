package file

import (
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
