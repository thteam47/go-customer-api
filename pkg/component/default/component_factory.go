package defaultcomponent

import (
	"github.com/thteam47/common-libs/confg"
	grpcauth "github.com/thteam47/common/grpcutil"
	"github.com/thteam47/common/handler"
	"github.com/thteam47/go-customer-api/errutil"
	"github.com/thteam47/go-customer-api/pkg/component"
)

type ComponentFactory struct {
	properties confg.Confg
	handle     *handler.Handler
}

func NewComponentFactory(properties confg.Confg, handle *handler.Handler) (*ComponentFactory, error) {
	inst := &ComponentFactory{
		properties: properties,
		handle:     handle,
	}

	return inst, nil
}

func (inst *ComponentFactory) CreateAuthService() *grpcauth.AuthInterceptor {
	authService := grpcauth.NewAuthInterceptor(inst.handle)
	return authService
}

func (inst *ComponentFactory) CreateTenantRepository() (component.TenantRepository, error) {
	tenantRepository, err := NewTenantRepositoryWithConfig(inst.properties.Sub("tenant-repository"))
	if err != nil {
		return nil, errutil.Wrapf(err, "NewTenantRepositoryWithConfig")
	}
	return tenantRepository, nil
}

func (inst *ComponentFactory) CreateUserService() (component.UserService, error) {
	userService, err := NewUserServiceWithConfig(inst.properties.Sub("user-service"))
	if err != nil {
		return nil, errutil.Wrapf(err, "NewUserServiceWithConfig")
	}
	return userService, nil
}