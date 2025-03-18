package adapters

import (
	"github.com/dcastellini/bff-lambda-service/internal/core/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (h *HTTPHandler) CreateProductHandler(ginCtx *gin.Context) {
	log := logger.GetLogger(ginCtx.Request.Context())

	registerValidation()

	request := new(domain.CreateProductRequest)

	if err := ginCtx.ShouldBind(request); err != nil {
		log.Errorf("error unmarshalling input request", err.Error())
		_ = ginCtx.Error(domain.ErrInvalidRequest)
		return
	}

	response, err := h.BillRemindersBFFService.CreateBillReminder(ginCtx.Request.Context(), request)
	if err != nil {
		_ = ginCtx.Error(err)
		return
	}

	ginCtx.JSON(http.StatusOK, response)
}

func (h *HTTPHandler) EditProductHandler(ginCtx *gin.Context) {
	log := logger.GetLogger(ginCtx.Request.Context())

	registerValidation()

	request := new(domain.EditProductRequest)

	productID := ginCtx.Param("uid")
	if strings.TrimSpace(productID) == "" {
		_ = ginCtx.Error(domain.ErrInvalidRequest)
		return
	}

	request.BillReminderID = billReminderID

	if err := ginCtx.ShouldBind(request); err != nil {
		log.Errorf("error unmarshalling request body", err.Error())
		_ = ginCtx.Error(domain.ErrInvalidRequest)
		return
	}

	response, err := h.BillRemindersBFFService.EditBillReminder(ginCtx.Request.Context(), request)
	if err != nil {
		_ = ginCtx.Error(err)
		return
	}

	ginCtx.JSON(http.StatusOK, response)
}

func (h *HTTPHandler) DeleteProductHandler(ginCtx *gin.Context) {
	log := logger.GetLogger(ginCtx.Request.Context())

	request := new(domain.DeleteProductRequest)

	if err := ginCtx.ShouldBindUri(request); err != nil {
		log.Errorf("error binding request's path param", err.Error())
		_ = ginCtx.Error(domain.ErrInvalidRequest)
		return
	}

	response, err := h.BillRemindersBFFService.DeleteBillReminder(ginCtx.Request.Context(), request)
	if err != nil {
		_ = ginCtx.Error(err)
		return
	}

	ginCtx.JSON(http.StatusOK, response)
}

func (h *HTTPHandler) GetProductsHandler(ginCtx *gin.Context) {
	log := logger.GetLogger(ginCtx.Request.Context())

	request := new(domain.GetProductsRequest)

	if err := ginCtx.ShouldBindQuery(request); err != nil {
		log.Errorf("error binding request's query params", err.Error())
		_ = ginCtx.Error(domain.ErrInvalidRequest)
		return
	}

	response, err := h.BillRemindersBFFService.GetBillReminders(ginCtx.Request.Context(), request)
	if err != nil {
		_ = ginCtx.Error(err)
		return
	}

	log.Debug("got bill reminders", log.Any(response))

	ginCtx.JSON(http.StatusOK, response)
}
