package controllers

import (
	"github.com/labstack/echo/v4"
	"ptm/internal/di"
	"ptm/internal/services"
	"ptm/pkg/utils/customError"
	"ptm/pkg/utils/response"
	"strconv"
)

type AuditLogController interface {
	GetLogs(e echo.Context) error
	GetLogsById(e echo.Context) error
}

type auditLogController struct {
	service services.AuditLogService
}

type logRequest struct {
	EntityId   uint   `json:"entity_id"`
	EntityType string `json:"entity_type"`
}

func NewAuditLogController() AuditLogController {

	service := di.Resolve[services.AuditLogService]()
	return &auditLogController{
		service: service,
	}
}

func (a *auditLogController) GetLogs(e echo.Context) error {
	var req logRequest

	if err := e.Bind(&req); err != nil {
		return customError.InternalServerError("Error binding request", err)
	}

	if err := e.Validate(&req); err != nil {
		return customError.BadRequest("error validating request", err)
	}

	logs, err := a.service.GetModelLogs(req.EntityType, req.EntityId)

	if err != nil {
		return err
	}

	return response.Ok(e, "Successful", logs)
}

func (a *auditLogController) GetLogsById(e echo.Context) error {
	idString := e.QueryParam("id")

	id, err := strconv.Atoi(idString)

	if err != nil {
		return customError.BadRequest("error converting id to int", err)
	}

	log, err := a.service.FindById(uint(id))
	if err != nil {
		return err
	}

	return response.Ok(e, "Successful", log)
}
