package service

type AuthenticateDomainService struct{}

func NewAuthenticateDomainService() *AuthenticateDomainService {
	return &AuthenticateDomainService{}
}

func (svc *AuthenticateDomainService) Login(username, password string) bool {
	return username == "admin" && password == "admin"
}

func (svc *AuthenticateDomainService) Logout() {}
