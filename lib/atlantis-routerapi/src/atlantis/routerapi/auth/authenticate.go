package auth

type Authorizer interface {
	SimpleAuth(user, secret string) error
	IsSuperUser(user, secret string) error
	AuthOk() bool //is the auth ok (connection to auth db/ldap healthy, etc)
}

func IsAllowed(user, secret string) error {
	return GetAuthorizer().SimpleAuth(user, secret)
}
