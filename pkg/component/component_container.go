package component

import (
	grpcauth "github.com/thteam47/common/grpcutil"
	"github.com/thteam47/go-customer-api/errutil"
	"github.com/thteam47/go-customer-api/pkg/db"
)

type ComponentsContainer struct {
	tenantRepository TenantRepository
	authService      *grpcauth.AuthInterceptor
	handler          *db.Handler
	userService      UserService
}

func NewComponentsContainer(componentFactory ComponentFactory) (*ComponentsContainer, error) {
	inst := &ComponentsContainer{}

	var err error
	inst.authService = componentFactory.CreateAuthService()
	inst.tenantRepository, err = componentFactory.CreateTenantRepository()
	inst.userService, err = componentFactory.CreateUserService()
	if err != nil {
		return nil, errutil.Wrap(err, "CreateUserService")
	}
	if err != nil {
		return nil, errutil.Wrap(err, "CreateTenantRepository")
	}
	return inst, nil
}

func (inst *ComponentsContainer) AuthService() *grpcauth.AuthInterceptor {
	return inst.authService
}

func (inst *ComponentsContainer) TenantRepository() TenantRepository {
	return inst.tenantRepository
}

func (inst *ComponentsContainer) UserService() UserService {
	return inst.userService
}

var errorCodeBadRequest = 400
