package echoutils

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/thaitanloi365/goutils/errs"
	"github.com/thaitanloi365/goutils/gormpaging"
)

var userKey = "user"
var userTokenKey = "user_token"

// JwtClaims claims
type JwtClaims struct {
	jwt.StandardClaims
}

// CustomContext custom echo context
type CustomContext struct {
	echo.Context
	DB *gorm.DB
}

// GetUserFromContext get user from request context
func (c *CustomContext) GetUserFromContext(i interface{}) error {
	if reflect.ValueOf(i).Kind() == reflect.Ptr {
		return fmt.Errorf("Input must be a pointer")
	}

	var rawUser = c.Get(userKey)
	data, err := jsoniter.Marshal(&rawUser)
	if err != nil {
		return err
	}

	err = jsoniter.Unmarshal(data, i)
	if err != nil {
		return err
	}

	return nil
}

// GetJwtClaims get jwtClaims
func (c CustomContext) GetJwtClaims() (JwtClaims, error) {
	var user = c.Get(userTokenKey).(*jwt.Token)
	var claims = user.Claims.(*JwtClaims)
	var jwtClaims JwtClaims

	var err = claims.Valid()
	if err != nil {
		if e, ok := err.(*echo.HTTPError); ok {

			if e.Code == http.StatusBadRequest {
				return jwtClaims, errs.ErrTokenMissing
			}

			if e.Code == http.StatusUnauthorized {
				return jwtClaims, errs.ErrTokenInvalid
			}
		}

		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return jwtClaims, errs.ErrTokenInvalid
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return jwtClaims, errs.ErrTokenExpired
			} else {
				return jwtClaims, errs.ErrTokenInvalid
			}
		}

		return jwtClaims, errs.ErrTokenInvalid
	}

	jwtClaims.Id = claims.Id
	jwtClaims.Audience = claims.Audience
	jwtClaims.Issuer = claims.Issuer

	return jwtClaims, nil
}

// Success respond success
func (c CustomContext) Success(i interface{}) error {
	var code = http.StatusOK
	return c.JSON(code, i)
}

// BindAndValidate bind and validate input
func (c CustomContext) BindAndValidate(i interface{}) error {
	var err = c.Bind(i)
	if err != nil {
		return err
	}

	err = c.Validate(i)
	if err != nil {
		return err
	}
	return nil
}

// GetPathParamUint get param from context
func (c CustomContext) GetPathParamUint(tag string, fallbackValue ...uint) uint {
	v, err := strconv.ParseUint(c.Param(tag), 10, 32)
	if err != nil && len(fallbackValue) > 0 {
		return fallbackValue[0]
	}

	return uint(v)
}

// GetPathParamString get param from context
func (c CustomContext) GetPathParamString(tag string) string {
	return c.Param(tag)
}

// GetQueryParamString get param as string from context
func (c CustomContext) GetQueryParamString(tag string, fallbackValue ...string) string {
	v := c.QueryParam(tag)
	if v == "" && len(fallbackValue) > 0 {
		return fallbackValue[0]
	}
	return v
}

// GetQueryParamInt get param as int from context
func (c CustomContext) GetQueryParamInt(tag string, fallbackValue ...int) int {
	v, err := strconv.Atoi(c.QueryParam(tag))
	if err != nil && len(fallbackValue) > 0 {
		return fallbackValue[0]
	}
	return v
}

// GetQueryParamBool get param as int from context
func (c CustomContext) GetQueryParamBool(tag string, fallbackValue ...bool) bool {
	v, err := strconv.ParseBool(c.QueryParam(tag))
	if err != nil && len(fallbackValue) > 0 {
		return fallbackValue[0]
	}
	return v
}

// Paging get param as int from context
func (c CustomContext) Paging(p *gormpaging.PaginationParam, result interface{}) error {
	var response = gormpaging.Paging(p, result)
	return c.JSON(http.StatusCreated, response)
}
