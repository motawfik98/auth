package controllers

import (
	"backend-auth/pkg/cache"
	"backend-auth/pkg/database"
	"backend-auth/pkg/messaging"
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
