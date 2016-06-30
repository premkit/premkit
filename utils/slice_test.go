package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveDuplicatesNoneFound(t *testing.T) {
	a := []string{
		"1", "2", "3", "4",
	}

	a = RemoveDuplicates(a)

	assert.Equal(t, 4, len(a), "there should be 4 items in the array")
	assert.Equal(t, "1", a[0], "the first item in the array should be 1")
	assert.Equal(t, "2", a[1], "the second item in the array should be 2")
	assert.Equal(t, "3", a[2], "the third item in the array should be 3")
	assert.Equal(t, "4", a[3], "the fourth item in the array should be 4")
}

func TestRemoveDuplicatesWithDuplicates(t *testing.T) {
	a := []string{
		"1", "2", "3", "4", "2", "5",
	}

	a = RemoveDuplicates(a)

	assert.Equal(t, 5, len(a), "there should be 5 items in the array")
	assert.Equal(t, "1", a[0], "the first item in the array should be 1")
	assert.Equal(t, "2", a[1], "the second item in the array should be 2")
	assert.Equal(t, "3", a[2], "the third item in the array should be 3")
	assert.Equal(t, "4", a[3], "the fourth item in the array should be 4")
	assert.Equal(t, "5", a[4], "the fifth item in the array should be 5")
}

func TestDiffArrays(t *testing.T) {
	a := []string{"1", "2", "3"}
	b := []string{"2", "4", "6"}

	c := DiffArrays(a, b)
	assert.Equal(t, 2, len(c), "resulting array should have a length of 2")
	assert.Equal(t, "1", c[0], "first element should be '1'")
	assert.Equal(t, "3", c[1], "second element should be '3'")
}
