package vault

type Credential struct {
	Domain   string
	Username string
	Password string
}

func CreateCredential(domain string, username string, password string) Credential {
	return Credential{
		Domain:   domain,
		Username: username,
		Password: password,
	}
}
