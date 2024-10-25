package openaiorgs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		token    string
		expected string
	}{
		{"DefaultBaseURL", "", "test-token", DefaultBaseURL},
		{"CustomBaseURL", "https://custom.api.com", "test-token", "https://custom.api.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.baseURL, tt.token)
			if client.BaseURL != tt.expected {
				t.Errorf("NewClient() BaseURL = %v, want %v", client.BaseURL, tt.expected)
			}
		})
	}
}

func TestGetSingle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name      string
		endpoint  string
		mockResp  *resty.Response
		mockErr   error
		expected  *TestStruct
		expectErr bool
	}{
		{
			name:     "Successful GET request",
			endpoint: "/test-endpoint",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"field1":"value1","field2":2}`)),
				},
			},
			expected:  &TestStruct{Field1: "value1", Field2: 2},
			expectErr: false,
		},
		{
			name:      "GET request with error",
			endpoint:  "/test-endpoint",
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:     "GET request with non-200 status code",
			endpoint: "/test-endpoint",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal server error"}`)),
				},
			},
			expectErr: true,
		},
		{
			name:     "GET request with invalid JSON response",
			endpoint: "/test-endpoint",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`invalid json`)),
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.EXPECT().R().Return(mockRequest).AnyTimes()
			mockRequest.EXPECT().ExpectContentType("application/json").Return(mockRequest).AnyTimes()
			mockRequest.EXPECT().Get(tt.endpoint).Return(tt.mockResp, tt.mockErr).AnyTimes()

			result, err := GetSingle[TestStruct](client.client, tt.endpoint)
			if (err != nil) != tt.expectErr {
				t.Errorf("GetSingle() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetSingle() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name        string
		endpoint    string
		queryParams map[string]string
		mockResp    *resty.Response
		mockErr     error
		expected    *ListResponse[TestStruct]
		expectErr   bool
	}{
		{
			name:     "Successful GET request",
			endpoint: "/test-endpoint",
			queryParams: map[string]string{
				"param1": "value1",
				"param2": "value2",
			},
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"object":"list","data":[{"field1":"value1","field2":2}]}`)),
				},
			},
			expected: &ListResponse[TestStruct]{
				Object: "list",
				Data:   []TestStruct{{Field1: "value1", Field2: 2}},
			},
			expectErr: false,
		},
		{
			name:      "GET request with error",
			endpoint:  "/test-endpoint",
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:     "GET request with non-200 status code",
			endpoint: "/test-endpoint",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal server error"}`)),
				},
			},
			expectErr: true,
		},
		{
			name:     "GET request with invalid JSON response",
			endpoint: "/test-endpoint",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`invalid json`)),
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.EXPECT().R().Return(mockRequest).AnyTimes()
			mockRequest.EXPECT().SetQueryParams(tt.queryParams).Return(mockRequest).AnyTimes()
			mockRequest.EXPECT().ExpectContentType("application/json").Return(mockRequest).AnyTimes()
			mockRequest.EXPECT().Get(tt.endpoint).Return(tt.mockResp, tt.mockErr).AnyTimes()

			result, err := Get[TestStruct](client.client, tt.endpoint, tt.queryParams)
			if (err != nil) != tt.expectErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Get() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name      string
		endpoint  string
		body      interface{}
		mockResp  *resty.Response
		mockErr   error
		expected  *TestStruct
		expectErr bool
	}{
		{
			name:     "Successful POST request",
			endpoint: "/test-endpoint",
			body:     map[string]string{"field1": "value1", "field2": "2"},
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"field1":"value1","field2":2}`)),
				},
			},
			expected:  &TestStruct{Field1: "value1", Field2: 2},
			expectErr: false,
		},
		{
			name:      "POST request with error",
			endpoint:  "/test-endpoint",
			body:      map[string]string{"field1": "value1", "field2": "2"},
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:     "POST request with non-200 status code",
			endpoint: "/test-endpoint",
			body:     map[string]string{"field1": "value1", "field2": "2"},
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal server error"}`)),
				},
			},
			expectErr: true,
		},
		{
			name:     "POST request with invalid JSON response",
			endpoint: "/test-endpoint",
			body:     map[string]string{"field1": "value1", "field2": "2"},
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`invalid json`)),
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.EXPECT().R().Return(mockRequest).AnyTimes()
			mockRequest.EXPECT().SetBody(tt.body).Return(mockRequest).AnyTimes()
			mockRequest.EXPECT().ExpectContentType("application/json").Return(mockRequest).AnyTimes()
			mockRequest.EXPECT().Post(tt.endpoint).Return(tt.mockResp, tt.mockErr).AnyTimes()

			result, err := Post[TestStruct](client.client, tt.endpoint, tt.body)
			if (err != nil) != tt.expectErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Post() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name      string
		endpoint  string
		mockResp  *resty.Response
		mockErr   error
		expectErr bool
	}{
		{
			name:     "Successful DELETE request",
			endpoint: "/test-endpoint",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
				},
			},
			expectErr: false,
		},
		{
			name:      "DELETE request with error",
			endpoint:  "/test-endpoint",
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:     "DELETE request with non-200 status code",
			endpoint: "/test-endpoint",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal server error"}`)),
				},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.EXPECT().R().Return(mockRequest).AnyTimes()
			mockRequest.EXPECT().Delete(tt.endpoint).Return(tt.mockResp, tt.mockErr).AnyTimes()

			err := Delete[TestStruct](client.client, tt.endpoint)
			if (err != nil) != tt.expectErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.expectErr)
			}
		})
	}
}
