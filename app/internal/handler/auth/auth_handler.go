package auth

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/gunzgo2mars/test-poke-service/app/internal/constants"
	"github.com/gunzgo2mars/test-poke-service/app/internal/core/model"
	"github.com/gunzgo2mars/test-poke-service/app/internal/core/service/auth"
	"github.com/gunzgo2mars/test-poke-service/app/pkg/response"
	"github.com/gunzgo2mars/test-poke-service/app/pkg/utils"
)

type IAuthHandler interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
}

type authHandler struct {
	authSvc auth.IAuthService

	validator utils.IValidator
}

func New(authSvc auth.IAuthService, validator utils.IValidator) IAuthHandler {
	return &authHandler{
		authSvc:   authSvc,
		validator: validator,
	}
}

func (h *authHandler) Register(c echo.Context) error {
	ctx := c.Request().Context()
	res := response.NewJSONMessage()
	var validatorErrs []string
	var req *model.AuthRequest

	if err := c.Bind(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			res.AddCode(constants.ERROR_REQUEST_CODE).AddMessage(err.Error()),
		)
	}

	errs := h.validator.SetupValidator(
		req,
		utils.Fields{
			"Username": h.validator.Rules(
				h.validator.Required,
				h.validator.Max(20),
				h.validator.Min(8),
				h.validator.Regexp(utils.AlphaNumericLowerRegex),
			),
			"Password": h.validator.Rules(
				h.validator.Required,
			),
		},
	).Validate()

	if errs != nil {
		for _, v := range errs {
			validatorErrs = append(validatorErrs, fmt.Sprintf("%s: %s", v.FieldName, v.Err.Error()))
		}

		return c.JSON(
			http.StatusBadRequest,
			response.BuildMessageWithErrors(
				validatorErrs,
				res.AddCode(constants.ERROR_REQUEST_CODE).AddMessage("Error validator"),
			),
		)
	}

	if err := h.authSvc.RegisterUser(ctx, req); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			res.AddCode(constants.ERROR_INTERNAL_CODE).AddMessage(err.Error()),
		)
	}

	return c.JSON(
		http.StatusOK,
		res.AddCode(constants.SUCCESS_CODE).AddMessage(constants.SUCCESS_MSG),
	)
}

func (h *authHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	res := response.NewJSONMessage()
	var validatorErrs []string
	var req *model.AuthRequest

	if err := c.Bind(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			res.AddCode(constants.ERROR_REQUEST_CODE).AddMessage(err.Error()),
		)
	}

	errs := h.validator.SetupValidator(
		req,
		utils.Fields{
			"Username": h.validator.Rules(
				h.validator.Required,
				h.validator.Max(20),
				h.validator.Min(8),
				h.validator.Regexp(utils.AlphaNumericLowerRegex),
			),
			"Password": h.validator.Rules(
				h.validator.Required,
			),
		},
	).Validate()

	if errs != nil {
		for _, v := range errs {
			validatorErrs = append(validatorErrs, fmt.Sprintf("%s: %s", v.FieldName, v.Err.Error()))
		}

		return c.JSON(
			http.StatusBadRequest,
			response.BuildMessageWithErrors(
				validatorErrs,
				res.AddCode(constants.ERROR_REQUEST_CODE).AddMessage("Error validator"),
			),
		)
	}

	token, err := h.authSvc.ValidatingUser(ctx, req)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			res.AddCode(constants.ERROR_INTERNAL_CODE).AddMessage(err.Error()),
		)
	}

	resp := &model.AuthResponse{
		AccessToken: token,
	}

	return c.JSON(
		http.StatusOK,
		response.BuildMessageWithData(
			resp,
			res.AddCode(constants.SUCCESS_CODE).AddMessage(constants.SUCCESS_MSG),
		),
	)
}
