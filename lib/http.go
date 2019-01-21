package configclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	server "github.com/cmiceli/configserver/lib"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type HTTPClient struct {
	serverLocation string
}

func NewHTTPClient(serverLocation string) server.Storage {
	return &HTTPClient{serverLocation: serverLocation}
}

func (h *HTTPClient) Get(identifier string) (server.Config, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", h.serverLocation, identifier))
	if err != nil {
		log.Printf("HTTPClient::Get - %v", err)
		return server.Config{}, err
	}
	defer resp.Body.Close()
	cfg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("HTTPClient::Get - %v", err)
		return server.Config{}, err
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Code(%s) with message '%s'", resp.StatusCode, string(cfg)))
		log.Printf("HTTPClient::Get - %v", err)
		return server.Config{}, err
	}
	var c server.Config
	err = json.Unmarshal(cfg, &c)
	if err != nil {
		log.Printf("HTTPClient::Get - %v", err)
		return server.Config{}, err
	}
	return c, nil
}

func (h *HTTPClient) Set(identifier string, cfg server.Config) error {
	buf, err := json.Marshal(cfg)
	if err != nil {
		log.Printf("HTTPClient::Set - %v", err)
		return err
	}
	resp, err := http.Post(fmt.Sprintf("%s/%s", h.serverLocation, identifier), "application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Printf("HTTPClient::Set - %v", err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("HTTPClient::Set - %v", err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Code(%s) with message '%s'", resp.StatusCode, string(body)))
		log.Printf("HTTPClient::Set - %v", err)
		return err
	}
	return nil
}

func (h *HTTPClient) LastUpdate(identifier string) (time.Time, error) {
	client := &http.Client{}
	req, err := http.NewRequest("OPTIONS", fmt.Sprintf("%s/%s", h.serverLocation, identifier), nil)
	if err != nil {
		log.Printf("HTTPClient::LastUpdate - %v", err)
		return time.Now(), err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("HTTPClient::LastUpdate - %v", err)
		return time.Now(), err
	}
	log.Printf("Was returned %v", resp)
	defer resp.Body.Close()
	cfg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("HTTPClient::LastUpdate - %v", err)
		return time.Now(), err
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Code(%s) with message '%s'", resp.StatusCode, string(cfg)))
		log.Printf("HTTPClient::LastUpdate - %v", err)
		return time.Now(), err
	}
	log.Printf("Recieved: %s", string(cfg))
	var t time.Time
	err = (&t).UnmarshalText(cfg)
	if err != nil {
		log.Printf("HTTPClient::LastUpdate - %v", err)
		return time.Now(), err
	}
	return t, nil

}
