package service

import (
	"github.com/Fista6k/Url-Shorterer.git/internal/domain"
	"github.com/Fista6k/Url-Shorterer.git/internal/dto"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func (s ShortererService) Shorten(input dto.CreateShortLinkInput) (dto.CreateShortLinkOutput, error) {
	var output dto.CreateShortLinkOutput

	oldLink, err := s.storage.FindByURL(input.Url)
	if err != nil {
		return output, domain.ErrNotFound
	}
	if oldLink != nil {
		return dto.CreateShortLinkOutput{
			OriginalUrl: string(oldLink.OriginalUrl),
			ShortUrl:    string(oldLink.ShortUrl),
		}, nil
	}

	alphabet := "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"

	short_url, err := gonanoid.Generate(alphabet, 8)
	if err != nil {
		return output, domain.ErrInternalProblems
	}

	url, err := s.storage.FindByShortCode(short_url)
	if err != nil {
		return output, domain.ErrNotFound
	}

	for url != nil {
		short_url, err = gonanoid.Generate(alphabet, 8)
		if err != nil {
			return output, domain.ErrInternalProblems
		}

		url, err = s.storage.FindByShortCode(short_url)
		if err != nil {
			return output, domain.ErrNotFound
		}
	}

	l, err := domain.NewLink(input.Url, short_url)
	if err != nil {
		return output, err
	}

	err = s.storage.Save(l)
	if err != nil {
		return output, err
	}

	return dto.CreateShortLinkOutput{
		ShortUrl:    short_url,
		OriginalUrl: input.Url,
	}, nil
}
