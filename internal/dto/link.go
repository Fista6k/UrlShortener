package dto

type RequestLink struct {
	OriginalUrl string `form: "original_url" binding:"required"`
}

type ResponseLink struct {
	ShortUrl string `json: "short_url"`
	Message  string `json: "message"`
}
