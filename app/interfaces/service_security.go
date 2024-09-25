package interfaces

type ServiceSecurity interface {
	GeneratePasswordHash(password string) string
	GenerateAuthToken(login, password string) (string, error)
	ParseAuthToken(token string) (string, error)
}
