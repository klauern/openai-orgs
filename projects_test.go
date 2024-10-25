package openaiorgs

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/go-resty/resty/v2"
)

func TestListProjects(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name           string
		limit          int
		after          string
		includeArchived bool
		mockResp       *resty.Response
		mockErr        error
		expected       *ListResponse[Project]
		expectErr      bool
	}{
		{
			name:           "Successful ListProjects request",
			limit:          10,
			after:          "test-after-id",
			includeArchived: true,
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"object":"list","data":[{"id":"project1","name":"Project 1","created_at":"2023-08-01T12:34:56Z","status":"active"}]}`)),
				},
			},
			expected: &ListResponse[Project]{
				Object: "list",
				Data: []Project{
					{
						ID:        "project1",
						Name:      "Project 1",
						CreatedAt: CustomTime(time.Date(2023, 8, 1, 12, 34, 56, 0, time.UTC)),
						Status:    "active",
					},
				},
			},
			expectErr: false,
		},
		{
			name:           "ListProjects request with error",
			limit:          10,
			after:          "test-after-id",
			includeArchived: true,
			mockErr:        fmt.Errorf("request error"),
			expectErr:      true,
		},
		{
			name:           "ListProjects request with non-200 status code",
			limit:          10,
			after:          "test-after-id",
			includeArchived: true,
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal server error"}`)),
				},
			},
			expectErr: true,
		},
		{
			name:           "ListProjects request with invalid JSON response",
			limit:          10,
			after:          "test-after-id",
			includeArchived: true,
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
			mockRequest.EXPECT().Get(ProjectsListEndpoint).Return(tt.mockResp, tt.mockErr).AnyTimes()

			result, err := client.ListProjects(tt.limit, tt.after, tt.includeArchived)
			if (err != nil) != tt.expectErr {
				t.Errorf("ListProjects() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ListProjects() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCreateProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name        string
		projectName string
		mockResp    *resty.Response
		mockErr     error
		expected    *Project
		expectErr   bool
	}{
		{
			name:        "Successful CreateProject request",
			projectName: "Project 1",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"id":"project1","name":"Project 1","created_at":"2023-08-01T12:34:56Z","status":"active"}`)),
				},
			},
			expected: &Project{
				ID:        "project1",
				Name:      "Project 1",
				CreatedAt: CustomTime(time.Date(2023, 8, 1, 12, 34, 56, 0, time.UTC)),
				Status:    "active",
			},
			expectErr: false,
		},
		{
			name:        "CreateProject request with error",
			projectName: "Project 1",
			mockErr:     fmt.Errorf("request error"),
			expectErr:   true,
		},
		{
			name:        "CreateProject request with non-200 status code",
			projectName: "Project 1",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal server error"}`)),
				},
			},
			expectErr: true,
		},
		{
			name:        "CreateProject request with invalid JSON response",
			projectName: "Project 1",
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
			mockRequest.EXPECT().Post(ProjectsListEndpoint).Return(tt.mockResp, tt.mockErr).AnyTimes()

			result, err := client.CreateProject(tt.projectName)
			if (err != nil) != tt.expectErr {
				t.Errorf("CreateProject() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("CreateProject() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRetrieveProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name      string
		projectID string
		mockResp  *resty.Response
		mockErr   error
		expected  *Project
		expectErr bool
	}{
		{
			name:      "Successful RetrieveProject request",
			projectID: "test-project-id",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"id":"project1","name":"Project 1","created_at":"2023-08-01T12:34:56Z","status":"active"}`)),
				},
			},
			expected: &Project{
				ID:        "project1",
				Name:      "Project 1",
				CreatedAt: CustomTime(time.Date(2023, 8, 1, 12, 34, 56, 0, time.UTC)),
				Status:    "active",
			},
			expectErr: false,
		},
		{
			name:      "RetrieveProject request with error",
			projectID: "test-project-id",
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:      "RetrieveProject request with non-200 status code",
			projectID: "test-project-id",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal server error"}`)),
				},
			},
			expectErr: true,
		},
		{
			name:      "RetrieveProject request with invalid JSON response",
			projectID: "test-project-id",
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
			mockRequest.EXPECT().Get(fmt.Sprintf("%s/%s", ProjectsListEndpoint, tt.projectID)).Return(tt.mockResp, tt.mockErr).AnyTimes()

			result, err := client.RetrieveProject(tt.projectID)
			if (err != nil) != tt.expectErr {
				t.Errorf("RetrieveProject() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RetrieveProject() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestModifyProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name      string
		projectID string
		newName   string
		mockResp  *resty.Response
		mockErr   error
		expected  *Project
		expectErr bool
	}{
		{
			name:      "Successful ModifyProject request",
			projectID: "test-project-id",
			newName:   "New Project Name",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"id":"project1","name":"New Project Name","created_at":"2023-08-01T12:34:56Z","status":"active"}`)),
				},
			},
			expected: &Project{
				ID:        "project1",
				Name:      "New Project Name",
				CreatedAt: CustomTime(time.Date(2023, 8, 1, 12, 34, 56, 0, time.UTC)),
				Status:    "active",
			},
			expectErr: false,
		},
		{
			name:      "ModifyProject request with error",
			projectID: "test-project-id",
			newName:   "New Project Name",
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:      "ModifyProject request with non-200 status code",
			projectID: "test-project-id",
			newName:   "New Project Name",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal server error"}`)),
				},
			},
			expectErr: true,
		},
		{
			name:      "ModifyProject request with invalid JSON response",
			projectID: "test-project-id",
			newName:   "New Project Name",
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
			mockRequest.EXPECT().Post(fmt.Sprintf("%s/%s", ProjectsListEndpoint, tt.projectID)).Return(tt.mockResp, tt.mockErr).AnyTimes()

			result, err := client.ModifyProject(tt.projectID, tt.newName)
			if (err != nil) != tt.expectErr {
				t.Errorf("ModifyProject() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ModifyProject() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestArchiveProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name      string
		projectID string
		mockResp  *resty.Response
		mockErr   error
		expected  *Project
		expectErr bool
	}{
		{
			name:      "Successful ArchiveProject request",
			projectID: "test-project-id",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"id":"project1","name":"Project 1","created_at":"2023-08-01T12:34:56Z","status":"archived"}`)),
				},
			},
			expected: &Project{
				ID:        "project1",
				Name:      "Project 1",
				CreatedAt: CustomTime(time.Date(2023, 8, 1, 12, 34, 56, 0, time.UTC)),
				Status:    "archived",
			},
			expectErr: false,
		},
		{
			name:      "ArchiveProject request with error",
			projectID: "test-project-id",
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:      "ArchiveProject request with non-200 status code",
			projectID: "test-project-id",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal server error"}`)),
				},
			},
			expectErr: true,
		},
		{
			name:      "ArchiveProject request with invalid JSON response",
			projectID: "test-project-id",
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
			mockRequest.EXPECT().Post(fmt.Sprintf("%s/%s/archive", ProjectsListEndpoint, tt.projectID)).Return(tt.mockResp, tt.mockErr).AnyTimes()

			result, err := client.ArchiveProject(tt.projectID)
			if (err != nil) != tt.expectErr {
				t.Errorf("ArchiveProject() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ArchiveProject() = %v, want %v", result, tt.expected)
			}
		})
	}
}
