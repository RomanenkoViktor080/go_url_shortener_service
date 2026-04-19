package domain

// URLRedirectDto represents the path parameters for the redirect endpoint
type URLRedirectDto struct {
	// Unique short identifier (hash) of the URL
	// in: path
	Hash string `uri:"hash" binding:"required" example:"abc123"`
}

// CreateShortUrlDto represents the request payload for creating a new short URL
type CreateShortUrlDto struct {
	// Original long URL to be shortened
	Url string `json:"url" binding:"required,url" example:"https://github.com/google/uuid"`
}

// ShortUrlDto represents the generated short URL data
type ShortUrlDto struct {
	Url string `json:"url"  example:"http://localhost:8080/asdg"`
}
