package cache

type NoCache struct {
}

func initializeNoCache() *NoCache {
	return new(NoCache)
}

func (cache *NoCache) SaveAccessRefreshTokens(userID uint, deviceID, accessToken, refreshToken string) error {
	return nil
}
