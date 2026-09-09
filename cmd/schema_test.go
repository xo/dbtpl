package cmd

import (
	"testing"

	xo "github.com/xo/dbtpl/types"
)

func TestIndexFuncNameDuplicates(t *testing.T) {
	tests := []struct {
		name        string
		indexes     []xo.Index
		tableName   string
		useIndexNames bool
		expectDuplicates bool
		expectedFuncs []string
	}{
		{
			name:      "duplicate functions from primary key and unique index on same column",
			tableName: "xo_test",
			useIndexNames: false,
			indexes: []xo.Index{
				{
					Name:     "xo_test_pkey",
					IsUnique: true,
					IsPrimary: true,
					Fields: []xo.Field{
						{Name: "id"},
					},
				},
				{
					Name:     "xo_test_id_key",
					IsUnique: true,
					IsPrimary: false,
					Fields: []xo.Field{
						{Name: "id"},
					},
				},
			},
			expectDuplicates: false, // Should now be fixed
			expectedFuncs: []string{"xo_test_by_id_pk", "xo_test_by_id_unique"}, // Only conflicts get suffixes
		},
		{
			name:      "no duplicates with different columns",
			tableName: "users",
			useIndexNames: false,
			indexes: []xo.Index{
				{
					Name:     "users_pkey",
					IsUnique: true,
					IsPrimary: true,
					Fields: []xo.Field{
						{Name: "id"},
					},
				},
				{
					Name:     "users_email_key",
					IsUnique: true,
					IsPrimary: false,
					Fields: []xo.Field{
						{Name: "email"},
					},
				},
			},
			expectDuplicates: false,
			expectedFuncs: []string{"user_by_id", "user_by_email"}, // No conflicts, so no suffixes added
		},
		{
			name:      "duplicates with composite indexes on same columns",
			tableName: "user_roles",
			useIndexNames: false,
			indexes: []xo.Index{
				{
					Name:     "user_roles_pkey",
					IsUnique: true,
					IsPrimary: true,
					Fields: []xo.Field{
						{Name: "user_id"},
						{Name: "role_id"},
					},
				},
				{
					Name:     "user_roles_unique_idx",
					IsUnique: true,
					IsPrimary: false,
					Fields: []xo.Field{
						{Name: "user_id"},
						{Name: "role_id"},
					},
				},
			},
			expectDuplicates: false, // Should now be fixed
			expectedFuncs: []string{"user_role_by_user_id_role_id_pk", "user_role_by_user_id_role_id_unique"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Use the new conflict-aware function assignment
			indexes := make([]xo.Index, len(test.indexes))
			copy(indexes, test.indexes)

			assignUniqueIndexFuncNames(indexes, test.tableName, test.useIndexNames)

			funcNames := make([]string, len(indexes))
			duplicateMap := make(map[string]int)
			for i, index := range indexes {
				funcNames[i] = index.Func
				duplicateMap[index.Func]++
			}

			// Check for expected function names
			if len(test.expectedFuncs) > 0 {
				if len(funcNames) != len(test.expectedFuncs) {
					t.Errorf("Expected %d function names, got %d", len(test.expectedFuncs), len(funcNames))
				}
				for i, expected := range test.expectedFuncs {
					if i < len(funcNames) && funcNames[i] != expected {
						t.Errorf("Expected function name %d to be %q, got %q", i, expected, funcNames[i])
					}
				}
			}

			// Check for duplicates
			hasDuplicates := false
			for _, count := range duplicateMap {
				if count > 1 {
					hasDuplicates = true
					break
				}
			}

			// This test should FAIL when duplicates are detected to expose the bug
			if hasDuplicates {
				t.Fatalf("DUPLICATE FUNCTION NAMES DETECTED - this exposes the bug! Function names: %v, duplicates: %v", funcNames, duplicateMap)
			}

			// Log for debugging
			t.Logf("Generated function names: %v", funcNames)
			if hasDuplicates {
				t.Logf("Duplicate function names detected: %v", duplicateMap)
			}
		})
	}
}

