package client

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// HTTPHelper defines an interface for helper methods that encapsulates http requests complexities
type HTTPHelper interface {
	Post(url string, data []byte) (http.Response, []byte, error)
	Get(url string) (http.Response, []byte, error)
	Delete(url string) (http.Response, []byte, error)
}

// BindmanHTTPHelper defines a struct that handles with HTTP requests for a bindman webhook client
type BindmanHTTPHelper struct{}

// Post wraps the call to http.NewRequest apis and properly submits a new HTTP POST request
func (bhh *BindmanHTTPHelper) Post(url string, data []byte) (http.Response, []byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		logrus.Errorf("HTTP request creation failed. err=%v", err)
		return http.Response{}, []byte{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	logrus.Debugf("POST request=%v", req)
	response, err1 := client.Do(req)
	if err1 != nil {
		logrus.Errorf("HTTP request invocation failed. err=%v", err1)
		return http.Response{}, []byte{}, err1
	}

	logrus.Debugf("Response: %v", response)
	datar, _ := ioutil.ReadAll(response.Body)
	logrus.Debugf("Response body: %v", datar)
	return *response, datar, nil
}

// Get wraps the call to http.NewRequest apis and properly submits a new HTTP GET request
func (bhh *BindmanHTTPHelper) Get(url string) (http.Response, []byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Errorf("HTTP request creation failed. err=%s", err)
		return http.Response{}, []byte{}, err
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	logrus.Debugf("GET request=%v", req)
	response, err1 := client.Do(req)
	if err1 != nil {
		logrus.Errorf("HTTP request invocation failed. err=%s", err1)
		return http.Response{}, []byte{}, err1
	}

	logrus.Debugf("Response: %v", response)
	datar, _ := ioutil.ReadAll(response.Body)
	logrus.Debugf("Response body: %v", datar)
	return *response, datar, nil
}

// Delete wraps the call to http.NewRequest apis and properly submits a new HTTP DELETE request
func (bhh *BindmanHTTPHelper) Delete(url string) (http.Response, []byte, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		logrus.Errorf("HTTP request creation failed. err=%v", err)
		return http.Response{}, []byte{}, err
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	logrus.Debugf("DELETE request=%v", req)
	response, err1 := client.Do(req)
	if err1 != nil {
		logrus.Errorf("HTTP request invocation failed. err=%v", err1)
		return http.Response{}, []byte{}, err1
	}

	logrus.Debugf("Response: %v", response)
	datar, _ := ioutil.ReadAll(response.Body)
	logrus.Debugf("Response body: %v", datar)
	return *response, datar, nil
}
