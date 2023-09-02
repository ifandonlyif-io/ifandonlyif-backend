package api

import (
	"fmt"
	"github.com/google/uuid"
	db "github.com/ifandonlyif-io/ifandonlyif-backend/db/sqlc"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"strconv"
)

type AdminResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (server *Server) AdminLogin(e echo.Context) error {
	response := AdminResponse{
		Message: "",
		Data:    nil,
		Error:   nil,
	}

	email := e.FormValue("email")
	password := e.FormValue("password")

	if email == "" || password == "" {
		response.Message = "login param invalid"
		return e.JSON(
			http.StatusUnprocessableEntity,
			response,
		)
	}

	request := LoginRequest{
		Email:    email,
		Password: password,
	}

	fmt.Println("admin login request: ", request)

	user, err := server.store.GetAdminUserByEmail(e.Request().Context(), request.Email)
	if err != nil {
		response.Message = "login failed"
		response.Error = err.Error()
		return e.JSON(
			http.StatusUnprocessableEntity,
			response,
		)
	}

	if user.Email == "" {
		response.Message = "email or password invalid"
		return e.JSON(
			http.StatusForbidden,
			response,
		)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(request.Password),
	)

	if err != nil {
		response.Message = "email or password invalid"
		return e.JSON(
			http.StatusForbidden,
			response,
		)
	}

	response.Message = "login success"
	user.Password = ""

	idAndToken, err := server.store.GetTokenByUserId(e.Request().Context(), user.ID)
	if err != nil {
		idAndToken, err = server.store.CreateUserToken(e.Request().Context(),
			db.CreateUserTokenParams{
				UserID: user.ID,
				Token:  generateToken(),
			},
		)
		return err
	} else {
		idAndToken, _ = server.store.UpdateUserToken(e.Request().Context(), db.UpdateUserTokenParams{
			Token:  generateToken(),
			UserID: user.ID,
		})
	}

	if err != nil {
		response.Message = "Internal Server Error"
		response.Error = err.Error()
	}

	response.Data = struct {
		User  db.AdminUser `json:"user"`
		Token string       `json:"token"`
	}{
		User:  user,
		Token: idAndToken.Token,
	}

	return e.JSON(http.StatusOK, response)
}

func generateToken() string {
	return strconv.FormatInt(rand.Int63(), 36)
}

// get token from header
func (server *Server) checkLoginStatus(e echo.Context) error {
	response := AdminResponse{
		Message: "",
		Data:    nil,
		Error:   nil,
	}

	token := e.Request().Header.Get("token")
	UserId, err := server.store.GetUserIdByToken(e.Request().Context(), token)
	if err != nil {
		response.Message = "token not found"
		response.Error = err.Error()
		return e.JSON(http.StatusForbidden, response)
	}

	user, err := server.store.GetAdminUserByID(e.Request().Context(), UserId.UserID)
	if err != nil {
		response.Message = "token does not match any user"
		response.Error = err.Error()
		return e.JSON(http.StatusForbidden, response)
	}

	response.Message = "token valid"
	user.Password = ""
	response.Data = user
	return e.JSON(http.StatusOK, response)
}

func (server *Server) GetAllUsers(e echo.Context) error {
	response := AdminResponse{
		Message: "",
		Data:    nil,
		Error:   nil,
	}

	if !server.isAdmin(e) {
		response.Message = "permission denied"
		return e.JSON(http.StatusForbidden, response)
	}

	users, err := server.store.GetAllAdminUsers(e.Request().Context())
	if err != nil {
		response.Message = "Internal Server Error"
		response.Error = err.Error()
		return e.JSON(http.StatusInternalServerError, response)
	}

	for _, user := range users {
		user.Password = ""
	}

	response.Data = struct {
		User []db.AdminUser `json:"users"`
	}{
		User: users,
	}

	return e.JSON(http.StatusOK, response)
}

type newUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (server *Server) NewUser(e echo.Context) error {
	response := AdminResponse{
		Message: "",
		Data:    nil,
		Error:   nil,
	}

	if !server.isAdmin(e) {
		response.Message = "permission denied"
		return e.JSON(http.StatusForbidden, response)
	}

	var request *newUserRequest
	err := e.Bind(&request)
	if err != nil {
		response.Message = "Invalid request"
		response.Error = err.Error()
		return e.JSON(http.StatusUnprocessableEntity, response)
	}

	if request.Email == "" || request.Name == "" || request.Password == "" {
		response.Message = "Invalid request some field is empty"
		return e.JSON(http.StatusUnprocessableEntity, response)
	}

	user, _ := server.store.GetAdminUserByEmail(e.Request().Context(), request.Email)
	if user.Email != "" {
		response.Message = "Email already exist"
		return e.JSON(http.StatusUnprocessableEntity, response)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		response.Message = "Internal Server Error"
		response.Error = err.Error()
		return e.JSON(http.StatusInternalServerError, response)
	}

	user, err = server.store.CreateAdminUser(e.Request().Context(), db.CreateAdminUserParams{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashedPassword),
		IsAdmin:  false,
	})

	if err != nil {
		response.Message = "Internal Server Error"
		response.Error = err.Error()
		return e.JSON(http.StatusInternalServerError, response)
	}

	user.Password = ""

	response.Message = "User created"
	response.Data = user
	return e.JSON(http.StatusOK, response)
}

func (server *Server) DeleteUser(e echo.Context) error {
	response := AdminResponse{
		Message: "",
		Data:    nil,
		Error:   nil,
	}

	if !server.isAdmin(e) {
		response.Message = "permission denied"
		return e.JSON(http.StatusForbidden, response)
	}

	id := e.Param("id")
	Uuid, err := uuid.Parse(id)

	if err != nil {
		response.Message = "id is not valid"
		response.Error = err.Error()
		return e.JSON(http.StatusUnprocessableEntity, response)
	}

	err = server.store.DeleteAdminUser(e.Request().Context(), Uuid)

	if err != nil {
		response.Message = "Internal Server Error"
		response.Error = err.Error()
		return e.JSON(http.StatusInternalServerError, response)
	}

	response.Message = "User deleted"
	return e.JSON(http.StatusOK, response)
}

func (server *Server) isAdmin(e echo.Context) bool {
	token := e.Request().Header.Get("token")
	UserId, err := server.store.GetUserIdByToken(e.Request().Context(), token)
	if err != nil {
		return false
	}

	user, err := server.store.GetAdminUserByID(e.Request().Context(), UserId.UserID)
	if err != nil {
		return false
	}

	if !user.IsAdmin {
		return false
	}

	return true
}
