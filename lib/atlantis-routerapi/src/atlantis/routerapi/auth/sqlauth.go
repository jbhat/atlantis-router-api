package auth

type SqlAuth struct {
	Name string
}

func (*SqlAuth) SimpleAuth(user, secret string) error {

	return nil

}

func (*SqlAuth) IsSuperUser(user, secret string) error {

	return nil
}

func (s *SqlAuth) AuthOk() bool {
	return true
}

func GetSqlAuthorizer() *SqlAuth {

	return &SqlAuth{Name: "Default"}
}
