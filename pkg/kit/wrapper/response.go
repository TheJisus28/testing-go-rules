package wrapper

import (
	"errors"
	"net/http"

	kiterrors "github.com/donca/user-crud/pkg/kit/errors"

	"github.com/labstack/echo/v4"
)

type Response[T any] struct {
	Status       string `json:"status"`
	StatusReason string `json:"status_reason"`
	Data         T      `json:"data"`
}

func codeToHTTPStatus(code string) int {
	switch code {
	case kiterrors.CodeNotFound:
		return http.StatusNotFound
	case kiterrors.CodeValidation:
		return http.StatusBadRequest
	case kiterrors.CodeUnauthorized:
		return http.StatusUnauthorized
	case kiterrors.CodeForbidden:
		return http.StatusForbidden
	case kiterrors.CodeConflict, kiterrors.CodeAlreadyExists:
		return http.StatusConflict
	case kiterrors.CodeUnprocessable:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}

func FromError(err error) (Response[any], int) {
	var domainErr *kiterrors.DomainError
	if errors.As(err, &domainErr) {
		return Response[any]{
			Status:       domainErr.Code,
			StatusReason: domainErr.Message,
			Data:         nil,
		}, codeToHTTPStatus(domainErr.Code)
	}
	return Response[any]{
		Status:       kiterrors.CodeInternal,
		StatusReason: "an unexpected error occurred",
		Data:         nil,
	}, http.StatusInternalServerError
}

func GenerateResponse[T any](c echo.Context, data T, err error) error {
	if err != nil {
		resp, status := FromError(err)
		return c.JSON(status, resp)
	}
	return c.JSON(http.StatusOK, Response[T]{
		Status:       kiterrors.StatusSuccess,
		StatusReason: kiterrors.StatusReasonOK,
		Data:         data,
	})
}
