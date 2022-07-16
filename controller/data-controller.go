package controller

import (
	"fmt"
	"net/http"
	"pro1/dto"
	"pro1/entity"
	"pro1/helper"
	"pro1/service"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type DataController interface {
	All(ctx *gin.Context)
	FindById(ctx *gin.Context)
	FindByString(ctx *gin.Context)
	Insert(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type dataController struct {
	dataService service.DataService
	jwtService  service.JWTService
}

func NewDataController(dataService service.DataService, jwtService service.JWTService) DataController {
	return &dataController{
		dataService: dataService,
		jwtService:  jwtService,
	}
}

func (c *dataController) All(ctx *gin.Context) {
	var data []entity.Data = c.dataService.All()
	res := helper.BuildResponse(true, "ok", data)
	ctx.JSON(200, res)
}

func (c *dataController) FindById(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("bo param id was found", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var data entity.Data = c.dataService.FindDataByID(id)
	if (data == entity.Data{}) {
		res := helper.BuildErrorResponse("Data not found", "No data with given id", helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
	} else {
		res := helper.BuildResponse(true, "ok", data)
		ctx.JSON(200, res)
	}
}

func (c *dataController) FindByString(ctx *gin.Context){
	authHeader := ctx.GetHeader("Authorization")
	_, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken == nil{
	   s_data := ctx.Param("q")
		 if s_data == ""{
			 res := helper.BuildErrorResponse("Data field is empty", "please write something", helper.EmptyObj{})
			 ctx.AbortWithStatusJSON(http.StatusNotFound, res)
		 } else {
			 var search_data = c.dataService.FindDataByString(s_data)
			 ctx.JSON(200, search_data)
		 }
	} else {
		panic(errToken.Error())
	}

}

func (c *dataController) Insert(ctx *gin.Context) {
	var dataCreated dto.DataCreateDTO
	errDTO := ctx.ShouldBind(&dataCreated)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
	} else {
		authHeader := ctx.GetHeader("Authorization")
		userID := c.getUserIDByToken(authHeader)
		convertedUserID, err := strconv.ParseUint(userID, 10, 64)
		if err == nil {
			dataCreated.UserID = convertedUserID
		}

		result := c.dataService.InsertData(dataCreated)
		response := helper.BuildResponse(true, "ok", result)
		ctx.JSON(http.StatusCreated, response)
	}

}

func (c *dataController) Update(ctx *gin.Context) {
	var dataUpdateDTO dto.DataUpdateDTO
	errDTO := ctx.ShouldBind(&dataUpdateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.dataService.IsAllowedToEdit(userID, dataUpdateDTO.ID) {
		id, errID := strconv.ParseUint(userID, 10, 64)
		if errID == nil {
			dataUpdateDTO.UserID = id
		}
		result := c.dataService.UpdateData(dataUpdateDTO)
		response := helper.BuildResponse(true, "ok", result)
		ctx.JSON(200, response)
	} else {
		response := helper.BuildErrorResponse("You dont have permission", "You are not the owner", helper.EmptyObj{})
		ctx.JSON(http.StatusForbidden, response)
	}

}

func (c *dataController) Delete(ctx *gin.Context) {
	var data entity.Data

	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		response := helper.BuildErrorResponse("Failed tou get id", "No param id were found", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
	}

	data.ID = id

	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.dataService.IsAllowedToEdit(userID, data.ID) {
		c.dataService.DeleteData(data)
		res := helper.BuildResponse(true, "Deleted", helper.EmptyObj{})
		ctx.JSON(200, res)
	} else {
		response := helper.BuildErrorResponse("You dont have permission", "You are not the owner", helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusForbidden, response)
	}

}


func (c *dataController) getUserIDByToken(token string) string {
	aToken, err := c.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}

	claims := aToken.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id
}
