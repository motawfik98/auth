package cronjobs

import (
	"backend-auth/pkg/database"
)

type CronJob struct {
	datasource *database.DB
}

func (cj *CronJob) SetDatasource(db *database.DB) {
	cj.datasource = db
}
