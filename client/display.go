package client

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type DisplayService service

func (s *DisplayService) ChangeLED(ctx context.Context, r, g, b int) error {
	l := s.client.logger.WithFields(log.Fields{
		"func": "DisplayService.ChangeLED",
	})

	payload := struct {
		Red   int `json:"red"`
		Blue  int `json:"blue"`
		Green int `json:"green"`
	}{
		Red:   r,
		Green: g,
		Blue:  b,
	}

	l.WithField("payload", payload).Debug("sending command")

	res, err := s.client.Post(ctx, "/api/led/change", payload)
	if err != nil {
		return errors.Wrap(err, "cannot change led color")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		s.client.logger.Errorf("DisplayService.ChangeLED: received '%d'", res.StatusCode)
		return errors.New("change led failed")
	}

	return nil
}
