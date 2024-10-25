package openaiorgs

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
)

func TestListProjectServiceAccounts(t *testing.T) {
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
		expected  *ListResponse[ProjectServiceAccount]
		expectErr bool
	}{
		{
			name:      "Successful ListProjectServiceAccounts request",
			projectID: "test-project-id",
			limit:     10,
			after:     "test-after-id",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"object":"list","data":[{"id":"account1","name":"Account 1","role":"role1","created_at":"2023-08-01T12:34:56Z"}]}`)),
				},
			},
			expected: &ListResponse[ProjectServiceAccount]{
				Object: "list",
				Data: []ProjectServiceAccount{
					{
						ID:        "account1",
						Name:      "Account 1",
						Role:      "role1",
						CreatedAt: CustomTime(time.Date(2023, 8, 1, 12, 34, 56, 0, time.UTC)),
					},
				},
			},
			expectErr: false,
		},
		{
			name:      "ListProjectServiceAccounts request with error",
			projectID: "test-project-id",
			limit:     10,
			after:     "test-after-id",
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:      "ListProjectServiceAccounts request with non-200 status code",
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
			name:      "ListProjectServiceAccounts request with invalid JSON response",
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
			mockRequest.EXPECT().Get(fmt.Sprintf(ProjectServiceAccountsListEndpoint, tt.projectID)).Return(tt.mockResp, tt.mockErr).AnyTimes()

			result, err := client.ListProjectServiceAccounts(tt.projectID, tt.limit, tt.after)
			if (err != nil) != tt.expectErr {
				t.Errorf("ListProjectServiceAccounts() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ListProjectServiceAccounts() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCreateProjectServiceAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name        string
		projectID   string
		accountName string
		mockResp    *resty.Response
		mockErr     error
		expected    *ProjectServiceAccount
		expectErr   bool
	}{
		{
			name:        "Successful CreateProjectServiceAccount request",
			projectID:   "test-project-id",
			accountName: "Account 1",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"id":"account1","name":"Account 1","role":"role1","created_at":"2023-08-01T12:34:56Z"}`)),
				},
			},
			expected: &ProjectServiceAccount{
				ID:        "account1",
				Name:      "Account 1",
				Role:      "role1",
				CreatedAt: CustomTime(time.Date(2023, 8, 1, 12, 34, 56, 0, time.UTC)),
			},
			expectErr: false,
		},
		{
			name:        "CreateProjectServiceAccount request with error",
			projectID:   "test-project-id",
			accountName: "Account 1",
			mockErr:     fmt.Errorf("request error"),
			expectErr:   true,
		},
		{
			name:        "CreateProjectServiceAccount request with non-200 status code",
			projectID:   "test-project-id",
			accountName: "Account 1",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal server error"}`)),
				},
			},
			expectErr: true,
		},
		{
			name:        "CreateProjectServiceAccount request with invalid JSON response",
			projectID:   "test-project-id",
			accountName: "Account 1",
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
			mockRequest.EXPECT().SetBody(gomock.Any()).Return(mockRequest).AnyTimes()
			mockRequest.EXPECT().ExpectContentType("application/json").Return(mockRequest).AnyTimes()
			mockRequest.EXPECT().Post(fmt.Sprintf(ProjectServiceAccountsListEndpoint, tt.projectID)).Return(tt.mockResp, tt.mockErr).AnyTimes()

			result, err := client.CreateProjectServiceAccount(tt.projectID, tt.accountName)
			if (err != nil) != tt.expectErr {
				t.Errorf("CreateProjectServiceAccount() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("CreateProjectServiceAccount() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRetrieveProjectServiceAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name      string
		projectID string
		accountID string
		mockResp  *resty.Response
		mockErr   error
		expected  *ProjectServiceAccount
		expectErr bool
	}{
		{
			name:      "Successful RetrieveProjectServiceAccount request",
			projectID: "test-project-id",
			accountID: "test-account-id",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"id":"account1","name":"Account 1","role":"role1","created_at":"2023-08-01T12:34:56Z"}`)),
				},
			},
			expected: &ProjectServiceAccount{
				ID:        "account1",
				Name:      "Account 1",
				Role:      "role1",
				CreatedAt: CustomTime(time.Date(2023, 8, 1, 12, 34, 56, 0, time.UTC)),
			},
			expectErr: false,
		},
		{
			name:      "RetrieveProjectServiceAccount request with error",
			projectID: "test-project-id",
			accountID: "test-account-id",
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:      "RetrieveProjectServiceAccount request with non-200 status code",
			projectID: "test-project-id",
			accountID: "test-account-id",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal server error"}`)),
				},
			},
			expectErr: true,
		},
		{
			name:      "RetrieveProjectServiceAccount request with invalid JSON response",
			projectID: "test-project-id",
			accountID: "test-account-id",
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
			mockRequest.EXPECT().Get(fmt.Sprintf(ProjectServiceAccountsListEndpoint+"/%s", tt.projectID, tt.accountID)).Return(tt.mockResp, tt.mockErr).AnyTimes()

			result, err := client.RetrieveProjectServiceAccount(tt.projectID, tt.accountID)
			if (err != nil) != tt.expectErr {
				t.Errorf("RetrieveProjectServiceAccount() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RetrieveProjectServiceAccount() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDeleteProjectServiceAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name      string
		projectID string
		accountID string
		mockResp  *resty.Response
		mockErr   error
		expectErr bool
	}{
		{
			name:      "Successful DeleteProjectServiceAccount request",
			projectID: "test-project-id",
			accountID: "test-account-id",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
				},
			},
			expectErr: false,
		},
		{
			name:      "DeleteProjectServiceAccount request with error",
			projectID: "test-project-id",
			accountID: "test-account-id",
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:      "DeleteProjectServiceAccount request with non-200 status code",
			projectID: "test-project-id",
			accountID: "test-account-id",
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
			mockRequest.EXPECT().Delete(fmt.Sprintf(ProjectServiceAccountsListEndpoint+"/%s", tt.projectID, tt.accountID)).Return(tt.mockResp, tt.mockErr).AnyTimes()

			err := client.DeleteProjectServiceAccount(tt.projectID, tt.accountID)
			if (err != nil) != tt.expectErr {
				t.Errorf("DeleteProjectServiceAccount() error = %v, wantErr %v", err, tt.expectErr)
			}
		})
	}
}
