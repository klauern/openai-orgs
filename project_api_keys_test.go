package openaiorgs

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
)

func TestListProjectApiKeys(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name      string
		projectID string
		limit     int
		after     string
		mockResp  *resty.Response
		mockErr   error
		expected  *ListResponse[ProjectApiKey]
		expectErr bool
	}{
		{
			name:      "Successful ListProjectApiKeys request",
			projectID: "test-project-id",
			limit:     10,
			after:     "test-after-id",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"object":"list","data":[{"id":"key1","name":"Key 1","redacted_value":"****","created_at":"2023-08-01T12:34:56Z","owner":{"id":"owner1","name":"Owner 1","type":"user"}}]}`)),
				},
			},
			expected: &ListResponse[ProjectApiKey]{
				Object: "list",
				Data: []ProjectApiKey{
					{
						ID:            "key1",
						Name:          "Key 1",
						RedactedValue: "****",
						CreatedAt:     CustomTime(time.Date(2023, 8, 1, 12, 34, 56, 0, time.UTC)),
						Owner: Owner{
							ID:   "owner1",
							Name: "Owner 1",
							Type: OwnerTypeUser,
						},
					},
				},
			},
			expectErr: false,
		},
		{
			name:      "ListProjectApiKeys request with error",
			projectID: "test-project-id",
			limit:     10,
			after:     "test-after-id",
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:      "ListProjectApiKeys request with non-200 status code",
			projectID: "test-project-id",
			limit:     10,
			after:     "test-after-id",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal server error"}`)),
				},
			},
			expectErr: true,
		},
		{
			name:      "ListProjectApiKeys request with invalid JSON response",
			projectID: "test-project-id",
			limit:     10,
			after:     "test-after-id",
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
			mockRequest.EXPECT().SetQueryParams(gomock.Any()).Return(mockRequest).AnyTimes()
			mockRequest.EXPECT().ExpectContentType("application/json").Return(mockRequest).AnyTimes()
			mockRequest.EXPECT().Get(fmt.Sprintf(ProjectApiKeysListEndpoint, tt.projectID)).Return(tt.mockResp, tt.mockErr).AnyTimes()

			result, err := client.ListProjectApiKeys(tt.projectID, tt.limit, tt.after)
			if (err != nil) != tt.expectErr {
				t.Errorf("ListProjectApiKeys() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ListProjectApiKeys() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRetrieveProjectApiKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name      string
		projectID string
		apiKeyID  string
		mockResp  *resty.Response
		mockErr   error
		expected  *ProjectApiKey
		expectErr bool
	}{
		{
			name:      "Successful RetrieveProjectApiKey request",
			projectID: "test-project-id",
			apiKeyID:  "test-api-key-id",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"id":"key1","name":"Key 1","redacted_value":"****","created_at":"2023-08-01T12:34:56Z","owner":{"id":"owner1","name":"Owner 1","type":"user"}}`)),
				},
			},
			expected: &ProjectApiKey{
				ID:            "key1",
				Name:          "Key 1",
				RedactedValue: "****",
				CreatedAt:     CustomTime(time.Date(2023, 8, 1, 12, 34, 56, 0, time.UTC)),
				Owner: Owner{
					ID:   "owner1",
					Name: "Owner 1",
					Type: OwnerTypeUser,
				},
			},
			expectErr: false,
		},
		{
			name:      "RetrieveProjectApiKey request with error",
			projectID: "test-project-id",
			apiKeyID:  "test-api-key-id",
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:      "RetrieveProjectApiKey request with non-200 status code",
			projectID: "test-project-id",
			apiKeyID:  "test-api-key-id",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal server error"}`)),
				},
			},
			expectErr: true,
		},
		{
			name:      "RetrieveProjectApiKey request with invalid JSON response",
			projectID: "test-project-id",
			apiKeyID:  "test-api-key-id",
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
			mockRequest.EXPECT().Get(fmt.Sprintf(ProjectApiKeysListEndpoint+"/%s", tt.projectID, tt.apiKeyID)).Return(tt.mockResp, tt.mockErr).AnyTimes()

			result, err := client.RetrieveProjectApiKey(tt.projectID, tt.apiKeyID)
			if (err != nil) != tt.expectErr {
				t.Errorf("RetrieveProjectApiKey() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RetrieveProjectApiKey() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDeleteProjectApiKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name      string
		projectID string
		apiKeyID  string
		mockResp  *resty.Response
		mockErr   error
		expectErr bool
	}{
		{
			name:      "Successful DeleteProjectApiKey request",
			projectID: "test-project-id",
			apiKeyID:  "test-api-key-id",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
				},
			},
			expectErr: false,
		},
		{
			name:      "DeleteProjectApiKey request with error",
			projectID: "test-project-id",
			apiKeyID:  "test-api-key-id",
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:      "DeleteProjectApiKey request with non-200 status code",
			projectID: "test-project-id",
			apiKeyID:  "test-api-key-id",
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
			mockRequest.EXPECT().Delete(fmt.Sprintf(ProjectApiKeysListEndpoint+"/%s", tt.projectID, tt.apiKeyID)).Return(tt.mockResp, tt.mockErr).AnyTimes()

			err := client.DeleteProjectApiKey(tt.projectID, tt.apiKeyID)
			if (err != nil) != tt.expectErr {
				t.Errorf("DeleteProjectApiKey() error = %v, wantErr %v", err, tt.expectErr)
			}
		})
	}
}
