package cache

import "fmt"

func userDeviceAccessTokenKey(userID uint, deviceID string) string {
	return fmt.Sprintf("access-token::%d::%s", userID, deviceID)
}

func userDeviceRefreshTokenKey(userID uint, deviceID string) string {
	return fmt.Sprintf("refresh-token::%d::%s", userID, deviceID)
}

func usedRefreshTokensKey() string {
	return "used-refresh-tokens"
}

func compromisedRefreshTokensKey() string {
	return "compromised-refresh-tokens"
}
