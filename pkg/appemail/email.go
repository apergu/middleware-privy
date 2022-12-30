package appemail

import (
	"context"

	"github.com/sirupsen/logrus"
	"gitlab.com/rteja-library3/remailer"
)

type AppDummyEmail struct{}

func (a AppDummyEmail) Send(ctx context.Context, messages remailer.Message) (err error) {
	logrus.
		WithFields(logrus.Fields{
			"at":      "AuthRegisterUsecaseGeneral.Register",
			"src":     "userRepo.Create",
			"to":      messages.To,
			"from":    messages.From,
			"subject": messages.Subject,
		}).
		Info(string(messages.Body.Message))

	return
}
