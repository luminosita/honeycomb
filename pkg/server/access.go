package server

import "strings"

type Access int

const (
	PUBLIC Access = iota
	RESTRICTED
)

func (m Access) String() string {
	return []string{"PUBLIC", "RESTRICTED"}[m]
}

func AccessFromString(key string) Access {
	return map[string]Access{"public": PUBLIC, "restricted": RESTRICTED}[strings.ToLower(key)]
}
