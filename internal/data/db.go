package data

type DB interface {
	Save(u *User)
	Update(u *User)
	FindAll()
}

type MokeDB struct {
}

func (db MokeDB) Save(u *User)   {}
func (db MokeDB) Update(u *User) {}
func (db MokeDB) FindAll()       {}
