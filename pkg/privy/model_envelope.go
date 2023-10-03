package privy

type Envelope interface {
	Created() bool
	Failed() error
}

type MainEnvelope struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (m MainEnvelope) Created() bool {
	return true
}

func (m MainEnvelope) Failed() error {
	return nil
}

type TopupEnvelope struct {
	MainEnvelope
	Data TopupCreateResponse `json:"data"`
}
