package domain

type LinkRepository interface {
	Save(*Link) error
	FindByShortCode(code string) (*Link, error)
	FindByURL(url string) (*Link, error)
}