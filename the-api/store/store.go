package store

import (
	"fmt"
	"log"
	"strings"
	"the-api/utils"

	"github.com/google/uuid"
)

type Operation struct {
	Id         string `json: "id"`
	Name       string `json: "name"`
	UpstreamId string `json: "upstreamId"`
}

// in-memory store
var operations []Operation

func init() {
	log.Printf("starting the-api store")
	reset()
}

// creates a new operation in store
//  bool success if creation was successful or false otherwise
//  Operation the newly created object
func Create(name string, upstreamId string) (bool, Operation) {
	log.Printf("creating new operation with name %s and upstreamId %s", name, upstreamId)

	if utils.IsEmpty(name) || utils.IsEmpty(upstreamId) {
		return false /*success*/, Operation{}
	}

	newOp := Operation{
		Id:         fmt.Sprintf("%v", uuid.New()),
		Name:       strings.TrimSpace(name),
		UpstreamId: strings.TrimSpace(upstreamId),
	}
	operations = append(operations, newOp)
	return true /*success*/, newOp
}

// gets all stored operations
//  []Operation all stored operations
func ReadAll() []Operation {
	log.Printf("fetching all operations")
	return operations
}

// get a stored operation by its id
//  bool success if found or false otherwise
//  Operation the operation object
func ReadById(id string) (bool, Operation) {
	log.Printf("fetching operation by id %s", id)

	for _, op := range operations {
		if op.Id == id {
			log.Printf("found operation %v", op)
			return true /*exists*/, op
		}
	}

	log.Printf("unable to find operation with id %v", id)
	return false, Operation{}
}

// resets the store to an empty state
func reset() {
	operations = make([]Operation, 0)
	log.Printf("resetting store to (len=%d, cap=%d)", len(operations), cap(operations))
}
