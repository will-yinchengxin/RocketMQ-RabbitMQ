package dto

import (
	_const "rock/const"
	"strconv"
)

type Data struct {
	ProjectId string `json:"projectId" form:"projectId"`
	MapId     int    `json:"mapId" form:"mapId"`
	MakeUrl   string `json:"makeUrl" form:"makeUrl"`
}

func GetData() []Data {
	ProducerData := []Data{}
	for i := 0; i < 5; i++ {
		data := Data{
			ProjectId: _const.ProjectId + strconv.Itoa(i),
			MapId:     i,
			MakeUrl:   "test" + strconv.Itoa(i),
		}
		ProducerData = append(ProducerData, data)
	}
	return ProducerData
}
