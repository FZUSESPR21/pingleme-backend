//  Copyright (c) 2021 PingLeMe Team. All rights reserved.

package api

import (
	"PingLeMe-Backend/model"
	"PingLeMe-Backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// FillInPerformance 填写绩效的接口
func FillInPerformance(c *gin.Context) {
	var service service.TeamPerformanceService
	if err := c.ShouldBind(&service); err == nil {
		service.PerformanceRepositoryInterface = &model.Repo
		res := service.FillInPerformance()
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusOK, ErrorResponse(err))
	}
}
