package component

import grpcauth "github.com/thteam47/common/grpcutil"

type ComponentFactory interface {
	CreateAuthService() *grpcauth.AuthInterceptor
	CreateTenantRepository() (TenantRepository, error)
	CreateUserService() (UserService, error)
}
