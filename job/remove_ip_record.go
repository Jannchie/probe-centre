package job

import (
	"time"

	"github.com/Jannchie/probe-centre/db"
	"github.com/Jannchie/probe-centre/model"
)

func RemoveIpRecord() {
	for range time.Tick(time.Hour) {
		const hoursOfDay = 24
		const daysOfWeek = 7
		db.DB.Delete(&model.IPRecord{}, "time < ?", time.Now().UTC().
			Add(-time.Hour*hoursOfDay*daysOfWeek))
	}
}
