package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
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

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		return respData, err
	}

	err = json.Unmarshal(jsonData, &respData)

	return respData, err
}

func ListTunnels() ([]string, error) {
	apiKey, err := GetEnvString("ApiKey")
	if err != nil {
		return nil, err
	}
	tunnelResp, err := fetchTunnels(apiKey)
	if err != nil {
		return nil, err
	} else {
		fmt.Printf("tunnels: %+v\n", tunnelResp)
		urls := []string{}
		if len(tunnelResp.Tunnels) > 0 {
			url := tunnelResp.Tunnels[0].PublicURL
			if url != "" {
				urls = append(urls, url)
			}
		}
		return urls, nil
	}

}

/*
func FUNCNAME [TypeParam1, TypeParam2]( params ) {}

OR

type TypeParam interface {
	type1 | type2
}
*/

type NgrokTunnel struct {
	Protocal   string
	PublicURL  string
	PublicPort int
	LocalPort  int
}

type ChanRetType interface {
	NgrokTunnel | *NgrokTunnel
}

type ChanRet[T ChanRetType] struct {
	Val T
	Err error
}

func GetNgrokTunnels[T *NgrokTunnel]() <-chan ChanRet[T] {

	outChan := make(chan ChanRet[T])

	apiKey, err := GetEnvString("ApiKey")
	if err != nil {
		outChan <- ChanRet[T]{Err: err}
	}
	go func() {
		tunnelResp, err := fetchTunnels(apiKey)
		if err != nil {
			outChan <- ChanRet[T]{Err: err}
		} else {
			fmt.Printf("tunnels: %+v\n", tunnelResp)

			if len(tunnelResp.Tunnels) > 0 {
				for idx := 0; idx < len(tunnelResp.Tunnels); idx++ {
					tunnel := tunnelResp.Tunnels[idx]
					if val, err := parseTunnelData(&tunnel); err != nil {
						outChan <- ChanRet[T]{Err: err}
					} else {
						outChan <- ChanRet[T]{Val: val}
					}
				}
			}
		}

		close(outChan)
	}()

	return outChan
}

func parseTunnelData(data *Tunnel) (*NgrokTunnel, error) {
	val := &NgrokTunnel{}

	url := data.PublicURL
	idx := strings.Index(url, "://")
	if idx > 0 {
		// found protocol
		val.Protocal = url[:idx]
		url = url[idx+3:]
	}
	idx = strings.LastIndex(url, ":")
	if idx > 0 {
		digiPart := url[idx+1:]
		// test digits
		port, err := strconv.Atoi(digiPart)
		if err != nil {
			return nil, fmt.Errorf("%v is not a port", digiPart)
		}
		val.PublicPort = port
		url = url[:idx]
	}
	val.PublicURL = url

	url = data.ForwardsTo
	idx = strings.LastIndex(url, ":")
	if idx > 0 {
		digiPart := url[idx+1:]
		// test digits
		port, err := strconv.Atoi(digiPart)
		if err != nil {
			return nil, fmt.Errorf("%v is not a port", digiPart)
		}
		val.LocalPort = port
	}
	return val, nil
}
