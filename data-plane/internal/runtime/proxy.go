package runtime

func RegisterProxy(serviceID string) (string, error) {
	return "http://proxy/" + serviceID, nil
}
