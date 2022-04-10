package store

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateOperationSucceedsWithValidName(t *testing.T) {
	actual /*success*/, op := Create("some-valid-name")
	assert.Equal(t, actual, true /*expected*/)

	exists, _ := ReadById(op.Id)
	defer reset()
	assert.True(t, exists)
}

func TestReadAllReturnsAllStoredOperations(t *testing.T) {
	operations := ReadAll()
	actualLen := len(operations)
	assert.Equal(t, actualLen, 0 /*expected len*/)

	Create("op-01")
	Create("op-02")
	Create("op-03")

	operations = ReadAll()
	defer reset()

	actualLen = len(operations)
	assert.Equal(t, actualLen, 3 /*expected len*/)
}

func TestReadByIdReturnsTheExpectedStoredOperation(t *testing.T) {
	Create("op-01")
	Create("op-02")
	_, createdOp := Create("op-03")

	exists, fetchedOp := ReadById(createdOp.Id)
	defer reset()

	assert.True(t, exists)
	assert.Equal(t, createdOp.Id, fetchedOp.Id)
	assert.Equal(t, fetchedOp.Name, "op-03")
}

func TestCreateOperationFailsWithEmptyName(t *testing.T) {
	actual /*success*/, _ := Create("")
	defer reset()
	assert.Equal(t, actual, false /*expected*/)
}
