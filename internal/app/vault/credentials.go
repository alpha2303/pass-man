package vault

type Credential struct {
	Username string
	Password string
}

func CreateCredential(username string, password string) Credential {
	return Credential{
		Username: username,
		Password: password,
	}
}
