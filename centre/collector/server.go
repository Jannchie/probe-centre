package collector

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jannchie/probe/common"
	"github.com/jannchie/probe/common/model"
	. "github.com/jannchie/probe/common/util"

	"github.com/gin-gonic/gin"
)

func PostRawDataHandle(c *gin.Context) {
	var form struct {
		TaskID uint64 `json:"task_id"`
		Data   string `json:"data"`
		URL    string `json:"url"`
	}
	err := c.ShouldBindJSON(&form)
	if ShouldReturn(c, err) {
		return
	}
	data := []byte(form.Data)
	url := form.URL
	taskID := form.TaskID
	if err = SaveRawData(data, url, taskID); ShouldReturn(c, err) {
		return
	}
	ReturnOK(c)
}

func SaveRawData(data []byte, url string, taskID uint64) error {
	if json.Valid(data) {
		rawData := model.RawJSONData{
			Data: data,
			URL:  url,
		}
		if err := saveRawJSONData(rawData, taskID); err != nil {
			return err
		}
	} else {
		rawData := model.RawData{
			Data: data,
			URL:  url,
		}
		if err := saveRawData(rawData, taskID); err != nil {
			return err
		}
	}
	return nil
}

func saveRawData(rawData model.RawData, taskID uint64) error {
	res := common.DB.Create(&rawData)
	if res.Error != nil {
		return res.Error
	}
	go common.DB.Delete(&model.Task{ID: taskID})
	return nil
}

func saveRawJSONData(rawData model.RawJSONData, taskID uint64) error {
	res := common.DB.Create(&rawData)
	if res.Error != nil {
		return res.Error
	}
	go common.DB.Delete(&model.Task{ID: taskID})
	return nil
}

func GetRawData(c *gin.Context) {
	var form struct {
		Path    string `form:"path"`
		Count   int    `form:"count"`
		Consume bool   `form:"consume"`
		Type    string `form:"type"`
	}
	err := c.ShouldBindQuery(&form)
	if ShouldReturn(c, err) {
		return
	}
	if form.Type == "test" {
		var res []model.RawData
		result := common.DB.Limit(form.Count).Find(&res, "url LIKE ?", fmt.Sprintf("%s%%", form.Path))
		if ShouldReturn(c, result.Error) {
			return
		}
		if form.Consume {
			var idList = make([]uint64, len(res))
			for idx, raw := range res {
				idList[idx] = raw.ID
			}
			common.DB.Debug().Delete(&res, "id IN ?", idList)
		}
		c.JSON(http.StatusOK, res)
	} else {
		var res []model.RawJSONData
		result := common.DB.Limit(form.Count).Find(&res, "url LIKE ?", fmt.Sprintf("%s%%", form.Path))
		if ShouldReturn(c, result.Error) {
			return
		}
		if form.Consume {
			var idList = make([]uint64, len(res))
			for idx, raw := range res {
				idList[idx] = raw.ID
			}
			common.DB.Debug().Delete(&res, "id IN ?", idList)
		}
		c.JSON(http.StatusOK, res)
	}
}

func Init(engine *gin.Engine) {
	_ = common.DB.AutoMigrate(&model.RawData{})
	_ = common.DB.AutoMigrate(&model.RawJSONData{})
	engine.Group("/raw").POST("", PostRawDataHandle).GET("", GetRawData)
}
