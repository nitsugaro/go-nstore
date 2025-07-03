package test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/nitsugaro/go-nstore"
)

type Password struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type User struct {
	*nstore.Metadata

	Name     string    `json:"name"`
	LastName string    `json:"last_name"`
	Password *Password `json:"password"`
}

func TestMain(t *testing.T) {
	storage, err := nstore.New[*User]("managed/users")
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	randomName := uuid.NewString()
	storage.Save(&User{Name: randomName, LastName: "Romero", Password: &Password{Value: "1234", Type: "Hash"}})
	storage.LoadFromDisk()
	results, total := storage.Query(func(t *User) bool {
		return t.Name == randomName
	}, 3)
	if total != 1 || results[0] == nil {
		t.Errorf("user not returned in query for name '%s'", randomName)
	}
}
