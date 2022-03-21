package crypto

type JWT struct {
	secret string
}

func (j *JWT) Create(a, e, t string) string {
	return j.secret + "something"
}
