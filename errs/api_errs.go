package errs

import "net/http"

var (
	// ErrRoleInvalid invalid role
	ErrRoleInvalid = New(10000, "Role is invalid", http.StatusUnauthorized)
	// ErrTokenInvalid invalid token
	ErrTokenInvalid = New(10001, "Token is invalid", http.StatusUnauthorized)
	// ErrTokenExpired expired token
	ErrTokenExpired = New(10002, "Token is expired", http.StatusUnauthorized)
	// ErrTokenMissing missing token
	ErrTokenMissing = New(10003, "Token is missing or malformed", http.StatusUnauthorized)
	// ErrSessionTimeout session timeout
	ErrSessionTimeout = New(10004, "Session timeout", http.StatusUnauthorized)
)
