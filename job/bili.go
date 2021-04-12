package job

import (
	"fmt"
	"time"

	"github.com/Jannchie/probe-centre/db"
	"github.com/Jannchie/probe-centre/model"
	"gorm.io/gorm/clause"
)

// API const
const (
	//BiliUserInfoAPI  = "http://api.bilibili.com/x/web-interface/card?mid=%d"
	//BiliUserStatAPI  = "http://api.bilibili.com/x/space/upstat?mid=%d"
	//BiliUserVideoAPI = "http://api.bilibili.com/x/space/arc/search?mid=%d&order=pubdate&pn=1&ps=50"
	BiliVideoRankAPI = "https://api.bilibili.com/x/web-interface/ranking/v2?rid=%d&type=%s"
)

func upsertTask(task *model.Task) {
	db.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(task)
}

func CreateBiliTasks() {
	rankTypes := [...]string{"rookie", "origin", "all"}
	tidList := [...]int{0, 5, 181, 155, 119, 217, 211, 160, 188, 36, 4, 129, 3, 1, 168}
	for _, rankType := range rankTypes {
		for _, tid := range tidList {
			url := fmt.Sprintf(BiliVideoRankAPI, tid, rankType)
			task := model.Task{
				URL:      url,
				Interval: time.Hour,
			}
			upsertTask(&task)
		}
	}
}
