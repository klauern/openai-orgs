package openaiorgs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/go-resty/resty/v2"
)

func TestListInvites(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name      string
		mockResp  *resty.Response
		mockErr   error
		expected  []Invite
		expectErr bool
	}{
		{
			name: "Successful ListInvites",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"object":"list","data":[{"id":"invite1","email":"test@example.com","role":"member","status":"pending"}]}`)),
				},
			},
			expected: []Invite{
				{ID: "invite1", Email: "test@example.com", Role: "member", Status: "pending"},
			},
			expectErr: false,
		},
		{
			name:      "ListInvites with error",
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name: "ListInvites with non-200 status code",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal server error"}`)),
				},
			},
			expectErr: true,
		},
		{
			name: "ListInvites with invalid JSON response",
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
			mockRequest.EXPECT().Get(InviteListEndpoint).Return(tt.mockResp, tt.mockErr).AnyTimes()

			result, err := client.ListInvites()
			if (err != nil) != tt.expectErr {
				t.Errorf("ListInvites() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ListInvites() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCreateInvite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name      string
		email     string
		role      RoleType
		mockResp  *resty.Response
		mockErr   error
		expected  *Invite
		expectErr bool
	}{
		{
			name:  "Successful CreateInvite",
			email: "test@example.com",
			role:  RoleTypeMember,
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"id":"invite1","email":"test@example.com","role":"member","status":"pending"}`)),
				},
			},
			expected:  &Invite{ID: "invite1", Email: "test@example.com", Role: "member", Status: "pending"},
			expectErr: false,
		},
		{
			name:      "CreateInvite with error",
			email:     "test@example.com",
			role:      RoleTypeMember,
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name:  "CreateInvite with non-200 status code",
			email: "test@example.com",
			role:  RoleTypeMember,
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal server error"}`)),
				},
			},
			expectErr: true,
		},
		{
			name:  "CreateInvite with invalid JSON response",
			email: "test@example.com",
			role:  RoleTypeMember,
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
			mockRequest.EXPECT().SetBody(map[string]string{"email": tt.email, "role": string(tt.role)}).Return(mockRequest).AnyTimes()
			mockRequest.EXPECT().Post(InviteListEndpoint).Return(tt.mockResp, tt.mockErr).AnyTimes()

			result, err := client.CreateInvite(tt.email, tt.role)
			if (err != nil) != tt.expectErr {
				t.Errorf("CreateInvite() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("CreateInvite() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRetrieveInvite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name      string
		id        string
		mockResp  *resty.Response
		mockErr   error
		expected  *Invite
		expectErr bool
	}{
		{
			name: "Successful RetrieveInvite",
			id:   "invite1",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"object":"list","data":[{"id":"invite1","email":"test@example.com","role":"member","status":"pending"}]}`)),
				},
			},
			expected:  &Invite{ID: "invite1", Email: "test@example.com", Role: "member", Status: "pending"},
			expectErr: false,
		},
		{
			name:      "RetrieveInvite with error",
			id:        "invite1",
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name: "RetrieveInvite with non-200 status code",
			id:   "invite1",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{"error":"internal server error"}`)),
				},
			},
			expectErr: true,
		},
		{
			name: "RetrieveInvite with invalid JSON response",
			id:   "invite1",
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
			mockRequest.EXPECT().Get(InviteListEndpoint + "/" + tt.id).Return(tt.mockResp, tt.mockErr).AnyTimes()

			result, err := client.RetrieveInvite(tt.id)
			if (err != nil) != tt.expectErr {
				t.Errorf("RetrieveInvite() error = %v, wantErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RetrieveInvite() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDeleteInvite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	client := &Client{client: mockClient}

	tests := []struct {
		name      string
		id        string
		mockResp  *resty.Response
		mockErr   error
		expectErr bool
	}{
		{
			name: "Successful DeleteInvite",
			id:   "invite1",
			mockResp: &resty.Response{
				RawResponse: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(`{}`)),
				},
			},
			expectErr: false,
		},
		{
			name:      "DeleteInvite with error",
			id:        "invite1",
			mockErr:   fmt.Errorf("request error"),
			expectErr: true,
		},
		{
			name: "DeleteInvite with non-200 status code",
			id:   "invite1",
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
			mockRequest.EXPECT().Delete(InviteListEndpoint + "/" + tt.id).Return(tt.mockResp, tt.mockErr).AnyTimes()

			err := client.DeleteInvite(tt.id)
			if (err != nil) != tt.expectErr {
				t.Errorf("DeleteInvite() error = %v, wantErr %v", err, tt.expectErr)
			}
		})
	}
}
