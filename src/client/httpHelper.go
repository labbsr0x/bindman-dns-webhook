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
	Put(url string, data []byte) (*http.Response, []byte, error)
	Post(url string, data []byte) (*http.Response, []byte, error)
	Get(url string) (*http.Response, []byte, error)
	Delete(url string) (*http.Response, []byte, error)
}

// BindmanHTTPHelper defines a struct that handles with HTTP requests for a bindman webhook client
type BindmanHTTPHelper struct{}

// Put wraps the call to http.NewRequest apis and properly submits a new HTTP POST request
func (bhh *BindmanHTTPHelper) Put(url string, data []byte) (*http.Response, []byte, error) {
	return bhh.request(url, "PUT", "application/json", data)
}

// Post wraps the call to http.NewRequest apis and properly submits a new HTTP POST request
func (bhh *BindmanHTTPHelper) Post(url string, data []byte) (*http.Response, []byte, error) {
	return bhh.request(url, "POST", "application/json", data)
}

// Get wraps the call to http.NewRequest apis and properly submits a new HTTP GET request
func (bhh *BindmanHTTPHelper) Get(url string) (*http.Response, []byte, error) {
	return bhh.request(url, "GET", "", nil)
}

// Delete wraps the call to http.NewRequest apis and properly submits a new HTTP DELETE request
func (bhh *BindmanHTTPHelper) Delete(url string) (*http.Response, []byte, error) {
	return bhh.request(url, "DELETE", "", nil)
}

// request defines a generic method to execute http requests
func (bhh *BindmanHTTPHelper) request(url, method, contentType string, payload []byte) (httpResponse *http.Response, data []byte, err error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err == nil {
		req.Header.Set("Content-Type", contentType)

		client := &http.Client{
			Timeout: time.Second * 10,
		}
		logrus.Debugf("%v request=%v", method, req)

		httpResponse, err = client.Do(req)
		if err == nil {
			defer func() {
				if closeError := httpResponse.Body.Close(); closeError != nil {
					logrus.Errorf("HTTP  %v response body close invocation failed. err=%v", method, err)
				}
			}()
			logrus.Debugf("Response: %v", httpResponse)
			data, _ = ioutil.ReadAll(httpResponse.Body)
			logrus.Debugf("Response body: %v", data)
			return
		}
		logrus.Errorf("HTTP  %v request invocation failed. err=%v", method, err)
		return
	}
	logrus.Errorf("HTTP %v request creation failed. err=%v", method, err)
	return
}
