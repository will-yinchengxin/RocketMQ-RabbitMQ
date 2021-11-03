package dto

import (
	"strconv"
)

var ProducerData []Data

type Data struct {
	ProjectId int  `json:"projectId" form:"projectId"`
	MapId     int `json:"mapId" form:"mapId"`
	MakeUrl   string `json:"makeUrl" form:"makeUrl"`
}

func GetData() []Data {
	for i := 0; i < 5; i++ {
		data := Data{
			ProjectId: i,
			MapId: i,
			MakeUrl: "test" + strconv.Itoa(i),
		}
		ProducerData = append(ProducerData, data)
	}
	return ProducerData
}