package openaiorgs

import (
	"fmt"
	"testing"
)

func TestListUsers(t *testing.T) {
	tests := []struct {
		name             string
		mockResponse     ListResponse[User]
		mockResponseCode int
		expectedError    error
		expectedLength   int
	}{
		{
			name: "Valid response",
			mockResponse: ListResponse[User]{
				Object:  "list",
				Data:    []User{{ID: "user_123", Name: "Test User"}},
				FirstID: "user_123",
				LastID:  "user_123",
				HasMore: false,
			},
			mockResponseCode: 200,
			expectedError:    nil,
			expectedLength:   1,
		},
		// error response
		{
			name:             "Error response",
			mockResponse:     ListResponse[User]{},
			mockResponseCode: 400,
			expectedError:    fmt.Errorf("API request failed with status code 400: {\"object\":\"\",\"data\":null,\"first_id\":\"\",\"last_id\":\"\",\"has_more\":false}"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestHelper(t)
			defer h.cleanup()

			h.mockResponse("GET", "/organization/users", tt.mockResponseCode, tt.mockResponse)

			users, err := h.client.ListUsers(10, "")
			if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
			if err == nil && tt.expectedError != nil {
				t.Errorf("Expected error %v, got nil", tt.expectedError)
			} else if err == nil {
				if len(users.Data) != tt.expectedLength {
					t.Errorf("Expected %d users, got %d", tt.expectedLength, len(users.Data))
				}
			}

			h.assertRequest("GET", "/organization/users", 1)
		})
	}
}

func TestRetrieveUser(t *testing.T) {
	tests := []struct {
		name             string
		userID           string
		mockUser         User
		mockResponseCode int
		expectedError    error
	}{
		{
			name:   "Valid user",
			userID: "user_123",
			mockUser: User{
				ID:   "user_123",
				Name: "Test User",
			},
			mockResponseCode: 200,
			expectedError:    nil,
		},
		// error response
		{
			name:             "Error response",
			userID:           "user_123",
			mockResponseCode: 400,
			expectedError:    fmt.Errorf("API request failed with status code 400: {\"object\":\"\",\"id\":\"\",\"name\":\"\",\"email\":\"\",\"role\":\"\",\"added_at\":-62135596800}"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestHelper(t)
			defer h.cleanup()

			h.mockResponse("GET", "/organization/users/"+tt.userID, tt.mockResponseCode, tt.mockUser)

			user, err := h.client.RetrieveUser(tt.userID)
			if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
			if err == nil && tt.expectedError != nil {
				t.Errorf("Expected error %v, got nil", tt.expectedError)
			} else if err == nil {
				if user == nil {
					t.Error("Expected user, got nil")
				}
				if tt.mockUser.ID != user.ID {
					t.Errorf("Expected ID %s, got %s", tt.mockUser.ID, user.ID)
				}
			}

			h.assertRequest("GET", "/organization/users/"+tt.userID, 1)
		})
	}
}

func TestModifyUser(t *testing.T) {
	tests := []struct {
		name             string
		userID           string
		newRole          string
		expectedError    error
		mockResponseCode int
	}{
		{
			name:             "Valid modification",
			userID:           "user_123",
			newRole:          "owner",
			expectedError:    nil,
			mockResponseCode: 200,
		},
		// error response
		{
			name:             "Error response",
			userID:           "user_123",
			newRole:          "owner",
			expectedError:    fmt.Errorf("failed to modify user role: API request failed with status code 400: null"),
			mockResponseCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestHelper(t)
			defer h.cleanup()

			h.mockResponse("POST", "/organization/users/"+tt.userID, tt.mockResponseCode, nil)

			err := h.client.ModifyUserRole(tt.userID, ParseRoleType(tt.newRole))
			if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
			if err == nil && tt.expectedError != nil {
				t.Errorf("Expected error %v, got nil", tt.expectedError)
			}

			h.assertRequest("POST", "/organization/users/"+tt.userID, 1)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name             string
		userID           string
		mockResponseCode int
		expectedError    error
	}{
		{
			name:             "Valid deletion",
			userID:           "user_123",
			mockResponseCode: 204,
			expectedError:    nil,
		},
		// error response
		{
			name:             "Error response",
			userID:           "user_123",
			mockResponseCode: 400,
			expectedError:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestHelper(t)
			defer h.cleanup()

			h.mockResponse("DELETE", "/organization/users/"+tt.userID, tt.mockResponseCode, nil)

			err := h.client.DeleteUser(tt.userID)
			if err != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
			if err == nil && tt.expectedError != nil {
				t.Errorf("Expected error %v, got nil", tt.expectedError)
			}

			h.assertRequest("DELETE", "/organization/users/"+tt.userID, 1)
		})
	}
}
