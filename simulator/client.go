package simulator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const baseURL = "https://api.control.events/v/1.0/"

type Client struct {
	httpClient *http.Client
	token      string
}

type Request map[string]interface{}
type Response map[string]interface{}

func New(token string) *Client {
	return &Client{
		httpClient: &http.Client{},
		token:      token,
	}
}

func (c Client) Get(apiUrl string, params Request) (Response, error) {

	url := baseURL + apiUrl + "?" + params.urlEncode()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		//Handle Error
	}
	req.Header = http.Header{
		"Authorization": {"Bearer " + c.token},
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		//Handle Error
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("bad response status code, status code is %d", resp.StatusCode))
	}

	var res Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c Client) Post(apiUrl string, data Request, params Request) (Response, error) {

	url := baseURL + apiUrl + "?" + params.urlEncode()
	reqBody := bytes.NewBuffer(data.jsonEncode())
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		//Handle Error
	}
	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + c.token},
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		//Handle Error
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("bad response status code, status code is %d", resp.StatusCode))
	}

	var res Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c Client) Put(apiUrl string, data Request, params Request) (Response, error) {

	url := baseURL + apiUrl + "?" + params.urlEncode()

	reqBody := bytes.NewBuffer(data.jsonEncode())
	req, err := http.NewRequest("PUT", url, reqBody)
	if err != nil {
		//Handle Error
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		//Handle Error
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, errors.New(fmt.Sprintf("bad response status code, status code is %d respBody: %s", resp.StatusCode, respBody))
	}

	var res Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (r Request) jsonEncode() []byte {
	obj, err := json.Marshal(r)
	if err != nil {
		return []byte{}
	}
	return obj
}

func (r Request) urlEncode() string {

	if r == nil {
		return ""
	}
	params := url.Values{}

	for k, v := range r {
		params.Add(k, v.(string))
	}
	return params.Encode()
}
