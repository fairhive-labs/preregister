package data

type DB interface {
	Save(u *User) error
	Count() (map[string]int, error)
	List(options ...int) ([]*User, error)
}
