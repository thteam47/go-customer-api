package grpcapp

import (
	"context"
	"fmt"

	pb "github.com/thteam47/common/api/customer-api"
	"github.com/thteam47/common/entity"
	grpcauth "github.com/thteam47/common/grpcutil"
	"github.com/thteam47/common/pkg/adapter"
	"github.com/thteam47/common/pkg/entityutil"
	"github.com/thteam47/go-customer-api/errutil"
	"github.com/thteam47/go-customer-api/pkg/component"
	"github.com/thteam47/go-customer-api/pkg/models"
	"github.com/thteam47/go-customer-api/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CustomerService struct {
	pb.CustomerServiceServer
	componentsContainer *component.ComponentsContainer
}

func NewCustomerService(componentsContainer *component.ComponentsContainer) *CustomerService {
	return &CustomerService{
		componentsContainer: componentsContainer,
	}
}
func getTenant(item *pb.Tenant) (*models.Tenant, error) {
	if item == nil {
		return nil, nil
	}
	tenant := &models.Tenant{}
	err := util.FromMessage(item, tenant)
	if err != nil {
		return nil, errutil.Wrap(err, "FromMessage")
	}
	return tenant, nil
}

func getTenants(items []*pb.Tenant) ([]*models.Tenant, error) {
	tenants := []*models.Tenant{}
	for _, item := range items {
		tenant, err := getTenant(item)
		if err != nil {
			return nil, errutil.Wrap(err, "getTenant")
		}
		tenants = append(tenants, tenant)
	}
	return tenants, nil
}

func makeTenant(item *models.Tenant) (*pb.Tenant, error) {
	if item == nil {
		return nil, nil
	}
	tenant := &pb.Tenant{}
	err := util.ToMessage(item, tenant)
	if err != nil {
		return nil, errutil.Wrap(err, "ToMessage")
	}
	return tenant, nil
}

func makeTenants(items []models.Tenant) ([]*pb.Tenant, error) {
	tenants := []*pb.Tenant{}
	for _, item := range items {
		tenant, err := makeTenant(&item)
		if err != nil {
			return nil, errutil.Wrap(err, "makeTenant")
		}
		tenants = append(tenants, tenant)
	}
	return tenants, nil
}

func (inst *CustomerService) Create(ctx context.Context, req *pb.TenantRequest) (*pb.TenantResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "@any", "@any", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	tenant, err := getTenant(req.Data)
	if err != nil {
		return nil, errutil.Wrap(err, "getTenant")
	}
	customerID, err := entityutil.GetUserId(userContext)
	if err != nil {
		return nil, errutil.Wrap(err, "entityutil.GetUserId")
	}
	tenant.CustomerId = customerID
	result, err := inst.componentsContainer.TenantRepository().Create(userContext, tenant)
	if err != nil {
		return nil, errutil.Wrap(err, "TenantRepository.Create")
	}
	item, err := makeTenant(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeTenant")
	}
	return &pb.TenantResponse{
		Data: item,
	}, nil
}

func (inst *CustomerService) GetById(ctx context.Context, req *pb.StringRequest) (*pb.TenantResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "customer-api", "get", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	err = inst.checkPermission(userContext, req.Value)
	if err != nil {
		return nil, err
	}
	result, err := inst.componentsContainer.TenantRepository().FindById(userContext, req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "TenantRepository.FindById")
	}
	item, err := makeTenant(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeTenant")
	}
	return &pb.TenantResponse{
		Data: item,
	}, nil
}

func (inst *CustomerService) GetByDomain(ctx context.Context, req *pb.StringRequest) (*pb.TenantResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "@any", "@any", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	result, err := inst.componentsContainer.TenantRepository().FindByDomain(userContext, req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "TenantRepository.FindByDomain")
	}

	item, err := makeTenant(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeTenant")
	}
	return &pb.TenantResponse{
		Data: item,
	}, nil
}

func (inst *CustomerService) GetAll(ctx context.Context, req *pb.ListRequest) (*pb.ListTenantResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "customer-api", "get", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	findRequest, err := adapter.GetFindRequest(req, req.RequestPayload)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, fmt.Sprint(err))
	}
	if !entityutil.ServiceOrAdminRole(userContext) {
		customerID, err := entityutil.GetUserId(userContext)
		if err != nil {
			return nil, errutil.Wrap(err, "entityutil.GetUserId")
		}
		findRequest.Filters = append(findRequest.Filters, entity.FindRequestFilter{
			Key:      "CustomerId",
			Operator: entity.FindRequestFilterOperatorEqualTo,
			Value:    customerID,
		})
	}
	result, err := inst.componentsContainer.TenantRepository().FindAll(userContext, findRequest)
	if err != nil {
		return nil, errutil.Wrap(err, "TenantRepository.FindAll")
	}
	item, err := makeTenants(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeTenants")
	}
	count, err := inst.componentsContainer.TenantRepository().Count(userContext, findRequest)
	if err != nil {
		return nil, errutil.Wrap(err, "TenantRepository.Count")
	}

	return &pb.ListTenantResponse{
		Data:  item,
		Total: count,
	}, nil
}

func (inst *CustomerService) Update(ctx context.Context, req *pb.UpdateTenantRequest) (*pb.TenantResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "customer-api", "update", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}

	tenant, err := getTenant(req.Data)
	if err != nil {
		return nil, errutil.Wrap(err, "getTenant")
	}
	err = inst.checkPermission(userContext, tenant.TenantId)
	if err != nil {
		return nil, err
	}
	result, err := inst.componentsContainer.TenantRepository().Update(userContext, tenant, nil)
	if err != nil {
		return nil, errutil.Wrap(err, "TenantRepository.UpdatebyId")
	}
	item, err := makeTenant(result)
	if err != nil {
		return nil, errutil.Wrap(err, "makeTenant")
	}
	return &pb.TenantResponse{
		Data: item,
	}, nil
}
func (inst *CustomerService) DeleteById(ctx context.Context, req *pb.StringRequest) (*pb.StringResponse, error) {
	userContext, err := inst.componentsContainer.AuthService().Authentication(ctx, req.Ctx.AccessToken, req.Ctx.DomainId, "customer-api", "delete", &grpcauth.AuthenOption{})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, errutil.Message(err))
	}
	err = inst.checkPermission(userContext, req.Value)
	if err != nil {
		return nil, err
	}
	err = inst.componentsContainer.TenantRepository().DeleteById(userContext, req.Value)
	if err != nil {
		return nil, errutil.Wrap(err, "TenantRepository.DeleteById")
	}
	return &pb.StringResponse{}, nil
}
