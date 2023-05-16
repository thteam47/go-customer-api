package component

import (
	"github.com/thteam47/common/entity"
)

type UserService interface {
	FindByLoginName(userContext entity.UserContext, name string) (*entity.User, error)
	FindById(userContext entity.UserContext, id string) (*entity.User, error)
	Create(userContext entity.UserContext, user *entity.User) (*entity.User, error)
	FindByEmail(userContext entity.UserContext, email string) (*entity.User, error)
	//VerifyUser(userContext entity.UserContext, id string) error
	// GetAll(userContext entity.UserContext, number int32, limit int32) ([]*models.User, error)
	// Count(userContext entity.UserContext) (int32, error)
	// GetOneByAttr(userContext entity.UserContext, data map[string]string) (*models.User, error)
	// Create(userContext entity.UserContext, user *models.User) (*models.User, error)
	// UpdatebyId(userContext entity.UserContext, user *models.User, id string) (*models.User, error)
	// DeleteById(userContext entity.UserContext, id string) error
}
