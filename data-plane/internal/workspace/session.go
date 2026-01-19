package workspace

type SessionWorkspace struct {
	SessionID string
}

func NewSessionWorkspace(sessionID string) SessionWorkspace {
	return SessionWorkspace{SessionID: sessionID}
}
