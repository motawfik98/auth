package controllers

import "backend-auth/database"

type Controller struct {
	db *database.DB
}

func (c *Controller) SetDB(db *database.DB) {
	c.db = db
}
