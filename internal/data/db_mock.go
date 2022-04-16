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

func (db mockDB) Count() (map[string]int, error) {
	m := map[string]int{
		"advisor":     0,
		"agent":       2,
		"client":      0,
		"contributor": 0,
		"investor":    0,
		"mentor":      1,
		"talent":      3,
	}
	return m, nil
}

var MockDB = mockDB{}

type mockErrDB struct {
}

func (db mockErrDB) Save(u *User) (err error) {
	m := fmt.Sprintf("🔥 Error saving User [ %v ] in DB\n", *u)
	fmt.Print(m)
	return errors.New(m)
}

func (db mockErrDB) Count() (map[string]int, error) {
	m := fmt.Sprintf("🔥 Error counting Users in DB\n")
	fmt.Print(m)
	return nil, errors.New(m)
}

var MockErrDB = mockErrDB{}
