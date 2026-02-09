package adapter

import "github.com/Fista6k/Url-Shorterer.git/internal/domain"

func (s storage) Save(link *domain.Link) error {
	s = append(s, link)
	return nil
}

func (s storage) FindByShortCode(code string) (*domain.Link, error) {
	for _, v := range s {
		if string(v.ShortUrl) == code {
			return v, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (s storage) FindByURL(url string) (*domain.Link, error) {
	for _, v := range s {
		if string(v.OriginalUrl) == url {
			return v, nil
		}
	}
	return nil, domain.ErrNotFound
}
