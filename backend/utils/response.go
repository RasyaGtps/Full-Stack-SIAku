package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct melakukan validasi terhadap struct
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// HandleValidationError menangani error validasi dan mengirim response
func HandleValidationError(c *gin.Context, err error) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)
		for _, e := range validationErrors {
			errors[e.Field()] = getValidationMessage(e)
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": errors,
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
	})
}

// getValidationMessage memberikan pesan error yang user-friendly
func getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return e.Field() + " is required"
	case "min":
		return e.Field() + " must be at least " + e.Param() + " characters"
	case "max":
		return e.Field() + " must not exceed " + e.Param() + " characters"
	case "email":
		return e.Field() + " must be a valid email"
	case "numeric":
		return e.Field() + " must be numeric"
	default:
		return e.Field() + " is invalid"
	}
}

// Response helpers
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

func CreatedResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"error":   message,
	})
}

// GetPageParam gets page parameter from query string
func GetPageParam(c *gin.Context) int {
	page := 1
	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}
	return page
}

// GetLimitParam gets limit parameter from query string
func GetLimitParam(c *gin.Context) int {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 && val <= 100 {
			limit = val
		}
	}
	return limit
}