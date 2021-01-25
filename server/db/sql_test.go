package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testDsn = "root:root@tcp(localhost:3306)/testscreens"

func TestSqlDb_AddRemove(t *testing.T) {
	db, err := NewSqlDb(testDsn)
	assert.NoError(t, err)
	if db != nil {
		defer db.Close()
	}

	id, err := db.Add("expected_path", "expected_desc")
	assert.NoError(t, err)
	assert.NotEqual(t, 0, id)

	err = db.Remove(id)
	assert.NoError(t, err)
}

func TestSqlDb_FindByDesc(t *testing.T) {
	db, err := NewSqlDb(testDsn)
	assert.NoError(t, err)
	if db != nil {
		defer db.Close()
	}

	id1, _ := db.Add("expected_path", "expect to find")
	defer db.Remove(id1)
	id2, _ := db.Add("expected_path", "must not find")
	defer db.Remove(id2)

	images, err := db.FindByDesc("expect to find")

	assert.NoError(t, err)
	assert.Len(t, images, 1)
	assert.Equal(t, "expect to find", images[0].Data.Desc)
}

func TestSqlDb_All(t *testing.T) {
	db, err := NewSqlDb(testDsn)
	assert.NoError(t, err)
	if db != nil {
		defer db.Close()
	}

	id1, _ := db.Add("test1", "")
	defer db.Remove(id1)
	id2, _ := db.Add("test2", "")
	defer db.Remove(id2)
	id3, _ := db.Add("test3", "")
	defer db.Remove(id3)

	images, err := db.All()

	assert.NoError(t, err)
	assert.Len(t, images, 3)
}