// TestIndexFuncNameGeneration tests the basic functionality of indexFuncName
func TestIndexFuncNameGeneration(t *testing.T) {
	tests := []struct {
		name        string
		index       xo.Index
		tableName   string
		useIndexNames bool
		expected    string
	}{
		{
			name:      "simple unique index",
			tableName: "users",
			useIndexNames: false,
			index: xo.Index{
				Name:     "users_email_key",
				IsUnique: true,
				Fields: []xo.Field{
					{Name: "email"},
				},
			},
			expected: "user_by_email", // Backwards compatible - no suffix unless conflict
		},
		{
			name:      "non-unique index",
			tableName: "posts",
			useIndexNames: false,
			index: xo.Index{
				Name:     "posts_status_idx",
				IsUnique: false,
				Fields: []xo.Field{
					{Name: "status"},
				},
			},
			expected: "posts_by_status",
		},
		{
			name:      "composite index",
			tableName: "user_posts",
			useIndexNames: false,
			index: xo.Index{
				Name:     "user_posts_user_id_created_at_idx",
				IsUnique: false,
				Fields: []xo.Field{
					{Name: "user_id"},
					{Name: "created_at"},
				},
			},
			expected: "user_posts_by_user_id_created_at",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := indexFuncName(test.index, test.tableName, test.useIndexNames)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

// TestBackwardsCompatibility ensures that function names remain the same when no conflicts exist
func TestBackwardsCompatibility(t *testing.T) {
	tests := []struct {
		name          string
		indexes       []xo.Index
		tableName     string
		expectedFuncs []string
	}{
		{
			name:      "single unique index - no suffix added",
			tableName: "users",
			indexes: []xo.Index{
				{
					Name:     "users_email_key",
					IsUnique: true,
					Fields: []xo.Field{
						{Name: "email"},
					},
				},
			},
			expectedFuncs: []string{"user_by_email"}, // No suffix - backwards compatible
		},
		{
			name:      "single primary key - no suffix added",
			tableName: "posts",
			indexes: []xo.Index{
				{
					Name:      "posts_pkey",
					IsUnique:  true,
					IsPrimary: true,
					Fields: []xo.Field{
						{Name: "id"},
					},
				},
			},
			expectedFuncs: []string{"post_by_id"}, // No suffix - backwards compatible
		},
		{
			name:      "different columns - no suffixes added",
			tableName: "articles",
			indexes: []xo.Index{
				{
					Name:      "articles_pkey",
					IsUnique:  true,
					IsPrimary: true,
					Fields: []xo.Field{
						{Name: "id"},
					},
				},
				{
					Name:     "articles_slug_key",
					IsUnique: true,
					Fields: []xo.Field{
						{Name: "slug"},
					},
				},
				{
					Name:     "articles_author_idx",
					IsUnique: false,
					Fields: []xo.Field{
						{Name: "author_id"},
					},
				},
			},
			expectedFuncs: []string{"article_by_id", "article_by_slug", "articles_by_author_id"}, // No suffixes
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			indexes := make([]xo.Index, len(test.indexes))
			copy(indexes, test.indexes)

			assignUniqueIndexFuncNames(indexes, test.tableName, false)

			funcNames := make([]string, len(indexes))
			for i, index := range indexes {
				funcNames[i] = index.Func
			}

			if len(funcNames) != len(test.expectedFuncs) {
				t.Errorf("Expected %d function names, got %d", len(test.expectedFuncs), len(funcNames))
			}

			for i, expected := range test.expectedFuncs {
				if i < len(funcNames) && funcNames[i] != expected {
					t.Errorf("Expected function name %d to be %q, got %q", i, expected, funcNames[i])
				}
			}

			t.Logf("Generated function names: %v", funcNames)
		})
	}
}