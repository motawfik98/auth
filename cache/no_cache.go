package cache

type NoCache struct {
}

func initializeNoCache() *NoCache {
	return new(NoCache)
}

func (cache *NoCache) SaveAccessRefreshTokens(userID uint, deviceID, accessToken, refreshToken string) error {
	return nil
}

func (cache *NoCache) MarkRefreshTokenAsUsed(refreshToken string) (int64, error) {
	return 0, nil
}
