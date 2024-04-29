package vault

type Credentials struct {
	username string
	password string
}

func CreateCredentials(username string, password string) Credentials {
	return Credentials{
		username: username,
		password: password,
	}
}
