package auth_service

import "hr-api/models"

type Auth struct {
	Username string
	Password string
}

func (a *Auth) Check() (int, error) {
	return models.CheckAuth(a.Username, a.Password)
}
