package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	BaseUrl   string = "https://api.ngrok.com"
	TunnelUrl string = BaseUrl + "/tunnels"
)

type URIID struct {
	ID  string `json:"id"`
	URI string `json:"uri"`
}

type Tunnel struct {
	ID            string `json:"id"`
	PublicURL     string `json:"public_url"`
	StartedAt     string `json:"started_at"`
	Region        string `json:"region"`
	TunnelSession URIID  `json:"tunnel_session"`
	Endpoint      URIID  `json:"endpoint"`
	ForwardsTo    string `json:"forwards_to"`
}

type TunnelResp struct {
	Tunnels []Tunnel `json:"tunnels"`
	URI     string   `json:"uri"`
	NextURI string   `json:"next_page_uri"`
}

func fetchTunnels(apiKey string) (TunnelResp, error) {

	var respData TunnelResp

	req, _ := http.NewRequest(`GET`, TunnelUrl, nil)

	// two headers
	req.Header.Set(`Authorization`, `Bearer `+apiKey)
	req.Header.Set(`ngrok-version`, `2`)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return respData, err
	}
	defer resp.Body.Close()

	jsonData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return respData, err
	}

	err = json.Unmarshal(jsonData, &respData)

	return respData, err
}
