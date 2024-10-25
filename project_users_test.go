package openaiorgs

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/go-resty/resty/v2"
)

func TestListProjectUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name        string
		projectID   string
		limit       int
		after       string
		mockResp    *resty.Response
		mockErr     error
		expected    *ListResponse[ProjectUser]
		expectErr   bool
	}{
		{
			name:      "Successful ListProjectUsers request",
			projectID: "test-project-id",
			limit:     10,
			after:     "test-after-id",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"object":"list","data":[{"id":"user1","name":"User 1","email":"user1@example.com","role":"role1","added_at":"2023-08-01T12:34:56Z"}]}`)),
				},
			},
			expected: &ListResponse[ProjectUser]{
				Object: "list",
				Data: []ProjectUser{
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
			name:      "ListProjectUsers request with error",
			projectID: "test-project-id",
			limit:     10,
			after:     "test-after-id",
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:      "ListProjectUsers request with non-200 status code",
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
			name:      "ListProjectUsers request with invalid JSON response",
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
			mockRequest.EXPECT().Get(fmt.Sprintf(ProjectUsersListEndpoint, tt.projectID)).Return(tt.mockResp, tt.mockErr).AnyTimes()

			result, err := client.ListProjectUsers(tt.projectID, tt.limit, tt.after)
			if (err != nil) != tt.expectErr {
				t.Errorf("ListProjectUsers() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ListProjectUsers() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCreateProjectUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name        string
		projectID   string
		userID      string
		role        RoleType
		mockResp    *resty.Response
		mockErr     error
		expected    *ProjectUser
		expectErr   bool
	}{
		{
			name:      "Successful CreateProjectUser request",
			projectID: "test-project-id",
			userID:    "test-user-id",
			role:      RoleTypeMember,
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"id":"user1","name":"User 1","email":"user1@example.com","role":"member","added_at":"2023-08-01T12:34:56Z"}`)),
				},
			},
			expected: &ProjectUser{
				ID:      "user1",
				Name:    "User 1",
				Email:   "user1@example.com",
				Role:    "member",
				AddedAt: CustomTime(time.Date(2023, 8, 1, 12, 34, 56, 0, time.UTC)),
			},
			expectErr: false,
		},
		{
			name:      "CreateProjectUser request with error",
			projectID: "test-project-id",
			userID:    "test-user-id",
			role:      RoleTypeMember,
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:      "CreateProjectUser request with non-200 status code",
			projectID: "test-project-id",
			userID:    "test-user-id",
			role:      RoleTypeMember,
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal server error"}`)),
				},
			},
			expectErr: true,
		},
		{
			name:      "CreateProjectUser request with invalid JSON response",
			projectID: "test-project-id",
			userID:    "test-user-id",
			role:      RoleTypeMember,
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
			mockRequest.EXPECT().Post(fmt.Sprintf(ProjectUsersListEndpoint, tt.projectID)).Return(tt.mockResp, tt.mockErr).AnyTimes()

			result, err := client.CreateProjectUser(tt.projectID, tt.userID, tt.role)
			if (err != nil) != tt.expectErr {
				t.Errorf("CreateProjectUser() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("CreateProjectUser() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRetrieveProjectUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name        string
		projectID   string
		userID      string
		mockResp    *resty.Response
		mockErr     error
		expected    *ProjectUser
		expectErr   bool
	}{
		{
			name:      "Successful RetrieveProjectUser request",
			projectID: "test-project-id",
			userID:    "test-user-id",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"id":"user1","name":"User 1","email":"user1@example.com","role":"member","added_at":"2023-08-01T12:34:56Z"}`)),
				},
			},
			expected: &ProjectUser{
				ID:      "user1",
				Name:    "User 1",
				Email:   "user1@example.com",
				Role:    "member",
				AddedAt: CustomTime(time.Date(2023, 8, 1, 12, 34, 56, 0, time.UTC)),
			},
			expectErr: false,
		},
		{
			name:      "RetrieveProjectUser request with error",
			projectID: "test-project-id",
			userID:    "test-user-id",
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:      "RetrieveProjectUser request with non-200 status code",
			projectID: "test-project-id",
			userID:    "test-user-id",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal server error"}`)),
				},
			},
			expectErr: true,
		},
		{
			name:      "RetrieveProjectUser request with invalid JSON response",
			projectID: "test-project-id",
			userID:    "test-user-id",
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
			mockRequest.EXPECT().Get(fmt.Sprintf(ProjectUsersListEndpoint+"/%s", tt.projectID, tt.userID)).Return(tt.mockResp, tt.mockErr).AnyTimes()

			result, err := client.RetrieveProjectUser(tt.projectID, tt.userID)
			if (err != nil) != tt.expectErr {
				t.Errorf("RetrieveProjectUser() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RetrieveProjectUser() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestModifyProjectUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name        string
		projectID   string
		userID      string
		role        RoleType
		mockResp    *resty.Response
		mockErr     error
		expected    *ProjectUser
		expectErr   bool
	}{
		{
			name:      "Successful ModifyProjectUser request",
			projectID: "test-project-id",
			userID:    "test-user-id",
			role:      RoleTypeOwner,
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"id":"user1","name":"User 1","email":"user1@example.com","role":"owner","added_at":"2023-08-01T12:34:56Z"}`)),
				},
			},
			expected: &ProjectUser{
				ID:      "user1",
				Name:    "User 1",
				Email:   "user1@example.com",
				Role:    "owner",
				AddedAt: CustomTime(time.Date(2023, 8, 1, 12, 34, 56, 0, time.UTC)),
			},
			expectErr: false,
		},
		{
			name:      "ModifyProjectUser request with error",
			projectID: "test-project-id",
			userID:    "test-user-id",
			role:      RoleTypeOwner,
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:      "ModifyProjectUser request with non-200 status code",
			projectID: "test-project-id",
			userID:    "test-user-id",
			role:      RoleTypeOwner,
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal server error"}`)),
				},
			},
			expectErr: true,
		},
		{
			name:      "ModifyProjectUser request with invalid JSON response",
			projectID: "test-project-id",
			userID:    "test-user-id",
			role:      RoleTypeOwner,
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
			mockRequest.EXPECT().Post(fmt.Sprintf(ProjectUsersListEndpoint+"/%s", tt.projectID, tt.userID)).Return(tt.mockResp, tt.mockErr).AnyTimes()

			result, err := client.ModifyProjectUser(tt.projectID, tt.userID, tt.role)
			if (err != nil) != tt.expectErr {
				t.Errorf("ModifyProjectUser() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ModifyProjectUser() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDeleteProjectUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name        string
		projectID   string
		userID      string
		mockResp    *resty.Response
		mockErr     error
		expectErr   bool
	}{
		{
			name:      "Successful DeleteProjectUser request",
			projectID: "test-project-id",
			userID:    "test-user-id",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
				},
			},
			expectErr: false,
		},
		{
			name:      "DeleteProjectUser request with error",
			projectID: "test-project-id",
			userID:    "test-user-id",
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:      "DeleteProjectUser request with non-200 status code",
			projectID: "test-project-id",
			userID:    "test-user-id",
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
			mockRequest.EXPECT().Delete(fmt.Sprintf(ProjectUsersListEndpoint+"/%s", tt.projectID, tt.userID)).Return(tt.mockResp, tt.mockErr).AnyTimes()

			err := client.DeleteProjectUser(tt.projectID, tt.userID)
			if (err != nil) != tt.expectErr {
				t.Errorf("DeleteProjectUser() error = %v, wantErr %v", err, tt.expectErr)
			}
		})
	}
}
