package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type humioConfig struct {
	Token    string
	Parser   string
	Endpoint string `default:"https://cloud.humio.com"`
}

/*
https://docs.humio.com/api/ingest-api/
  {
    "type": "accesslog",
    "fields": {
      "host": "webhost1"
    },
    "messages": [
       "192.168.1.21 - user1 [02/Nov/2017:13:48:26 +0000] \"POST /humio/api/v1/dataspaces/humio/ingest HTTP/1.1\" 200 0 \"-\" \"useragent\" 0.015 664 0.015",
       "192.168.1.49 - user1 [02/Nov/2017:13:48:33 +0000] \"POST /humio/api/v1/dataspaces/developer/ingest HTTP/1.1\" 200 0 \"-\" \"useragent\" 0.014 657 0.014",
       "192.168.1..21 - user2 [02/Nov/2017:13:49:09 +0000] \"POST /humio/api/v1/dataspaces/humio HTTP/1.1\" 200 0 \"-\" \"useragent\" 0.013 565 0.013",
       "192.168.1.54 - user1 [02/Nov/2017:13:49:10 +0000] \"POST /humio/api/v1/dataspaces/humio/queryjobs HTTP/1.1\" 200 0 \"-\" \"useragent\" 0.015 650 0.015"
    ]
  }
*/
type humioMsg struct {
	// Type The parser Humio will use to parse the messages
	Type string `json:"type,omitempty"`
	// Fields Annotate each of the messages with these key-values. Values must be strings.
	Fields map[string]string `json:"fields,omitempty"`
	// Tags Annotate each of the messages with these key-values as Tags. Please see other documentation on tags before using.
	Tags map[string]interface{} `json:"tags,omitempty"`

	// Messages	The raw strings representing the events. Each string will be parsed by the parser specified by type.
	Messages []string `json:"messages,omitempty"`
}

func send(log logrus.FieldLogger, out *humioMsg) (int, error) {
	bs, err := json.Marshal([]*humioMsg{out})
	if err != nil {
		return 0, errors.Wrap(err, "Failed to marshal json")
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bs))
	if err != nil {
		return 0, errors.Wrap(err, "Failed to make a new request object")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Humio.Token))
	req.Header.Set("Accept", "text/plain")

	rsp, err := client.Do(req)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to make the request")
	}
	defer rsp.Body.Close()
	log.WithField("status_code", rsp.StatusCode).Debug("Got response from humio")

	if rsp.StatusCode != http.StatusOK {
		val, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return rsp.StatusCode, errors.Wrap(err, "Failed to read response body")
		}
		log.WithField("status_code", rsp.StatusCode).Warnf("Error from humio: %s", string(val))
	}

	return rsp.StatusCode, nil
}
