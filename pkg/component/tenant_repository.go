package component

import (
	"github.com/thteam47/common/entity"
	"github.com/thteam47/go-customer-api/pkg/models"
)

type TenantRepository interface {
	FindAll(userContext entity.UserContext, findRequest *entity.FindRequest) ([]models.Tenant, error)
	Count(userContext entity.UserContext, findRequest *entity.FindRequest) (int32, error)
	FindById(userContext entity.UserContext, id string) (*models.Tenant, error)
	FindByDomain(userContext entity.UserContext, domain string) (*models.Tenant, error)
	Create(userContext entity.UserContext, data *models.Tenant) (*models.Tenant, error)
	Update(userContext entity.UserContext, data *models.Tenant, updateRequest *entity.UpdateRequest) (*models.Tenant, error)
	DeleteById(userContext entity.UserContext, id string) error
}
