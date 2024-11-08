package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestIsFileThere(t *testing.T) {
	testFile := "./test.txt"
	_, err := os.Create(testFile)
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(testFile)

	if !isFileThere(testFile) {
		t.Errorf("File is supposed to be there.")
	}
	if isFileThere("?") {
		t.Errorf("File is not supposed to be there")
	}
}

func TestErrorHandler(t *testing.T) {
	testCases := []struct {
		method         string
		url            string
		formData       string
		expectedText   string
		expectedStatus int
	}{
		{
			method:         "GET",
			url:            "/asc11-a47",
			formData:       "",
			expectedText:   "404 page not found",
			expectedStatus: 404,
		}, {
			method:         "POST",
			url:            "/",
			formData:       "text=hey&banner=sun",
			expectedText:   "500 internal server error",
			expectedStatus: 500,
		}, {
			method:         "POST",
			url:            "/",
			formData:       "text=&banner=shadow",
			expectedText:   "400 bad request",
			expectedStatus: 400,
		},
		{
			method:         "HEAD",
			url:            "/",
			expectedText:   "405 method not allowed",
			expectedStatus: 405,
		},
		{
			method:         "PUT",
			url:            "/ascii-art",
			expectedText:   "405 method not allowed",
			expectedStatus: 405,
		},
		{
			method:         "DELETE",
			url:            "/",
			expectedText:   "405 method not allowed",
			expectedStatus: 405,
		},
		{
			method:         "CONNECT",
			url:            "/",
			expectedText:   "405 method not allowed",
			expectedStatus: 405,
		},
		{
			method:         "OPTIONS",
			url:            "/",
			expectedText:   "405 method not allowed",
			expectedStatus: 405,
		},
		{
			method:         "TRACE",
			url:            "/",
			expectedText:   "405 method not allowed",
			expectedStatus: 405,
		},
		{
			method:         "PATCH",
			url:            "/",
			expectedText:   "405 method not allowed",
			expectedStatus: 405,
		},
	}

	for _, tc := range testCases {
		req, err := http.NewRequest(tc.method, tc.url, strings.NewReader(tc.formData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		respRec := httptest.NewRecorder()
		handler := http.HandlerFunc(homeHandler)
		handler.ServeHTTP(respRec, req)

		if status := respRec.Code; status != tc.expectedStatus {
			t.Errorf("Expected HTTP Status Code %v got %v.", tc.expectedStatus, status)
		}
		if !strings.Contains(respRec.Body.String(), tc.expectedText) {
			t.Errorf("Expected content not found. %s", tc.expectedText)
		}
	}
}

func TestGetFormInput(t *testing.T) {
	testCases := []struct {
		formData       string
		expectedText   string
		expectedBanner string
		expectedError  string
	}{
		{
			formData:       "text=hello&banner=shadow",
			expectedText:   "hello",
			expectedBanner: "shadow",
			expectedError:  "",
		},
		{
			formData:       "banner=shadow",
			expectedText:   "",
			expectedBanner: "",
			expectedError:  "400",
		},
		{
			formData:       "text=hello",
			expectedText:   "",
			expectedBanner: "",
			expectedError:  "400",
		},
	}
	for _, tc := range testCases {
		req, err := http.NewRequest("POST", "/", strings.NewReader(tc.formData))
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		text, banner, errForm := getFormInputs(req)

		if text != tc.expectedText ||
			banner != tc.expectedBanner ||
			errForm != tc.expectedError {
			t.Errorf("Failed reading form: %s", tc.formData)
		}

	}
}

func TestHomeHandlerGet(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	respRec := httptest.NewRecorder()
	handler := http.HandlerFunc(homeHandler)
	handler.ServeHTTP(respRec, req)

	if status := respRec.Code; status != http.StatusOK {
		t.Errorf("GET: Expected HTTP Status Code 200.")
	}
	if !strings.Contains(respRec.Body.String(), "Enter text\nhere") {
		t.Errorf("GET: Expected content not found.")
	}
}

func TestHomeHandlerPost(t *testing.T) {
	testCases := []struct {
		formData       string
		expectedText   string
		expectedStatus int
	}{
		{
			formData: "text=hello&banner=standard",
			expectedText: ` _              _   _          
| |            | | | |         
| |__     ___  | | | |   ___   
|  _ \   / _ \ | | | |  / _ \  
| | | | |  __/ | | | | | (_) | 
|_| |_|  \___| |_| |_|  \___/  
                               `,
			expectedStatus: 200,
		},
		{
			formData: "text=\n\nhello There\n\nAuditor!&banner=shadow",
			expectedText: `                                                                                      
_|                _| _|                _|_|_|_|_| _|                                  
_|_|_|     _|_|   _| _|   _|_|             _|     _|_|_|     _|_|   _|  _|_|   _|_|   
_|    _| _|_|_|_| _| _| _|    _|           _|     _|    _| _|_|_|_| _|_|     _|_|_|_| 
_|    _| _|       _| _| _|    _|           _|     _|    _| _|       _|       _|       
_|    _|   _|_|_| _| _|   _|_|             _|     _|    _|   _|_|_| _|         _|_|_| 
                                                                                      
                                                                                      

                                                            
  _|_|                  _| _|   _|                       _| 
_|    _| _|    _|   _|_|_|    _|_|_|_|   _|_|   _|  _|_| _| 
_|_|_|_| _|    _| _|    _| _|   _|     _|    _| _|_|     _| 
_|    _| _|    _| _|    _| _|   _|     _|    _| _|          
_|    _|   _|_|_|   _|_|_| _|     _|_|   _|_|   _|       _| 
                                                            
                                                            `,
			expectedStatus: 200,
		},
		{
			formData:       "text=äöå&banner=standard",
			expectedText:   `character not a printable ASCII char: ä`,
			expectedStatus: 200,
		},
		{
			formData:       "text=\n\n\n&banner=standard",
			expectedText:   `no character to convert`,
			expectedStatus: 200,
		},
	}
	for _, tc := range testCases {
		req, err := http.NewRequest("POST", "/ascii-art", strings.NewReader(tc.formData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		respRec := httptest.NewRecorder()
		handler := http.HandlerFunc(homeHandler)
		handler.ServeHTTP(respRec, req)

		if status := respRec.Code; status != tc.expectedStatus {
			t.Errorf("Expected HTTP Status Code %v got %v.", tc.expectedStatus, status)
		}
		if !strings.Contains(respRec.Body.String(), tc.expectedText) {
			t.Errorf("Expected content not found. %s", tc.expectedText)
		}
	}
}
