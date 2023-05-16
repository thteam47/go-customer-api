package defaultcomponent

import (
	"strings"

	"github.com/thteam47/common-libs/confg"
	"github.com/thteam47/common-libs/mongoutil"
	"github.com/thteam47/common/entity"
	"github.com/thteam47/common/pkg/mongorepository"
	"github.com/thteam47/go-customer-api/errutil"
	"github.com/thteam47/go-customer-api/pkg/models"
)

type TenantRepository struct {
	config         *TenantRepositoryConfig
	baseRepository *mongorepository.BaseRepository
}

type TenantRepositoryConfig struct {
	MongoClientWrapper *mongoutil.MongoClientWrapper `mapstructure:"mongo-client-wrapper"`
}

func NewTenantRepositoryWithConfig(properties confg.Confg) (*TenantRepository, error) {
	config := TenantRepositoryConfig{}
	err := properties.Unmarshal(&config)
	if err != nil {
		return nil, errutil.Wrap(err, "Unmarshal")
	}

	mongoClientWrapper, err := mongoutil.NewBaseMongoClientWrapperWithConfig(properties.Sub("mongo-client-wrapper"))
	if err != nil {
		return nil, errutil.Wrap(err, "NewBaseMongoClientWrapperWithConfig")
	}
	return NewTenantRepository(&TenantRepositoryConfig{
		MongoClientWrapper: mongoClientWrapper,
	})
}

func NewTenantRepository(config *TenantRepositoryConfig) (*TenantRepository, error) {
	inst := &TenantRepository{
		config: config,
	}

	var err error
	inst.baseRepository, err = mongorepository.NewBaseRepository(&mongorepository.BaseRepositoryConfig{
		MongoClientWrapper: inst.config.MongoClientWrapper,
		Prototype:          models.Tenant{},
		MongoIdField:       "Id",
		IdField:            "TenantId",
	})
	if err != nil {
		return nil, errutil.Wrap(err, "mongorepository.NewBaseRepository")
	}

	return inst, nil
}
func (inst *TenantRepository) FindAll(userContext entity.UserContext, findRequest *entity.FindRequest) ([]models.Tenant, error) {
	result := []models.Tenant{}
	err := inst.baseRepository.FindAll(userContext, findRequest, &mongorepository.FindOptions{}, &result)
	if err != nil {
		return nil, errutil.Wrap(err, "baseRepository.FindAll")
	}
	return result, nil
}

func (inst *TenantRepository) Count(userContext entity.UserContext, findRequest *entity.FindRequest) (int32, error) {
	result, err := inst.baseRepository.Count(userContext, findRequest, &mongorepository.FindOptions{})
	if err != nil {
		return 0, errutil.Wrap(err, "baseRepository.Count")
	}
	return int32(result), nil
}

func (inst *TenantRepository) FindById(userContext entity.UserContext, id string) (*models.Tenant, error) {
	result := &models.Tenant{}
	err := inst.baseRepository.FindOneByAttribute(userContext, "TenantId", id, &mongorepository.FindOptions{}, &result)
	if err != nil {
		return nil, errutil.Wrap(err, "baseRepository.FindOneByAttribute")
	}
	return result, nil
}

func (inst *TenantRepository) FindByDomain(userContext entity.UserContext, domain string) (*models.Tenant, error) {
	result := &models.Tenant{}
	domain = strings.ToLower(strings.TrimSpace(domain))
	err := inst.baseRepository.FindOneByAttribute(userContext, "Domain", domain, &mongorepository.FindOptions{}, &result)
	if err != nil {
		return nil, errutil.Wrap(err, "baseRepository.FindOneByAttribute")
	}
	return result, nil
}

func (inst *TenantRepository) Create(userContext entity.UserContext, data *models.Tenant) (*models.Tenant, error) {
	if data.Meta == nil {
		data.Meta = map[string]string{}
	}
	data.Domain = strings.ToLower(strings.TrimSpace(data.Domain))
	err := inst.baseRepository.Create(userContext, data, nil)
	if err != nil {
		return nil, errutil.Wrap(err, "baseRepository.Create")
	}
	return data, nil
}

func (inst *TenantRepository) Update(userContext entity.UserContext, data *models.Tenant, updateRequest *entity.UpdateRequest) (*models.Tenant, error) {
	if data.Meta == nil {
		data.Meta = map[string]string{}
	}

	excludedProperties := []string{
		"CreatedTime",
	}

	err := inst.baseRepository.UpdateOneByAttribute(userContext, "TenantId", data.TenantId, data, updateRequest, &mongorepository.UpdateOptions{
		ExcludedProperties: excludedProperties,
	})
	if err != nil {
		return nil, errutil.Wrap(err, "baseRepository.UpdateOneByAttribute")
	}
	return data, nil
}
func (inst *TenantRepository) DeleteById(userContext entity.UserContext, id string) error {
	err := inst.baseRepository.DeleteOneByAttribute(userContext, "TenantId", id)
	if err != nil {
		return errutil.Wrap(err, "baseRepository.DeleteOneByAttribute")
	}
	return nil
}