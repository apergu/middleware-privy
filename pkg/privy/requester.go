package privy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/rapperror"
	"golang.org/x/net/context/ctxhttp"
)

type requester struct {
	hc *http.Client
}

func (r *requester) Do(ctx context.Context, req *http.Request, env Envelope) error {
	resp, err := ctxhttp.Do(ctx, r.hc, req)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":  "requester.Do",
				"src": "ctxhttp.Do",
			}).
			Error(err)

		return rapperror.ErrInternalServerError(
			"",
			err.Error(),
			"",
			nil,
		)
	}
	defer resp.Body.Close()

	_ = json.NewDecoder(resp.Body).Decode(env)
	// _ = json.Unmarshal(bts, env)

	logrus.
		WithFields(logrus.Fields{
			"at":  "requester.Do",
			"src": "json.NewDecoder(resp.Body).Decode(env)",
			"cde": resp.StatusCode,
		}).
		Info(fmt.Sprintf("%+v", env))

	if resp.StatusCode == http.StatusInternalServerError {
		logrus.
			WithFields(logrus.Fields{
				"at":  "requester.Do",
				"src": "http.StatusInternalServerError",
			}).
			Error(err)

		return rapperror.ErrInternalServerError(
			"",
			"Something went wrong",
			"",
			env,
		)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		logrus.
			WithFields(logrus.Fields{
				"at":         "requester.Do",
				"src":        "resp.StatusCode < 200 || resp.StatusCode > 299",
				"StatusCode": resp.StatusCode,
			}).
			Error(err)

		return rapperror.ErrInternalServerError(
			"",
			"Something went wrong",
			"",
			env,
		)
	}

	if env.Failed() != nil {
		logrus.
			WithFields(logrus.Fields{
				"at":    "requester.Do",
				"src":   "env.Failed()",
				"error": env.Failed(),
			}).
			Error(err)

		return rapperror.ErrInternalServerError(
			"",
			env.Failed().Error(),
			"",
			env,
		)
	}

	return nil
}
