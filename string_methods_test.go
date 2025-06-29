package openaiorgs

import (
	"strings"
	"testing"
	"time"
)

func TestString_Methods(t *testing.T) {
	testTime := time.Date(2024, 3, 14, 12, 0, 0, 0, time.UTC)

	t.Run("AdminAPIKey.String", func(t *testing.T) {
		key := AdminAPIKey{
			Object:        "organization.api_key",
			ID:            "key_123",
			Name:          "Test API Key",
			CreatedAt:     UnixSeconds(testTime),
			LastUsedAt:    UnixSeconds(testTime.Add(-1 * time.Hour)),
			RedactedValue: "sk-...abc123",
			Scopes:        []string{"read", "write"},
		}

		result := key.String()
		if !strings.Contains(result, "key_123") {
			t.Errorf("Expected string to contain ID 'key_123', got: %s", result)
		}
		if !strings.Contains(result, "Test API Key") {
			t.Errorf("Expected string to contain name 'Test API Key', got: %s", result)
		}
		if !strings.Contains(result, "scopes:read,write") {
			t.Errorf("Expected string to contain scopes, got: %s", result)
		}
	})

	t.Run("User.String", func(t *testing.T) {
		user := User{
			Object:    "organization.user",
			ID:        "user_123",
			Name:      "John Doe",
			Email:     "john@example.com",
			Role:      "owner",
			AddedAt:   UnixSeconds(testTime),
		}

		result := user.String()
		if !strings.Contains(result, "user_123") {
			t.Errorf("Expected string to contain ID 'user_123', got: %s", result)
		}
		if !strings.Contains(result, "John Doe") {
			t.Errorf("Expected string to contain name 'John Doe', got: %s", result)
		}
		if !strings.Contains(result, "john@example.com") {
			t.Errorf("Expected string to contain email 'john@example.com', got: %s", result)
		}
		if !strings.Contains(result, "owner") {
			t.Errorf("Expected string to contain role 'owner', got: %s", result)
		}
	})

	t.Run("Project.String", func(t *testing.T) {
		project := Project{
			ID:        "proj_123",
			Object:    "organization.project",
			Name:      "test-project",
			CreatedAt: UnixSeconds(testTime),
			Status:    "active",
		}

		result := project.String()
		if !strings.Contains(result, "proj_123") {
			t.Errorf("Expected string to contain ID 'proj_123', got: %s", result)
		}
		if !strings.Contains(result, "test-project") {
			t.Errorf("Expected string to contain name 'test-project', got: %s", result)
		}
		if !strings.Contains(result, "active") {
			t.Errorf("Expected string to contain status 'active', got: %s", result)
		}
	})

	t.Run("Invite.String", func(t *testing.T) {
		invite := Invite{
			ObjectType: "organization.invite",
			ID:         "invite_123",
			Email:      "test@example.com",
			Role:       "reader",
			Status:     "pending",
			CreatedAt:  UnixSeconds(testTime),
			ExpiresAt:  UnixSeconds(testTime.Add(24 * time.Hour)),
		}

		result := invite.String()
		if !strings.Contains(result, "invite_123") {
			t.Errorf("Expected string to contain ID 'invite_123', got: %s", result)
		}
		if !strings.Contains(result, "test@example.com") {
			t.Errorf("Expected string to contain email 'test@example.com', got: %s", result)
		}
		if !strings.Contains(result, "reader") {
			t.Errorf("Expected string to contain role 'reader', got: %s", result)
		}
		if !strings.Contains(result, "pending") {
			t.Errorf("Expected string to contain status 'pending', got: %s", result)
		}
	})

	t.Run("Certificate.String", func(t *testing.T) {
		active := true
		cert := Certificate{
			Object:    "organization.certificate",
			ID:        "cert_123",
			Name:      "Test Certificate",
			Active:    &active,
			CreatedAt: UnixSeconds(testTime),
			CertificateDetails: CertificateDetails{
				ValidAt:   UnixSeconds(testTime),
				ExpiresAt: UnixSeconds(testTime.Add(365 * 24 * time.Hour)),
			},
		}

		result := cert.String()
		if !strings.Contains(result, "cert_123") {
			t.Errorf("Expected string to contain ID 'cert_123', got: %s", result)
		}
		if !strings.Contains(result, "Test Certificate") {
			t.Errorf("Expected string to contain name 'Test Certificate', got: %s", result)
		}
	})

	t.Run("Owner.String", func(t *testing.T) {
		owner := Owner{
			Type: OwnerTypeUser,
			User: &User{
				ID:   "user_123",
				Name: "John Doe",
			},
		}

		result := owner.String()
		if !strings.Contains(result, "user") {
			t.Errorf("Expected string to contain type 'user', got: %s", result)
		}
		if !strings.Contains(result, "user:") {
			t.Errorf("Expected string to contain user type, got: %s", result)
		}
	})

	t.Run("ProjectUser.String", func(t *testing.T) {
		projectUser := ProjectUser{
			Object:  "organization.project.user",
			ID:      "user_123",
			Name:    "John Doe",
			Email:   "john@example.com",
			Role:    "member",
			AddedAt: UnixSeconds(testTime),
		}

		result := projectUser.String()
		if !strings.Contains(result, "user_123") {
			t.Errorf("Expected string to contain ID 'user_123', got: %s", result)
		}
		if !strings.Contains(result, "John Doe") {
			t.Errorf("Expected string to contain name 'John Doe', got: %s", result)
		}
		if !strings.Contains(result, "member") {
			t.Errorf("Expected string to contain role 'member', got: %s", result)
		}
	})

	t.Run("ProjectServiceAccount.String", func(t *testing.T) {
		serviceAccount := ProjectServiceAccount{
			Object:     "organization.project.service_account",
			ID:         "sa_123",
			Name:       "Test Service Account",
			Role:       "member",
			CreatedAt:  UnixSeconds(testTime),
		}

		result := serviceAccount.String()
		if !strings.Contains(result, "sa_123") {
			t.Errorf("Expected string to contain ID 'sa_123', got: %s", result)
		}
		if !strings.Contains(result, "Test Service Account") {
			t.Errorf("Expected string to contain name 'Test Service Account', got: %s", result)
		}
		if !strings.Contains(result, "member") {
			t.Errorf("Expected string to contain role 'member', got: %s", result)
		}
	})

	t.Run("ProjectApiKey.String", func(t *testing.T) {
		apiKey := ProjectApiKey{
			Object:        "organization.project.api_key",
			Name:          "Test API Key",
			CreatedAt:     UnixSeconds(testTime),
			RedactedValue: "sk-...xyz789",
		}

		result := apiKey.String()
		if !strings.Contains(result, "Test API Key") {
			t.Errorf("Expected string to contain name 'Test API Key', got: %s", result)
		}
		if !strings.Contains(result, "no owner") {
			t.Errorf("Expected string to contain owner info, got: %s", result)
		}
	})

	t.Run("ProjectRateLimit.String", func(t *testing.T) {
		rateLimit := ProjectRateLimit{
			Object:                      "project_rate_limit",
			ID:                          "rl_123",
			Model:                       "gpt-4",
			MaxRequestsPer1Minute:       100,
			MaxTokensPer1Minute:         50000,
			MaxImagesPer1Minute:         10,
			MaxAudioMegabytesPer1Minute: 25,
			MaxRequestsPer1Day:          1000,
			Batch1DayMaxInputTokens:     1000000,
		}

		result := rateLimit.String()
		if !strings.Contains(result, "rl_123") {
			t.Errorf("Expected string to contain ID 'rl_123', got: %s", result)
		}
		if !strings.Contains(result, "gpt-4") {
			t.Errorf("Expected string to contain model 'gpt-4', got: %s", result)
		}
		if !strings.Contains(result, "100") {
			t.Errorf("Expected string to contain requests limit '100', got: %s", result)
		}
	})
}