package workers

import (
	"backend-auth/internal/logger"
	"encoding/json"
	"errors"
	"fmt"
)

func (w *Worker) InvalidateCompromisedRefreshTokens(encodedBody []byte) error {
	fmt.Println("Received a job")
	var body map[string]string
	err := json.Unmarshal(encodedBody, &body)
	if err != nil {
		logger.LogFailure(err, fmt.Sprintf("Failed to unmarshal the encoded body %s", encodedBody))
		return err
	}
	token, found := body["refresh_token"]
	if !found {
		err := errors.New("cannot find the refresh_token in the message")
		logger.LogFailure(
			err,
			fmt.Sprintf("refresh_token key was not found in the body: %v", body),
		)
		return err
	}
	generatedTokens := w.datasource.GetGeneratedRefreshTokenChain(token)
	generatedTokenStrings := make([]string, len(generatedTokens))
	for index, generatedToken := range generatedTokens {
		generatedTokenStrings[index] = generatedToken.RefreshToken
	}
	return w.cache.Connection.MarkRefreshTokensAsCompromised(generatedTokenStrings)
}
