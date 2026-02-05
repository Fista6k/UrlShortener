package dto

type CreateShortLinkInput struct {
	Url string `json: url`
}

type CreateShortLinkOutput struct {
	ShortUrl    string `json: short_url`
	OriginalUrl string `json: url`
}
