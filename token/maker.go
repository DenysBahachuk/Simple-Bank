package token

import "time"

type Maker interface {
	//CreateToken creates a new token for a specific username and durtion
	CreateToken(username string, duration time.Duration) (string, error)

	//ValidateToken checks if a token is valid
	ValidateToken(token string) (*Payload, error)
}
