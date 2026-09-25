package application

import "context"

// DeleteLink is use-case for deleting link entry
func (s *Service) DeleteLink(ctx context.Context, id int64) error {
	return s.links.DeleteLink(ctx, id)
}
