package service

import (
	"PingLeMe-Backend/model"
	"PingLeMe-Backend/serializer"
)

type CreateTeamService struct {
	model.TeamRepositoryInterface
	Name          string `form:"name" json:"name" binding:"required"`
	GroupLeaderID int    `form:"group_leader_id" json:"group_leader_id" binding:"required"`
	ClassID       int    `form:"class_id" json:"class_id" binding:"required"`
}

func (service *CreateTeamService) CreateTeam() serializer.Response {
	//TODO 1.创建者是否已有团队 2.班级是否存在 3.队名是否重复
	if isAllowed, _ := CheckStatus(uint(service.ClassID), "team"); !isAllowed {
		serializer.ParamErr("创建团队功能未开放", nil)
	}

	if service.UserHasTeam(uint(service.GroupLeaderID)) {
		return serializer.ParamErr("已有团队，不可重复创建", nil)
	}

	has, err := service.SetTeam(model.Team{
		Name:          service.Name,
		GroupLeaderID: service.GroupLeaderID,
		ClassID:       service.ClassID,
	})
	if err != nil {
		return serializer.DBErr("数据获取错误", err)
	}

	if has != 1 {
		return serializer.DBErr("has != 1 错误", err)
	}
	//TODO 返回Team_id
	return serializer.Response{
		Code: 0,
		Msg:  "创建团队成功",
	}
}
