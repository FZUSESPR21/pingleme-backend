package service

import (
	"PingLeMe-Backend/model"
	"PingLeMe-Backend/serializer"
)

type TeamDetailService struct {
	model.UserRepositoryInterface
	model.TeamRepositoryInterface
}

type TeamDetail struct {
	TeamID    uint       `json:"team_id"`
	TeamName  string     `json:"team_name"`
	Teammates []TeamMate `json:"teammates"`
}

type TeamMate struct {
	ID   uint   `json:"id"`
	UID  string `json:"uid"`
	Name string `json:"name"`
}

func (service *TeamDetailService) GetTeamDetail(TeamID uint) serializer.Response {
	teammates, err := service.GetTeammates(TeamID)
	teamInfo, err1 := service.GetTeam(TeamID)
	if err != nil {
		return serializer.ServerInnerErr("", err)
	}
	if err1 != nil {
		return serializer.ServerInnerErr("", err1)
	}

	teamMates := make([]TeamMate, 0)
	for _, teammate := range teammates {
		teamMates = append(teamMates, TeamMate{
			ID:   teammate.ID,
			UID:  teammate.UID,
			Name: teammate.UserName,
		})
	}

	return serializer.Response{
		Code: 0,
		Data: TeamDetail{
			TeamID:    teamInfo.ID,
			TeamName:  teamInfo.Name,
			Teammates: teamMates,
		},
		Msg: "Success",
	}
}
