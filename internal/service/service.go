package service

import ""

type Storage interface {
	Save(link *Link) error
	FindByShortCode(code string) (*Link, error)
	FindByURL(url string) (*Link, error)
}
