package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/aerosystems/nix-trainee-5-6-7-8/internal/helpers"
	"github.com/labstack/echo/v4"
)

type LoginRequestBody struct {
	Email    string `json:"email" xml:"email" example:"example@gmail.com"`
	Password string `json:"password" xml:"password" example:"P@ssw0rd"`
}

// Login godoc
// @Summary login user by credentionals
// @Description Password should contain:
// @Description - minimum of one small case letter
// @Description - minimum of one upper case letter
// @Description - minimum of one digit
// @Description - minimum of one special character
// @Description - minimum 8 characters length
// @Tags users
// @Accept  json
// @Accept  xml
// @Produce application/json
// @Produce application/xml
// @Param login body handlers.LoginRequestBody true "raw request body"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Router /users/login [post]
func (h *BaseHandler) Login(c echo.Context) error {
	var requestPayload LoginRequestBody

	if err := c.Bind(&requestPayload); err != nil {
		return WriteResponse(c, http.StatusBadRequest, NewErrorPayload(err))
	}

	addr, err := helpers.ValidateEmail(requestPayload.Email)
	if err != nil {
		err = errors.New("email is not valid")
		return WriteResponse(c, http.StatusBadRequest, NewErrorPayload(err))
	}

	email := helpers.NormalizeEmail(addr)

	// Minimum of one small case letter
	// Minimum of one upper case letter
	// Minimum of one digit
	// Minimum of one special character
	// Minimum 8 characters length
	err = helpers.ValidatePassword(requestPayload.Password)
	if err != nil {
		return WriteResponse(c, http.StatusBadRequest, NewErrorPayload(err))
	}

	// validate against database
	user, err := h.userRepo.FindByEmail(email)
	if err != nil {
		return WriteResponse(c, http.StatusBadRequest, NewErrorPayload(err))
	}

	if !user.IsActive {
		err := errors.New("user has did not confirm registration")
		return WriteResponse(c, http.StatusBadRequest, NewErrorPayload(err))
	}

	valid, err := h.userRepo.PasswordMatches(user, requestPayload.Password)
	if err != nil || !valid {
		err := errors.New("invalid credentials")
		return WriteResponse(c, http.StatusBadRequest, NewErrorPayload(err))
	}

	// create pair of JWT tokens
	ts, err := h.tokensRepo.CreateToken(user.ID)
	if err != nil {
		fmt.Println("sdsdsd")
		return WriteResponse(c, http.StatusBadRequest, NewErrorPayload(err))
	}

	// add refresh token UUID to cache
	err = h.tokensRepo.CreateCacheKey(user.ID, ts)
	if err != nil {
		return WriteResponse(c, http.StatusBadRequest, NewErrorPayload(err))
	}

	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}

	payload := Response{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", requestPayload.Email),
		Data:    tokens,
	}
	return WriteResponse(c, http.StatusAccepted, payload)
}