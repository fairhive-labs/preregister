package data

type DB interface {
	Save(u *User)
	Update(u *User)
	FindAll()
}

type mockDB struct {
}

func (db mockDB) Save(u *User)   {}
func (db mockDB) Update(u *User) {}
func (db mockDB) FindAll()       {}

var MockDB = mockDB{}
