package sns

import (
	"bytes"
	"log"
	"net/http"
)

func NewLocalSNS(cfg *Config) *SNS {
	return &SNS{
		client:   nil,
		topicARN: cfg.ARN,
		local:    true,
	}
}

func (s *SNS) publishLocalMsg(body []byte) error {

	resp, err := http.Post(s.topicARN, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	log.Println("[PublishLocalMsg] Msg published :", resp)

	return nil
}
