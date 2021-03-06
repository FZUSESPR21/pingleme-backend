//  Copyright (c) 2021 PingLeMe Team. All rights reserved.

package service

import (
	"PingLeMe-Backend/model"
	"PingLeMe-Backend/serializer"
	"PingLeMe-Backend/util"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"os"
)

type StudentImportService struct {
	model.UserRepositoryInterface
	model.ClassRepositoryInterface
}

type ImportMessage struct {
	TotalRows   int                 `json:"total_rows"`
	SuccessRows int                 `json:"success_rows"`
	FailedRows  int                 `json:"failed_rows"`
	ErrorRecord map[int]ErrorRecord `json:"error_record"`
}

type ErrorRecord struct {
	RowCnt       int      `json:"row_cnt"`
	RowUID       string   `json:"row_uid"`
	ErrRowUID    bool     `json:"err_row_uid"`
	RowName      string   `json:"row_name"`
	ErrRowName   bool     `json:"err_row_name"`
	RowClass     string   `json:"row_class"`
	ErrRowClass  bool     `json:"err_row_class"`
	RowPasswd    string   `json:"row_passwd"`
	ErrRowPasswd bool     `json:"err_row_passwd"`
	ErrMsg       []string `json:"err_msg"`
	IsRowIllegal bool     `json:"row_illegal"`
}

// Import 导入学生
// Excel 格式
// 学号  姓名  班级  密码
func (service *StudentImportService) Import(filepath string) serializer.Response {
	util.Log().Debug(filepath)
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			util.Log().Error(err.Error())
		}
	}(filepath)

	file, err1 := excelize.OpenFile(filepath)
	if err1 != nil {
		return serializer.ServerInnerErr("excelize error", err1)
	}

	rows, err2 := file.GetRows("Sheet1")
	if err2 != nil {
		return serializer.ServerInnerErr("excelize error", err2)
	}

	tmpMap := make(map[string]model.Class)
	errMsgs := make(map[int]ErrorRecord)

	totalRow := 0
	successRow := 0
	failedRow := 0

	for index, row := range rows {
		flag := true
		if index == 0 {
			continue
		}
		totalRow = totalRow + 1

		if len(row) < 3 {
			e := ErrorRecord{
				IsRowIllegal: true,
			}
			errMsgs[index+1] = e
		}

		var class model.Class
		if _, ok := tmpMap[row[2]]; !ok {
			c, err := service.GetClassByName(row[2])
			if err != nil {
				if e, ok := errMsgs[index+1]; ok {
					if !e.ErrRowClass {
						e.ErrRowClass = true
						e.ErrMsg = append(e.ErrMsg, err.Error())
					}
				} else {
					failedRow = failedRow + 1
					e := ErrorRecord{
						RowCnt:       index + 1,
						RowUID:       row[0],
						ErrRowUID:    false,
						RowName:      row[1],
						ErrRowName:   false,
						RowClass:     row[2],
						ErrRowClass:  true,
						RowPasswd:    row[3],
						ErrRowPasswd: false,
						ErrMsg:       make([]string, 0),
					}
					e.ErrMsg = append(e.ErrMsg, err.Error())
					errMsgs[index+1] = e
					continue
				}
			} else {
				tmpMap[row[2]] = c
				class = c
			}
		} else {
			class = tmpMap[row[2]]
		}

		user := model.User{
			UID:      row[0],
			UserName: row[1],
			Role:     model.RoleStudent,
		}

		password := "12345678"
		if len(row) == 4 {
			password = row[3]
		}

		if err := user.SetPassword(password); err != nil {
			flag = false
			util.Log().Error(err.Error())
			if e, ok := errMsgs[index+1]; ok {
				if !e.ErrRowPasswd {
					e.ErrRowPasswd = true
					e.ErrMsg = append(e.ErrMsg, err.Error())
				}
			} else {
				e := ErrorRecord{
					RowCnt:       index + 1,
					RowUID:       row[0],
					ErrRowUID:    false,
					RowName:      row[1],
					ErrRowName:   false,
					RowClass:     row[2],
					ErrRowClass:  false,
					RowPasswd:    password,
					ErrRowPasswd: true,
					ErrMsg:       make([]string, 0),
				}
				e.ErrMsg = append(e.ErrMsg, err.Error())
				errMsgs[index+1] = e
			}
		}
		if err := service.SetUser(&user); err != nil {
			flag = false
			util.Log().Error(err.Error())
			if e, ok := errMsgs[index+1]; ok {
				e.ErrMsg = append(e.ErrMsg, err.Error())
			} else {
				e := ErrorRecord{
					RowCnt:       index + 1,
					RowUID:       row[0],
					ErrRowUID:    false,
					RowName:      row[1],
					ErrRowName:   false,
					RowClass:     row[2],
					ErrRowClass:  false,
					RowPasswd:    password,
					ErrRowPasswd: false,
					ErrMsg:       make([]string, 0),
				}
				e.ErrMsg = append(e.ErrMsg, err.Error())
				errMsgs[index+1] = e
			}
		} else {
			if err := service.AddStudent(class, user); err != nil {
				flag = false
				util.Log().Error(err.Error())
				if e, ok := errMsgs[index+1]; ok {
					e.ErrMsg = append(e.ErrMsg, err.Error())
				} else {
					e := ErrorRecord{
						RowCnt:       index + 1,
						RowUID:       row[0],
						ErrRowUID:    false,
						RowName:      row[1],
						ErrRowName:   false,
						RowClass:     row[2],
						ErrRowClass:  false,
						RowPasswd:    password,
						ErrRowPasswd: false,
						ErrMsg:       make([]string, 0),
					}
					e.ErrMsg = append(e.ErrMsg, err.Error())
					errMsgs[index+1] = e
				}
			}
		}

		if flag {
			successRow = successRow + 1
		} else {
			failedRow = failedRow + 1
		}
	}
	return serializer.Response{
		Code: 0,
		Data: ImportMessage{
			TotalRows:   totalRow,
			SuccessRows: successRow,
			FailedRows:  failedRow,
			ErrorRecord: errMsgs,
		},
	}
}
