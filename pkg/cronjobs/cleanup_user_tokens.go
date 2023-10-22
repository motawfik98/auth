package cronjobs

import "time"

func (cj *CronJob) CleanupUserTokens() error {
	currentTime := time.Now()
	return cj.datasource.CleanupUserTokens(currentTime)
}

func (cj *CronJob) CleanupUsedRefreshToken() error {
	currentTime := time.Now()
	return cj.datasource.CleanupUsedRefreshToken(currentTime)
}

func (cj *CronJob) CleanupGeneratedRefreshToken() error {
	currentTime := time.Now()
	return cj.datasource.CleanupGeneratedRefreshToken(currentTime)
}

func (cj *CronJob) CleanupInvalidatedRefreshToken() error {
	currentTime := time.Now()
	return cj.datasource.CleanupInvalidatedRefreshToken(currentTime)
}
