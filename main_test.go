package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"task-axxon/transport"
	"task-axxon/view"
	"testing"
)

// TestSimple checks the following:
// - are the status codes the same
// - content length between the request via proxy and the original one
// - whether the UUID is in storage
func TestSimple(t *testing.T) {
	testRequests := []view.Request{
		{
			URL:    "https://httpbin.org/ip",
			Method: http.MethodGet,
		},
		{
			URL:    "https://httpbin.org/status/404",
			Method: http.MethodGet,
		},
		{
			URL:    "https://httpbin.org/post",
			Method: http.MethodPost,
			Body: struct {
				Name string
			}{
				"John",
			},
		},
		{
			URL:    "https://httpbin.org/response-headers",
			Method: http.MethodGet,
			Headers: map[string][]string{
				"name": {"John"},
			},
		},
	}

	for i, req := range testRequests {
		t.Run(fmt.Sprintf("test #%d", i), func(t *testing.T) {
			if req.Headers != nil {
				url, err := url.Parse(req.URL)
				if err != nil {
					t.Errorf("could not parse url for headers: %s", err)
					return
				}
				q := url.Query()
				for k, values := range req.Headers {
					for _, v := range values {
						q.Add(k, v)
					}
				}
				url.RawQuery = q.Encode()
				req.URL = url.String()
			}

			bodyProxy := new(bytes.Buffer)
			err := json.NewEncoder(bodyProxy).Encode(req)
			if err != nil {
				t.Errorf("could not encode body for proxy: %s", err)
				return
			}

			requestProxy := httptest.NewRequest(http.MethodPost, "/", bodyProxy)
			responseProxy := httptest.NewRecorder()
			transport.ProxyHTTP(responseProxy, requestProxy)

			log.Println(responseProxy.Body.String())

			bodyOriginal := new(bytes.Buffer)
			err = json.NewEncoder(bodyOriginal).Encode(req.Body)
			if err != nil {
				t.Errorf("could not encode body for request to original: %s", err)
				return
			}
			requestOriginal, err := http.NewRequest(req.Method, req.URL, bodyOriginal)
			if err != nil {
				t.Errorf("could not build a request to original: %s", err)
				return
			}
			responseOriginal, err := http.DefaultClient.Do(requestOriginal)
			defer responseOriginal.Body.Close()
			if err != nil {
				t.Errorf("could not request to original: %s", err)
				return
			}
			var responseProxyView view.Response
			err = json.NewDecoder(responseProxy.Body).Decode(&responseProxyView)
			if err != nil {
				t.Errorf("could not decode response proxy: %s", err)
				return
			}
			if responseOriginal.StatusCode != responseProxyView.Status {
				t.Errorf("Want status '%d', got '%d'", responseOriginal.StatusCode, responseProxyView.Status)
				return
			}
			if responseOriginal.ContentLength != responseProxyView.Length {
				t.Errorf("Want content length '%d', got '%d'", responseOriginal.ContentLength, responseProxyView.Length)
				return
			}
			for k, values := range req.Headers {
				for _, v := range values {
					valGot := responseProxyView.Headers.Get(k)
					if valGot != v {
						fmt.Println(req.URL)
						t.Errorf("headers are not synced: %s != %s", valGot, v)
						return
					}
				}
			}

			requestUUID := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?uuid=%s", responseProxyView.UUID), nil)
			responseUUID := httptest.NewRecorder()
			transport.GetByUUID(responseUUID, requestUUID)
			if responseUUID.Result().StatusCode != http.StatusOK {
				t.Errorf("could not find uuid in store: %s", responseProxyView.UUID)
				return
			}
		})
	}
}
