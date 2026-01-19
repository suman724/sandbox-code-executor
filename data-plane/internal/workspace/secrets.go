package workspace

import "errors"

type SecretPolicy struct {
	AllowPersist bool
	Secrets      map[string]string
}

func MaterializeSecrets(policy SecretPolicy) error {
	if len(policy.Secrets) == 0 {
		return nil
	}
	if !policy.AllowPersist {
		return errors.New("secret persistence denied")
	}
	return nil
}
