package workspace

type SecretPolicy struct {
	AllowPersist bool
}

func MaterializeSecrets(policy SecretPolicy) error {
	_ = policy
	return nil
}
