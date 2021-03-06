//  Copyright (c) 2021 PingLeMe Team. All rights reserved.

package api

import (
	"PingLeMe-Backend/model"
	"PingLeMe-Backend/serializer"
	"PingLeMe-Backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CreateHomework 创建作业的接口
func CreateHomework(c *gin.Context) {
	var homeworkService service.HomeworkService
	if err := c.ShouldBind(&homeworkService); err == nil {
		homeworkService.HomeworkRepositoryInterface = &model.Repo
		res := homeworkService.CreateHomework()
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, ErrorResponse(err))
	}
}

// ViewHomework 查看作业的接口
func ViewHomework(c *gin.Context) {
	var homeworkDetailService service.HomeworkDetailService
	if err := c.ShouldBind(&homeworkDetailService); err == nil {
		homeworkDetailService.HomeworkRepositoryInterface = &model.Repo
		res := homeworkDetailService.ViewHomework(homeworkDetailService.HomeworkID)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, ErrorResponse(err))
	}
}

// ViewHomeworkList 查看作业列表的接口
func ViewHomeworkList(c *gin.Context) {
	var homeworkListService service.HomeworkListService
	if err := c.ShouldBind(&homeworkListService); err == nil {
		homeworkListService.HomeworkRepositoryInterface = &model.Repo
		homeworkListService.ClassRepositoryInterface = &model.Repo
		res := homeworkListService.ViewHomeworkList()
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, ErrorResponse(err))
	}
}

func GetHomeworkSubmissionList(c *gin.Context) {
	var service service.WorkListService
	id := c.Query("id")
	homeworkID, err := strconv.Atoi(id)
	if err == nil {
		service.WorkSubmissionRepositoryInterface = &model.Repo
		res := service.GetWorkList(homeworkID)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, serializer.ParamErr("", err))
	}
}
