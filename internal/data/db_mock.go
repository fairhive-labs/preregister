package data

import (
	"errors"
	"fmt"
)

type mockDB struct {
}

func (db mockDB) Save(u *User) (err error) {
	fmt.Printf("💾 User [ %v ] saved in DB\n", *u)
	return
}

var MockDB = mockDB{}

type mockErrDB struct {
}

func (db mockErrDB) Save(u *User) (err error) {
	m := fmt.Sprintf("🔥 Error saving User [ %v ] in DB\n", *u)
	fmt.Print(m)
	return errors.New(m)
}

var MockErrDB = mockErrDB{}
