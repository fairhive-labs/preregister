package data

type DB interface {
	Save(u *User)
}
