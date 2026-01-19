package runtime

import "errors"

type DependencyPolicy struct {
	Allowlist []string
	Requested []string
}

func ValidateDependencies(policy DependencyPolicy) error {
	allowed := map[string]struct{}{}
	for _, item := range policy.Allowlist {
		allowed[item] = struct{}{}
	}
	for _, req := range policy.Requested {
		if _, ok := allowed[req]; !ok {
			return errors.New("dependency not allowlisted")
		}
	}
	return nil
}
