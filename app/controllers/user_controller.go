package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gkazioka/clientsapi/app/domain"
	"github.com/gkazioka/clientsapi/app/middlewares"
	"github.com/gkazioka/clientsapi/app/services"
	"github.com/gkazioka/clientsapi/app/types"
	"github.com/gkazioka/clientsapi/app/utils"
)

type UserController struct {
	UserService services.UserService
}

func (uc *UserController) APIStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"requestAmount": middlewares.GetRequestAmount(),
		"uptime":        utils.GetUptime(),
	})
}

func (uc *UserController) ListClients(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, uc.UserService.ListAll(c))
}

func (uc *UserController) AddClients(c *gin.Context) {
	var newClient domain.User
	if err := c.BindJSON(&newClient); err != nil {
		return
	}
	error := uc.UserService.Save(c, newClient)
	if error == nil {
		c.IndentedJSON(http.StatusCreated, newClient)
	} else {
		if error == types.ErrorAlreadyExists {
			c.IndentedJSON(http.StatusConflict, gin.H{"message": "Client already exists"})
		}
		if errors.Is(error, types.ErrorInvalidDocument) ||
			errors.Is(error, types.ErrorInvalidCpf) ||
			errors.Is(error, types.ErrorInvalidCnpj) {
			fmt.Fprintln(os.Stdout, "Invalid document error.")
			c.AbortWithStatusJSON(http.StatusBadRequest, error.Error())
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, error)
		}
	}
}

func (uc *UserController) FindByCode(c *gin.Context) {
	userFound, error := uc.UserService.FindUserByCode(c, c.Param("code"))
	if errors.Is(error, types.ErrorInvalidDocument) ||
		errors.Is(error, types.ErrorInvalidCpf) ||
		errors.Is(error, types.ErrorInvalidCnpj) {
		fmt.Fprintln(os.Stdout, "Invalid document error.")
		c.AbortWithStatusJSON(http.StatusBadRequest, error.Error())
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, error)
	}

	if userFound == nil {
		c.Status(http.StatusNotFound)
	} else {
		c.IndentedJSON(http.StatusOK, userFound)
	}
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{UserService: userService}
}
