package service

import (
	"github.com/Fista6k/Url-Shorterer.git/internal/domain"
	"github.com/Fista6k/Url-Shorterer.git/internal/dto"
)

func (s ShortererService) Redirect(shortCode string) (dto.CreateShortLinkOutput, error) {
	var output dto.CreateShortLinkOutput
	link, err := s.storage.FindByShortCode(shortCode)
	if err != nil {
		return output, domain.ErrNotFound
	}
	return dto.CreateShortLinkOutput{
		OriginalUrl: string(link.OriginalUrl),
		ShortUrl:    string(link.ShortUrl),
	}, nil
}
