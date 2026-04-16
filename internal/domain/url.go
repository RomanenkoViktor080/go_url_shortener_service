package domain

type CreateShortUrlDto struct {
	Url string `json:"url" binding:"required,url"`
}
