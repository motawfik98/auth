package servers

import (
	"backend-auth/pkg/cache"
	"backend-auth/pkg/database"
	"backend-auth/pkg/messaging"
)

type Server struct {
	datasource *database.DB
	cache      *cache.Cache
	messaging  *messaging.Messaging
}

func (s *Server) SetDatasource(db *database.DB) {
	s.datasource = db
}

func (s *Server) SetCache(cache *cache.Cache) {
	s.cache = cache
}

func (s *Server) SetMessaging(messaging *messaging.Messaging) {
	s.messaging = messaging
}
