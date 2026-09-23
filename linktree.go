package sinac

import "context"

// LinktreeService groups the Linktree endpoints.
type LinktreeService struct{ http *httpCore }

// Profile is GET /linktree/profiles/{handle}.
func (s *LinktreeService) Profile(ctx context.Context, handle string) (*LinktreeProfile, error) {
	r, err := linktreeProfile(handle)
	return getJSON[LinktreeProfile](ctx, s.http, r, err)
}
