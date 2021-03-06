package command

import (
	"bytes"
	"command/http/arguments_builder"
	"command/http/errors"
	"command/http/types"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"test_runner_meta/interfaces"
	"time"
)

type Command struct {
	req       *http.Request
	arguments types.Arguments
	settings  types.Settings
}

func New(settings types.Settings) interfaces.Command {
	return &Command{settings: settings}
}

func (c *Command) Run(arguments string) (map[string]interface{}, error) {
	c.arguments = arguments_builder.Build(arguments)
	err := c.createRequest()
	if err != nil {
		return nil, err
	}

	responseBody, err := c.makeRequest()
	if err != nil {
		return nil, err
	}

	return c.decodeResponseBody(responseBody)
}

func (c *Command) createRequest() error {
	url, body, err := c.getURLAndBody()
	if err != nil {
		return err
	}

	c.req, err = http.NewRequest(c.settings.GetMethod(), url, body)
	if err != nil {
		return err
	}

	c.setDefaultHeaders()
	c.setHeadersFromSettings()
	c.setCookiesFromSettings()

	return nil
}

func (c Command) getURLAndBody() (string, io.Reader, error) {
	switch {
	case c.settings.ShouldPassArgumentsInURL():
		if c.arguments.IsSlashSeparated() {
			return c.buildURLForSlashSeparatedArguments(), nil, nil
		} else {
			url, err := c.tryBuildURLForAmpersandSeparatedArguments()

			return url, nil, err
		}

	default:
		return c.buildDefaultURL(), bytes.NewBufferString(c.arguments.Value()), nil
	}
}

func (c Command) buildURLForSlashSeparatedArguments() string {
	return fmt.Sprintf(
		"%s/%s/%s",
		c.settings.GetBaseURL(),
		c.settings.GetEndpoint(),
		c.arguments.Value(),
	)
}

func (c Command) tryBuildURLForAmpersandSeparatedArguments() (string, error) {
	ampersandSeparatedArguments, err := c.arguments.AmpersandSeparated()

	return fmt.Sprintf(
		"%s/%s?%s",
		c.settings.GetBaseURL(),
		c.settings.GetEndpoint(),
		ampersandSeparatedArguments,
	), err
}

func (c Command) buildDefaultURL() string {
	return fmt.Sprintf(
		"%s/%s",
		c.settings.GetBaseURL(),
		c.settings.GetEndpoint(),
	)
}

func (c Command) setDefaultHeaders() {
	c.req.Header.Set("Content-Type", "application/json; charset=utf-8")
}

func (c Command) setHeadersFromSettings() {
	for key, value := range c.settings.GetHeaders() {
		c.req.Header.Add(key, value)
	}
}

func (c Command) setCookiesFromSettings() {
	for _, cookie := range c.settings.GetCookies() {
		c.req.AddCookie(cookie)
	}
}

func (c Command) makeRequest() (io.ReadCloser, error) {
	res, err := (&http.Client{
		Timeout: time.Second * time.Duration(c.settings.GetTimeout()),
	}).Do(c.req)
	if err != nil {
		return nil, err
	}

	if 200 <= res.StatusCode && res.StatusCode <= 299 {
		return res.Body, nil
	} else {
		return nil, errors.UnsuccessfulStatus(res.StatusCode)
	}
}

func (c Command) decodeResponseBody(data io.ReadCloser) (map[string]interface{}, error) {
	responseBody, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
