package runtime

import "errors"

func RegisterProxy(serviceID string) (string, error) {
	if serviceID == "" {
		return "", errors.New("missing service id")
	}
	return "http://proxy/" + serviceID, nil
}

func RevokeProxy(serviceID string) error {
	if serviceID == "" {
		return errors.New("missing service id")
	}
	return nil
}
