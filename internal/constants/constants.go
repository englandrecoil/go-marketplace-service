package constants

import "time"

const (
	MinTitleLength                    = 1
	MaxTitleLength                    = 50
	MinDescLength                     = 10
	MaxDescLength                     = 750
	MinPrice                          = 0
	MaxPrice                          = 99999999
	MaxImageSize                      = 10 * 1024 * 1024
	TokenExpirationTime time.Duration = time.Minute * 15
	MinEntropyBits                    = 60
	MinLoginLength                    = 5
	MaxLoginLength                    = 32
)
