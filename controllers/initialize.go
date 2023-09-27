package controllers

import "backend-auth/database"

type Controller struct {
	datasource *database.DB
}

func (c *Controller) SetDatasource(db *database.DB) {
	c.datasource = db
}
