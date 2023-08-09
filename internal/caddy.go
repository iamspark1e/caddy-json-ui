package caddy_manager_internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ref: https://github.com/qdm12/caddy-ui-server
type CaddyServer interface {
	GetCaddyConfig() (jsonContent []byte, err error)
	SetCaddyConfig(jsonContent []byte) (err error)
}

type caddyServer struct {
	caddyAPIEndpoint string
}

func NewCaddyServer(caddyAPIEndpoint string) CaddyServer {
	return &caddyServer{
		caddyAPIEndpoint: caddyAPIEndpoint,
	}
}

func DoHTTPRequest(r *http.Request) (int, []byte, error) {
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return 0, nil, err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, resBody, nil
}

func (p *caddyServer) GetCaddyConfig() (jsonContent []byte, err error) {
	r, err := http.NewRequest(http.MethodGet, p.caddyAPIEndpoint+"/config", nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")
	status, jsonContent, err := DoHTTPRequest(r)
	if err != nil {
		return nil, err
	}
	// p.logger.Info("Caddy (get config) responded HTTP status %d with content: %s", status, string(jsonContent))
	if status != http.StatusOK {
		return nil, fmt.Errorf("HTTP status code %d", status)
	}
	return jsonContent, nil
}

func (p *caddyServer) SetCaddyConfig(jsonContent []byte) (err error) {
	r, err := http.NewRequest(http.MethodPost, p.caddyAPIEndpoint+"/load", bytes.NewBuffer(jsonContent))
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/json")
	status, jsonContent, err := DoHTTPRequest(r)
	if err != nil {
		return err
	}
	// p.logger.Info("Caddy (set config) responded HTTP status %d with content: %s", status, string(jsonContent))
	if status == http.StatusOK {
		return nil
	}
	response := struct {
		Error string `json:"error"`
	}{}
	if err := json.Unmarshal(jsonContent, &response); err != nil {
		return err
	}
	if response.Error != "" {
		resErr := errors.New(response.Error)
		return resErr
	}
	return nil
}
