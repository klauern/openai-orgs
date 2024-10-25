package openaiorgs

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
)

func TestListUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name      string
		limit     int
		after     string
		mockResp  *resty.Response
		mockErr   error
		expected  *ListResponse[Users]
		expectErr bool
	}{
		{
			name:  "Successful ListUsers request",
			limit: 10,
			after: "test-after-id",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"object":"list","data":[{"id":"user1","name":"User 1","email":"user1@example.com","role":"role1","added_at":"2023-08-01T12:34:56Z"}]}`)),
				},
			},
			expected: &ListResponse[Users]{
				Object: "list",
				Data: []Users{
					{
						ID:      "user1",
						Name:    "User 1",
						Email:   "user1@example.com",
						Role:    "role1",
						AddedAt: CustomTime(time.Date(2023, 8, 1, 12, 34, 56, 0, time.UTC)),
					},
				},
			},
			expectErr: false,
		},
		{
			name:      "ListUsers request with error",
			limit:     10,
			after:     "test-after-id",
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:  "ListUsers request with non-200 status code",
			limit: 10,
			after: "test-after-id",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal server error"}`)),
				},
			},
			expectErr: true,
		},
		{
			name:  "ListUsers request with invalid JSON response",
			limit: 10,
			after: "test-after-id",
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
			mockRequest.EXPECT().Get(UsersListEndpoint).Return(tt.mockResp, tt.mockErr).AnyTimes()

			result, err := client.ListUsers(tt.limit, tt.after)
			if (err != nil) != tt.expectErr {
				t.Errorf("ListUsers() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ListUsers() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRetrieveUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name      string
		userID    string
		mockResp  *resty.Response
		mockErr   error
		expected  *Users
		expectErr bool
	}{
		{
			name:   "Successful RetrieveUser request",
			userID: "test-user-id",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"id":"user1","name":"User 1","email":"user1@example.com","role":"role1","added_at":"2023-08-01T12:34:56Z"}`)),
				},
			},
			expected: &Users{
				ID:      "user1",
				Name:    "User 1",
				Email:   "user1@example.com",
				Role:    "role1",
				AddedAt: CustomTime(time.Date(2023, 8, 1, 12, 34, 56, 0, time.UTC)),
			},
			expectErr: false,
		},
		{
			name:      "RetrieveUser request with error",
			userID:    "test-user-id",
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:   "RetrieveUser request with non-200 status code",
			userID: "test-user-id",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal server error"}`)),
				},
			},
			expectErr: true,
		},
		{
			name:   "RetrieveUser request with invalid JSON response",
			userID: "test-user-id",
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
			mockRequest.EXPECT().Get(UsersListEndpoint+"/"+tt.userID).Return(tt.mockResp, tt.mockErr).AnyTimes()

			result, err := client.RetrieveUser(tt.userID)
			if (err != nil) != tt.expectErr {
				t.Errorf("RetrieveUser() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RetrieveUser() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name      string
		userID    string
		mockResp  *resty.Response
		mockErr   error
		expectErr bool
	}{
		{
			name:   "Successful DeleteUser request",
			userID: "test-user-id",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
				},
			},
			expectErr: false,
		},
		{
			name:      "DeleteUser request with error",
			userID:    "test-user-id",
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:   "DeleteUser request with non-200 status code",
			userID: "test-user-id",
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
			mockRequest.EXPECT().Delete(UsersListEndpoint+"/"+tt.userID).Return(tt.mockResp, tt.mockErr).AnyTimes()

			err := client.DeleteUser(tt.userID)
			if (err != nil) != tt.expectErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestModifyUserRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name      string
		userID    string
		role      RoleType
		mockResp  *resty.Response
		mockErr   error
		expectErr bool
	}{
		{
			name:   "Successful ModifyUserRole request",
			userID: "test-user-id",
			role:   RoleTypeOwner,
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"id":"user1","name":"User 1","email":"user1@example.com","role":"owner","added_at":"2023-08-01T12:34:56Z"}`)),
				},
			},
			expectErr: false,
		},
		{
			name:      "ModifyUserRole request with error",
			userID:    "test-user-id",
			role:      RoleTypeOwner,
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:   "ModifyUserRole request with non-200 status code",
			userID: "test-user-id",
			role:   RoleTypeOwner,
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal server error"}`)),
				},
			},
			expectErr: true,
		},
		{
			name:   "ModifyUserRole request with invalid JSON response",
			userID: "test-user-id",
			role:   RoleTypeOwner,
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
			mockRequest.EXPECT().Post(UsersListEndpoint+"/"+tt.userID).Return(tt.mockResp, tt.mockErr).AnyTimes()

			err := client.ModifyUserRole(tt.userID, tt.role)
			if (err != nil) != tt.expectErr {
				t.Errorf("ModifyUserRole() error = %v, wantErr %v", err, tt.expectErr)
			}
		})
	}
}
