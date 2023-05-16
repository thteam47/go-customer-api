package defaultcomponent

import (
	"context"
	"time"

	"github.com/thteam47/common-libs/confg"
	v1 "github.com/thteam47/common/api/identity-api"
	"github.com/thteam47/common/entity"
	"github.com/thteam47/go-customer-api/errutil"
	"github.com/thteam47/go-customer-api/util"
	"google.golang.org/grpc"
)

type UserService struct {
	config *UserServiceConfig
	client v1.IdentityServiceClient
}

type UserServiceConfig struct {
	Address     string        `mapstructure:"address"`
	Timeout     time.Duration `mapstructure:"timeout"`
	AccessToken string        `mapstructure:"access_token"`
}

func NewUserServiceWithConfig(properties confg.Confg) (*UserService, error) {
	config := UserServiceConfig{}
	err := properties.Unmarshal(&config)
	if err != nil {
		return nil, errutil.Wrap(err, "Unmarshal")
	}
	return NewUserService(&config)
}

func NewUserService(config *UserServiceConfig) (*UserService, error) {
	inst := &UserService{
		config: config,
	}
	conn, err := grpc.Dial(config.Address, grpc.WithInsecure())
	if err != nil {
		return nil, errutil.Wrapf(err, "grpc.Dial")
	}
	client := v1.NewIdentityServiceClient(conn)
	inst.client = client
	return inst, nil
}

func (inst *UserService) requestCtx(userContext entity.UserContext) *v1.Context {
	return &v1.Context{
		AccessToken: inst.config.AccessToken,
		DomainId:    userContext.DomainId(),
	}
}

func getUser(item *v1.User) (*entity.User, error) {
	if item == nil {
		return nil, nil
	}
	user := &entity.User{}
	err := util.FromMessage(item, user)
	if err != nil {
		return nil, errutil.Wrap(err, "FromMessage")
	}
	return user, nil
}

func getUsers(items []*v1.User) ([]*entity.User, error) {
	users := []*entity.User{}
	for _, item := range items {
		user, err := getUser(item)
		if err != nil {
			return nil, errutil.Wrap(err, "getUser")
		}
		users = append(users, user)
	}
	return users, nil
}

func makeUser(item *entity.User) (*v1.User, error) {
	if item == nil {
		return nil, nil
	}
	user := &v1.User{}
	err := util.ToMessage(item, user)
	if err != nil {
		return nil, errutil.Wrap(err, "ToMessage")
	}
	return user, nil
}

func makeUsers(items []*entity.User) ([]*v1.User, error) {
	users := []*v1.User{}
	for _, item := range items {
		user, err := makeUser(item)
		if err != nil {
			return nil, errutil.Wrap(err, "makeUser")
		}
		users = append(users, user)
	}
	return users, nil
}

func (inst *UserService) FindByLoginName(userContext entity.UserContext, name string) (*entity.User, error) {
	result, err := inst.client.GetByLoginName(context.Background(), &v1.StringRequest{
		Ctx:   inst.requestCtx(userContext),
		Value: name,
	})
	if err != nil {
		return nil, errutil.Wrapf(err, "client.GetByLoginName")
	}
	item, err := getUser(result.Data)
	if err != nil {
		return nil, errutil.Wrapf(err, "getUser")
	}
	return item, nil
}

func (inst *UserService) FindByEmail(userContext entity.UserContext, email string) (*entity.User, error) {
	result := &v1.UserResponse{}
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), inst.config.Timeout)
	defer cancel()
	result, err = inst.client.GetByEmail(ctx, &v1.StringRequest{
		Ctx:   inst.requestCtx(userContext),
		Value: email,
	})
	if err != nil {
		return nil, errutil.Wrapf(err, "client.GetByEmail")
	}
	item, err := getUser(result.Data)
	if err != nil {
		return nil, errutil.Wrapf(err, "getUser")
	}
	return item, nil
}

func (inst *UserService) FindById(userContext entity.UserContext, id string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), inst.config.Timeout)
	defer cancel()
	result, err := inst.client.GetById(ctx, &v1.StringRequest{
		Ctx:   inst.requestCtx(userContext),
		Value: id,
	})
	if err != nil {
		return nil, errutil.Wrapf(err, "client.GetById")
	}
	item, err := getUser(result.Data)
	if err != nil {
		return nil, errutil.Wrapf(err, "getUser")
	}
	return item, nil
}

func (inst *UserService) Create(userContext entity.UserContext, user *entity.User) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), inst.config.Timeout)
	defer cancel()
	item, err := makeUser(user)
	if err != nil {
		return nil, errutil.Wrapf(err, "getUser")
	}
	result, err := inst.client.Create(ctx, &v1.UserRequest{
		Ctx:  inst.requestCtx(userContext),
		Data: item,
	})
	if err != nil {
		return nil, errutil.Wrapf(err, "client.GetById")
	}
	itemApi, err := getUser(result.Data)
	if err != nil {
		return nil, errutil.Wrapf(err, "getUser")
	}
	return itemApi, nil
}
func (inst *UserService) GetAll(userContext entity.UserContext, number int32, limit int32) ([]*entity.User, error) {

	return nil, nil
}

func (inst *UserService) Count(userContext entity.UserContext) (int32, error) {

	return 0, nil
}

func (inst *UserService) GetOneByAttr(userContext entity.UserContext, data map[string]string) (*entity.User, error) {

	return nil, nil
}

func (inst *UserService) UpdatebyId(userContext entity.UserContext, user *entity.User, id string) (*entity.User, error) {

	return nil, nil
}

func (inst *UserService) DeleteById(userContext entity.UserContext, id string) error {

	return nil
}

//func (inst *UserService) VerifyUser(userContext entity.UserContext, id string) error {
//	userId := id
//	if id == "" {
//		userId = userContext.GetUserId()
//	}
//	ctx, cancel := context.WithTimeout(context.Background(), inst.config.Timeout)
//	defer cancel()
//	_, err := inst.client.ApproveUser(ctx, &v1.ApproveUserRequest{
//		Ctx:    inst.requestCtx(),
//		UserId: userId,
//		Status: "verified",
//	})
//	if err != nil {
//		return errutil.Wrapf(err, "client.GetById")
//	}
//
//	return nil
//}
