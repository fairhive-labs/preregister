package data

type DB interface {
	Save(u *User)
}

type mockDB struct {
}

func (db mockDB) Save(u *User) {}

var MockDB = mockDB{}
