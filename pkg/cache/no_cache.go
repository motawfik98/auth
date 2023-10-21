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
	return 1, nil
}

func (cache *NoCache) IsUsedRefreshToken(refreshToken string) (bool, error) {
	return false, nil
}

func (cache *NoCache) MarkRefreshTokensAsCompromised(refreshTokens []string) error {
	return nil
}

func (cache *NoCache) IsCompromisedRefreshToken(refreshTokens string) (bool, error) {
	return false, nil
}
