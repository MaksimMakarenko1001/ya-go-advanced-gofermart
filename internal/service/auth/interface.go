package auth

type JwtRepository interface {
	JwtValidate(tokenString string) (ok bool, err error)
}
