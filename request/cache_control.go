package request

type CacheControlType string

const CacheControlEphemeral CacheControlType = "ephemeral"

type CacheControl struct {
	Type CacheControlType `json:"type"`
	TTL  string           `json:"ttl,omitempty"`
}

type CacheOption func(*CacheControl)

func CacheTTL(ttl string) CacheOption {
	return func(cacheControl *CacheControl) {
		cacheControl.TTL = ttl
	}
}
