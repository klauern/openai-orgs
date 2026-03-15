package cmd

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/urfave/cli/v3"
)

func TestPrintTableData(t *testing.T) {
	tests := []struct {
		name  string
		data  TableData
		check func(t *testing.T, output string)
	}{
		{
			name: "empty table",
			data: TableData{
				Headers: []string{"ID", "Name"},
				Rows:    [][]string{},
			},
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "ID | Name") {
					t.Errorf("Expected headers in output, got: %s", output)
				}
				lines := strings.Split(strings.TrimSpace(output), "\n")
				if len(lines) != 2 {
					t.Errorf("Expected 2 lines (header + separator) for empty table, got %d lines: %s", len(lines), output)
				}
			},
		},
		{
			name: "single row",
			data: TableData{
				Headers: []string{"ID", "Name"},
				Rows:    [][]string{{"123", "Alice"}},
			},
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "ID | Name") {
					t.Errorf("Expected headers in output, got: %s", output)
				}
				if !strings.Contains(output, "123 | Alice") {
					t.Errorf("Expected row data in output, got: %s", output)
				}
				lines := strings.Split(strings.TrimSpace(output), "\n")
				if len(lines) != 3 {
					t.Errorf("Expected 3 lines (header + separator + 1 row), got %d lines: %s", len(lines), output)
				}
			},
		},
		{
			name: "multiple rows",
			data: TableData{
				Headers: []string{"ID", "Name", "Role"},
				Rows: [][]string{
					{"123", "Alice", "owner"},
					{"456", "Bob", "member"},
					{"789", "Charlie", "reader"},
				},
			},
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "ID | Name | Role") {
					t.Errorf("Expected headers in output, got: %s", output)
				}
				if !strings.Contains(output, "123 | Alice | owner") {
					t.Errorf("Expected first row in output, got: %s", output)
				}
				if !strings.Contains(output, "789 | Charlie | reader") {
					t.Errorf("Expected last row in output, got: %s", output)
				}
				lines := strings.Split(strings.TrimSpace(output), "\n")
				if len(lines) != 5 {
					t.Errorf("Expected 5 lines (header + separator + 3 rows), got %d lines: %s", len(lines), output)
				}
			},
		},
		{
			name: "separator matches header length",
			data: TableData{
				Headers: []string{"ID", "Name"},
				Rows:    [][]string{},
			},
			check: func(t *testing.T, output string) {
				lines := strings.Split(strings.TrimSpace(output), "\n")
				if len(lines) < 2 {
					t.Fatalf("Expected at least 2 lines, got %d", len(lines))
				}
				headerLen := len(lines[0])
				separatorLen := len(lines[1])
				if headerLen != separatorLen {
					t.Errorf("Separator length (%d) should match header length (%d)", separatorLen, headerLen)
				}
				// Verify separator is all dashes
				for _, c := range lines[1] {
					if c != '-' {
						t.Errorf("Separator should only contain dashes, got: %s", lines[1])
						break
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				printTableData(tt.data)
			})
			tt.check(t, output)
		})
	}
}

func TestDefaultNewClient(t *testing.T) {
	// Build a minimal CLI command with the api-key flag
	cmd := &cli.Command{
		Name: "test",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "api-key",
				Value: "test-token",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := defaultNewClient(ctx, cmd)
			if client == nil {
				t.Fatal("expected non-nil client")
			}
			return nil
		},
	}
	err := cmd.Run(context.Background(), []string{"test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewClientFunc(t *testing.T) {
	// Verify newClientFunc is defaultNewClient by default
	cmd := &cli.Command{
		Name: "test",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "api-key",
				Value: "test-token",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := newClient(ctx, cmd)
			if client == nil {
				t.Fatal("expected non-nil client")
			}
			return nil
		},
	}
	err := cmd.Run(context.Background(), []string{"test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWrapError(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		err       error
		wantNil   bool
		wantMsg   string
	}{
		{
			name:      "nil error returns nil",
			operation: "test operation",
			err:       nil,
			wantNil:   true,
		},
		{
			name:      "non-nil error wraps with context",
			operation: "list users",
			err:       fmt.Errorf("connection refused"),
			wantMsg:   "failed to list users: connection refused",
		},
		{
			name:      "different operation context",
			operation: "delete resource",
			err:       fmt.Errorf("not found"),
			wantMsg:   "failed to delete resource: not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wrapError(tt.operation, tt.err)
			if tt.wantNil {
				if result != nil {
					t.Errorf("wrapError() = %v, want nil", result)
				}
				return
			}
			if result == nil {
				t.Errorf("wrapError() = nil, want error")
				return
			}
			if result.Error() != tt.wantMsg {
				t.Errorf("wrapError() = %q, want %q", result.Error(), tt.wantMsg)
			}
		})
	}
}

func TestBaseCommand(t *testing.T) {
	t.Run("Command returns correct structure", func(t *testing.T) {
		bc := &BaseCommand{
			Name:  "test-cmd",
			Usage: "Test command description",
			Commands: []*cli.Command{
				{Name: "sub1", Usage: "Sub command 1"},
				{Name: "sub2", Usage: "Sub command 2"},
			},
		}

		cmd := bc.Command()
		if cmd.Name != "test-cmd" {
			t.Errorf("Command().Name = %s, want test-cmd", cmd.Name)
		}
		if cmd.Usage != "Test command description" {
			t.Errorf("Command().Usage = %s, want 'Test command description'", cmd.Usage)
		}
		if len(cmd.Commands) != 2 {
			t.Errorf("Command().Commands length = %d, want 2", len(cmd.Commands))
		}
	})

	t.Run("Subcommands returns commands", func(t *testing.T) {
		subs := []*cli.Command{
			{Name: "list", Usage: "List items"},
			{Name: "create", Usage: "Create item"},
			{Name: "delete", Usage: "Delete item"},
		}
		bc := &BaseCommand{
			Name:     "resources",
			Usage:    "Manage resources",
			Commands: subs,
		}

		result := bc.Subcommands()
		if len(result) != 3 {
			t.Errorf("Subcommands() length = %d, want 3", len(result))
		}
		if result[0].Name != "list" {
			t.Errorf("Subcommands()[0].Name = %s, want list", result[0].Name)
		}
	})

	t.Run("empty commands", func(t *testing.T) {
		bc := &BaseCommand{
			Name:     "empty",
			Usage:    "Empty command",
			Commands: []*cli.Command{},
		}

		cmd := bc.Command()
		if len(cmd.Commands) != 0 {
			t.Errorf("Command().Commands length = %d, want 0", len(cmd.Commands))
		}

		subs := bc.Subcommands()
		if len(subs) != 0 {
			t.Errorf("Subcommands() length = %d, want 0", len(subs))
		}
	})
}
