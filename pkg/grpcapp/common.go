package grpcapp

import (
	"net/http"

	"github.com/thteam47/common/entity"
	"github.com/thteam47/common/pkg/entityutil"
	"github.com/thteam47/go-customer-api/errutil"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (inst *CustomerService) checkPermission(userContext entity.UserContext, tenantId string) error {
	if !entityutil.IsServiceToken(userContext) {
		tenant, err := inst.componentsContainer.TenantRepository().FindById(userContext, tenantId)
		if err != nil {
			return errutil.Wrap(err, "TenantRepository.FindById")
		}
		if tenant == nil {
			return errutil.Wrap(err, "Tenant not found")
		}
		customerID, err := entityutil.GetUserId(userContext)
		if err != nil {
			return errutil.Wrap(err, "entityutil.GetUserId")
		}
		if tenant.CustomerId != customerID {
			return status.Errorf(codes.PermissionDenied, http.StatusText(http.StatusForbidden))
		}
	}
	return nil
}
