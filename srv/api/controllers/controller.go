package apicontroller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ranjbar-dev/gowin/tools/logger"
)

type Controller struct {
}

// ===== response structs ===== //

type ErrorResponse struct {
	Message string `json:"message"`
}

// ===== responses ===== //

func (controller *Controller) ok(c *gin.Context, data any) {

	c.JSON(http.StatusOK, data)
}

func (controller *Controller) badRequest(c *gin.Context, msg string) {

	logger.Warn("Api bad response").Message("api server response was bad response").Params(map[string]any{
		"path":   c.Request.URL.Path,
		"method": c.Request.Method,
		"query":  c.Request.URL.Query(),
		"header": c.Request.Header,
	}).Log()

	c.JSON(http.StatusBadRequest, gin.H{
		"message": msg,
	})
}

func (controller *Controller) error(c *gin.Context, err error) {

	logger.Warn("Api error response").Message("api server response was error response").Params(map[string]any{
		"path":   c.Request.URL.Path,
		"method": c.Request.Method,
		"query":  c.Request.URL.Query(),
		"header": c.Request.Header,
	}).Log()

	c.JSON(http.StatusInternalServerError, gin.H{
		"message": err.Error(),
	})
}

// ===== validations ===== //

// int32, int64, string, bool
func (controller *Controller) queries(c *gin.Context, validations map[string]string) (map[string]any, bool) {

	result := make(map[string]any)

	for key, itemType := range validations {

		value := c.Query(key)
		if value == "" {

			controller.badRequest(c, "query "+key+" not found")
			return nil, false
		}

		// string type
		if itemType == "string" {

			result[key] = value
			continue
		}

		// []int32] type
		if itemType == "[]int32" {

			intList := make([]int32, 0)
			for _, v := range strings.Split(value, ",") {
				temp, err := strconv.ParseInt(v, 10, 32)
				if err != nil {
					controller.badRequest(c, "could not convert "+v+" to int32")
					return nil, false
				}
				intList = append(intList, int32(temp))
			}
			result[key] = intList
			continue
		}

		// int32 type
		if itemType == "int32" {

			temp, err := strconv.ParseInt(value, 10, 32)
			if err != nil {

				controller.badRequest(c, "could not convert "+value+" to int32")
				return nil, false
			}

			result[key] = int32(temp)
			continue
		}

		// int16 type
		if itemType == "int16" {

			temp, err := strconv.ParseInt(value, 10, 16)
			if err != nil {

				controller.badRequest(c, "could not convert "+value+" to int16")
				return nil, false
			}

			result[key] = int16(temp)
			continue
		}

		// int64 type
		if itemType == "int64" {

			temp, err := strconv.ParseInt(value, 10, 64)
			if err != nil {

				controller.badRequest(c, "could not convert "+value+" to int64")
				return nil, false
			}

			result[key] = int64(temp)
			continue
		}

		// int type
		if itemType == "int" {

			temp, err := strconv.Atoi(value)
			if err != nil {

				controller.badRequest(c, "could not convert "+value+" to int")
				return nil, false
			}

			result[key] = temp
			continue
		}

		// bool type
		if itemType == "bool" {

			temp, err := strconv.ParseBool(value)
			if err != nil {

				controller.badRequest(c, "could not convert "+value+" to bool")
				return nil, false
			}

			result[key] = temp
			continue
		}

		controller.badRequest(c, "unknown type "+itemType)
		return nil, false
	}

	return result, true
}

func NewController() *Controller {

	return &Controller{}
}
