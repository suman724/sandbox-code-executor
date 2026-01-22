package runtime

import "github.com/google/uuid"

func generateSessionToken() string {
	return uuid.NewString()
}
