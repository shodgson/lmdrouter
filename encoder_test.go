package lmdrouterv2

import (
	"errors"
	"net/http"
	"testing"

	"github.com/jgroeneveld/trial/assert"
)

func TestHandleError(t *testing.T) {
	t.Run("Handle an HTTPError", func(t *testing.T) {
		res, _ := HandleError(HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid input",
		})
		assert.Equal(t, http.StatusBadRequest, res.StatusCode, "status code must be correct")
		assert.Equal(t, `{"code":400,"message":"Invalid input"}`, res.Body, "body must be correct")
	})

	t.Run("Handle an HTTPError when ExposeServerErrors is true", func(t *testing.T) {
		ExposeServerErrors = true
		res, _ := HandleError(HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "database down",
		})
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode, "status code must be correct")
		assert.Equal(t, `{"code":500,"message":"database down"}`, res.Body, "body must be correct")
	})

	t.Run("Handle an HTTPError when ExposeServerErrors is false", func(t *testing.T) {
		ExposeServerErrors = false
		res, _ := HandleError(HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "database down",
		})
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode, "status code must be correct")
		assert.Equal(t, `{"code":500,"message":"Internal Server Error"}`, res.Body, "body must be correct")
	})

	t.Run("Handle a general error when ExposeServerErrors is true", func(t *testing.T) {
		ExposeServerErrors = true
		res, _ := HandleError(errors.New("database down"))
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode, "status code must be correct")
		assert.Equal(t, `{"code":500,"message":"database down"}`, res.Body, "body must be correct")
	})

	t.Run("Handle a general error when ExposeServerErrors is false", func(t *testing.T) {
		ExposeServerErrors = false
		res, _ := HandleError(errors.New("database down"))
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode, "status code must be correct")
		assert.Equal(t, `{"code":500,"message":"Internal Server Error"}`, res.Body, "body must be correct")
	})
}
