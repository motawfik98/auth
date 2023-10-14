package controllers

import (
	"backend-auth/cache"
	"backend-auth/database"
	"backend-auth/messaging"
)

type Controller struct {
	datasource *database.DB
	cache      *cache.Cache
	messaging  *messaging.Messaging
}

func (c *Controller) SetDatasource(db *database.DB) {
	c.datasource = db
}

func (c *Controller) SetCache(cache *cache.Cache) {
	c.cache = cache
}

func (c *Controller) SetMessaging(messaging *messaging.Messaging) {
	c.messaging = messaging
}
