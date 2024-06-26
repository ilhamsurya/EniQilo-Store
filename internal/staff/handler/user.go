package handler

import (
	"net/http"
	"projectsphere/eniqlo-store/internal/staff/entity"
	"projectsphere/eniqlo-store/internal/staff/service"
	"projectsphere/eniqlo-store/pkg/protocol/msg"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc service.UserService
}

func NewUserHandler(userSvc service.UserService) UserHandler {
	return UserHandler{
		userSvc: userSvc,
	}
}

func (h UserHandler) Register(c *gin.Context) {
	payload := new(entity.UserParam)

	err := c.ShouldBindJSON(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, msg.BadRequest(err.Error()))
		return
	}

	resp, err := h.userSvc.Register(c.Request.Context(), payload)
	if err != nil {
		respError := msg.UnwrapRespError(err)
		c.JSON(respError.Code, respError)
		return
	}

	c.JSON(http.StatusCreated, msg.ReturnResult("User registered successfully", resp))
}

func (h UserHandler) Login(c *gin.Context) {
	payload := new(entity.UserLoginParam)

	err := c.ShouldBindJSON(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, msg.BadRequest(err.Error()))
		return
	}

	resp, err := h.userSvc.Login(c.Request.Context(), payload)
	if err != nil {
		respError := msg.UnwrapRespError(err)
		c.JSON(respError.Code, respError)
		return
	}

	c.JSON(http.StatusOK, msg.ReturnResult("User logged successfully", resp))
}
