package scope_test

import (
	"compass/scope"
	"testing"
)

func TestCreatePermission(t *testing.T) {
	tests := []struct {
		name     string
		conf     byte
		expected scope.Permission
	}{
		{
			name: "No permissions set",
			conf: 0,
			expected: scope.Permission{
				All:    false,
				Access: false,
				Create: false,
				Delete: false,
				Get:    false,
				List:   false,
				Modify: false,
			},
		},
		{
			name: "All permissions set",
			conf: scope.S_ALL | scope.S_ACCESS | scope.S_CREATE | scope.S_DELETE | scope.S_GET | scope.S_LIST | scope.S_MODIFY,
			expected: scope.Permission{
				All:    true,
				Access: true,
				Create: true,
				Delete: true,
				Get:    true,
				List:   true,
				Modify: true,
			},
		},
		{
			name: "Only Access permission set",
			conf: scope.S_ACCESS,
			expected: scope.Permission{
				All:    false,
				Access: true,
				Create: false,
				Delete: false,
				Get:    false,
				List:   false,
				Modify: false,
			},
		},
		{
			name: "Create and Modify permissions set",
			conf: scope.S_CREATE | scope.S_MODIFY,
			expected: scope.Permission{
				All:    false,
				Access: false,
				Create: true,
				Delete: false,
				Get:    false,
				List:   false,
				Modify: true,
			},
		},
		{
			name: "Access, Delete, and List permissions set",
			conf: scope.S_ACCESS | scope.S_DELETE | scope.S_LIST,
			expected: scope.Permission{
				All:    false,
				Access: true,
				Create: false,
				Delete: true,
				Get:    false,
				List:   true,
				Modify: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scope.CreatePermission(tt.conf)

			if result.All != tt.expected.All {
				t.Errorf("expected All = %v, got %v", tt.expected.All, result.All)
			}
			if result.Access != tt.expected.Access {
				t.Errorf("expected Access = %v, got %v", tt.expected.Access, result.Access)
			}
			if result.Create != tt.expected.Create {
				t.Errorf("expected Create = %v, got %v", tt.expected.Create, result.Create)
			}
			if result.Delete != tt.expected.Delete {
				t.Errorf("expected Delete = %v, got %v", tt.expected.Delete, result.Delete)
			}
			if result.Get != tt.expected.Get {
				t.Errorf("expected Get = %v, got %v", tt.expected.Get, result.Get)
			}
			if result.List != tt.expected.List {
				t.Errorf("expected List = %v, got %v", tt.expected.List, result.List)
			}
			if result.Modify != tt.expected.Modify {
				t.Errorf("expected Modify = %v, got %v", tt.expected.Modify, result.Modify)
			}
		})
	}
}
