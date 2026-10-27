// Package responses builds every JSON body the API returns, so clients can rely
// on one envelope: {"data": ...} on success, {"error": {...}} on failure.
package responses

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"

	"leochat/app/apperr"
	"leochat/app/facades"
)

type envelope struct {
	Data any `json:"data"`
}

type errorEnvelope struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func OK(ctx http.Context, data any) http.Response {
	return ctx.Response().Json(http.StatusOK, envelope{Data: data})
}

func Created(ctx http.Context, data any) http.Response {
	return ctx.Response().Json(http.StatusCreated, envelope{Data: data})
}

func NoContent(ctx http.Context) http.Response {
	return ctx.Response().NoContent()
}

func Fail(ctx http.Context, status int, code, message string) http.AbortableResponse {
	return ctx.Response().Json(status, errorEnvelope{Error: ErrorBody{Code: code, Message: message}})
}

func BadRequest(ctx http.Context) http.Response {
	return Fail(ctx, http.StatusBadRequest, "bad_request", "Format request tidak valid.")
}

// ValidationFailed reports the first message for each invalid field.
func ValidationFailed(ctx http.Context, errs validation.Errors) http.Response {
	fields := make(map[string]string)
	for field, messages := range errs.All() {
		for _, message := range messages {
			fields[field] = message
			break
		}
	}

	return ctx.Response().Json(http.StatusUnprocessableEntity, errorEnvelope{Error: ErrorBody{
		Code:    "validation_failed",
		Message: "Data yang dikirim belum valid.",
		Fields:  fields,
	}})
}

// Error maps a service error to its HTTP response. Anything that is not an
// apperr.Error is unexpected: it is logged and hidden behind a generic 500.
func Error(ctx http.Context, err error) http.Response {
	if appErr, ok := apperr.As(err); ok {
		return Fail(ctx, statusFor(appErr.Kind), appErr.Code, appErr.Message)
	}

	facades.Log().WithContext(ctx).Error(err)
	return Fail(ctx, http.StatusInternalServerError, "internal_error", "Terjadi kesalahan di server. Coba lagi nanti.")
}

func statusFor(kind apperr.Kind) int {
	switch kind {
	case apperr.KindInvalid:
		return http.StatusUnprocessableEntity
	case apperr.KindUnauthorized:
		return http.StatusUnauthorized
	case apperr.KindForbidden:
		return http.StatusForbidden
	case apperr.KindNotFound:
		return http.StatusNotFound
	case apperr.KindConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
