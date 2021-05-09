package collector

import (
	"fmt"
	"net/http"

	. "github.com/jannchie/probe/centre/common"
	. "github.com/jannchie/probe/centre/common/model"
	"github.com/jannchie/probe/centre/common/util"
	. "github.com/jannchie/probe/centre/common/util"

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
	rawData := RawData{
		Data: form.Data,
		URL:  form.URL,
	}
	res := DB.Create(&rawData)
	if ShouldReturn(c, res.Error) {
		return
	}
	go DB.Delete(&Task{ID: form.TaskID})
	ReturnOK(c)
}

func GetRawData(c *gin.Context) {
	var form struct {
		Path    string `form:"path"`
		Count   int    `form:"count"`
		Consume bool   `form:"consume"`
	}
	err := c.ShouldBindQuery(&form)
	if ShouldReturn(c, err) {
		return
	}

	var res []RawData
	result := DB.Debug().Limit(form.Count).Find(&res, "url LIKE ?", fmt.Sprintf("%s%%", form.Path))
	if util.ShouldReturn(c, result.Error) {
		return
	}

	if form.Consume {
		var idList = make([]uint64, len(res))
		for idx, raw := range res {
			idList[idx] = raw.ID
		}
		DB.Debug().Delete(&res, "id IN ?", idList)
	}

	c.JSON(http.StatusOK, res)
}

func Init(engine *gin.Engine) {
	_ = DB.AutoMigrate(&RawData{})
	engine.Group("/raw").POST("", PostRawDataHandle).GET("", GetRawData)
}
