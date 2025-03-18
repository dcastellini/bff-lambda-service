package arg

import (
	"context"
	"github.com/dcastellini/bff-lambda-service/internal/core/domain"
	"github.com/dcastellini/lib-service/arg-api/pkg/client"

	"sync"
)

type ProductServiceAdapterOptions func(s *productServiceAdapter)

type productServiceAdapter struct {
	apiClient     client.
	options       []ProductServiceAdapterOptions
	apiClientOnce sync.Once
}

func NewProductServiceAdapter(options ...ProductServiceAdapterOptions) *productServiceAdapter {
	return &productServiceAdapter{
		options: options,
	}
}

func (a *productServiceAdapter) newAPIClient() {
	if a.apiClient != nil {
		return
	}
	a.apiClient = client.NewBillingAPIClient()
}

func (a *productServiceAdapter) CreateProduct(
	ctx context.Context,
	request *domain.CreateProductRequest,
) error {
	a.lazyInit()

	response, err := a.apiClient.CreateBillReminder(ctx, apiRequest)

	if err != nil {
		switch response.Code {
		case createBillReminderErrorCode:
			return domain.ErrCreateBillReminder
		case invalidAgendaErrorCode:
			return domain.ErrInvalidAgenda
		case adhesionAgendaErrorCode:
			return domain.ErrAdhesionAgenda
		default:
			return domain.ErrUnknownBillReminder
		}
	}
	return nil
}

func (a *billRemindersServiceAdapter) EditBillReminder(
	ctx context.Context,
	request *domain.EditBillReminderRequest,
) error {
	a.lazyInit()

	apiRequest := apiDomain.EditBillReminderRequest{
		BillReminderID: request.BillReminderID,
		Alias:          request.Alias,
		Active:         request.Active,
	}

	response, err := a.apiClient.EditBillReminder(ctx, apiRequest)
	if err != nil {
		switch response.Code {
		case editBillReminderErrorCode:
			return domain.ErrEditBillReminder
		case noAdhesionsErrorCode:
			return domain.ErrNoAdhesions
		case providerUnknownErrorCode:
			return domain.ErrUnknownProvider
		default:
			return domain.ErrUnknownBillReminder
		}
	}
	return nil
}

func (a *billRemindersServiceAdapter) DeleteBillReminder(ctx context.Context, billReminderID string) error {
	return domain.ErrServiceNotImplemented
}

func (a *billRemindersServiceAdapter) GetBillReminders(ctx context.Context, clientID string) ([]domain.BillReminder, error) {
	a.lazyInit()

	response, err := a.apiClient.GetBillReminders(ctx, clientID)
	if err != nil {
		switch response.Code {
		case getBillRemindersErrorCode:
			return nil, domain.ErrGetBillReminders
		case noAdhesionsErrorCode:
			return nil, domain.ErrNoAdhesions
		case unknownAdhesionsErrorCode:
			return nil, domain.ErrUnknownProvider
		default:
			return nil, domain.ErrUnknownBillReminder
		}

	}
	return mapBillReminders(response), nil
}

func (a *productServiceAdapter) lazyInit() {
	a.apiClientOnce.Do(func() {
		a.loadOptions()
		a.newAPIClient()
	})
}

func (a *productServiceAdapter) loadOptions() {
	for _, option := range a.options {
		option(a)
	}
}

func WithCustomProductAPIClient(apiClient client.IBillRemindersAPI) ProductServiceAdapterOptions {
	return func(a *productServiceAdapter) {
		a.apiClient = apiClient
	}
}
