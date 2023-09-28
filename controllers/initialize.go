package controllers

import (
	"backend-auth/cache"
	"backend-auth/database"
)

type Controller struct {
	datasource *database.DB
	cache      *cache.Cache
}

func (c *Controller) SetDatasource(db *database.DB) {
	c.datasource = db
}

func (c *Controller) SetCache(cache *cache.Cache) {
	c.cache = cache
}
