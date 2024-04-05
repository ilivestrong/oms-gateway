package internal

import "time"

type Options struct {
	ListenAddressHTTPPort       string
	OrderServiceListenAddress   string
	ProductServiceListenAddress string
	TokenExpirationInSeconds    time.Time
}
