package login

import (
	"bufio"
	"encoding/binary"
)

const (
	MayProceed            = 0
	Timeout               = 1
	LoginSuccess          = 2
	InvalidCredentials    = 3
	DisabledAccount       = 4
	AlreadyLoggedIn       = 5
	ClientUpdated         = 6
	WorldFull             = 7
	LoginServerOffline    = 8
	LimitExceeded         = 9
	BadSessionId          = 10
	LoginServerRejected   = 11
	RequireMembersAccount = 12
	CouldNotCompleteLogin = 13
	ServerBeingUpdated    = 14
)

type LoginHandshakeResponse struct {
	LoginStatus      byte
	Unknown          int64
	ServerSessionKey int64
}

func (l *LoginHandshakeResponse) Write(writer *bufio.Writer) {
	err := binary.Write(writer, binary.BigEndian, l)
	if err != nil {
		panic(err)
	}
}
