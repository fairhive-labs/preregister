package data

import "fmt"

type mockDB struct {
}

func (db mockDB) Save(u *User) {
	fmt.Printf("💾 User [ %v ] saved in DB", *u)
}

var MockDB = mockDB{}
