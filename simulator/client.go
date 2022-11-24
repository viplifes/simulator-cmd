package simulator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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

	fullUrl := baseURL + apiUrl + "?" + params.urlEncode()
	log.Printf("[REQ GET %s]", fullUrl)
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	return c.request(req)
}

func (c Client) Post(apiUrl string, data interface{}, params Request) (Response, error) {

	fullUrl := baseURL + apiUrl + "?" + params.urlEncode()

	reqBody, jsonData, err := toBody(data)
	if err != nil {
		return nil, errors.New("json_error")
	}
	log.Printf("[REQ POST %s] %s", fullUrl, jsonData)
	req, err := http.NewRequest("POST", fullUrl, reqBody)
	if err != nil {
		return nil, err
	}
	return c.request(req)
}

func (c Client) Put(apiUrl string, data interface{}, params Request) (Response, error) {

	fullUrl := baseURL + apiUrl + "?" + params.urlEncode()

	reqBody, jsonData, err := toBody(data)
	if err != nil {
		return nil, errors.New("json_error")
	}
	log.Printf("[REQ PUT %s] %s", fullUrl, jsonData)
	req, err := http.NewRequest("PUT", fullUrl, reqBody)
	if err != nil {
		return nil, err
	}
	return c.request(req)
}

func (c Client) Delete(apiUrl string, data interface{}, params Request) (Response, error) {

	fullUrl := baseURL + apiUrl + "?" + params.urlEncode()
	reqBody, jsonData, err := toBody(data)
	if err != nil {
		return nil, errors.New("json_error")
	}
	log.Printf("[REQ DELETE %s] %s", fullUrl, jsonData)
	req, err := http.NewRequest("DELETE", fullUrl, reqBody)
	if err != nil {
		return nil, err
	}
	return c.request(req)
}

func (c Client) request(req *http.Request) (Response, error) {

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("[RES %d] %s", resp.StatusCode, respBody)
	if resp.StatusCode != 200 {

		var errJson Response
		if err := json.Unmarshal(respBody, &errJson); err != nil {
			return nil, errors.New(fmt.Sprintf("bad response status code, status code is %d", resp.StatusCode))
		}
		errMessage := errJson["message"]
		return nil, errors.New(fmt.Sprintf("[error %d] %s", resp.StatusCode, errMessage))
	}

	var res Response
	if err := json.Unmarshal(respBody, &res); err != nil {
		return nil, err
	}

	return res, nil
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

func toBody(data interface{}) (io.Reader, []byte, error) {
	var val io.Reader
	var valJson []byte
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, nil, errors.New("json_error")
		}
		val = bytes.NewBuffer(jsonData)
		valJson = jsonData
	}
	return val, valJson, nil
}
