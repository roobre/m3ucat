package main

import "testing"

func TestGlobalDiretivePrefixFromEnv(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name       string
		env        string
		directives map[string]bool
	}{
		{
			name: "Empty list",
			env:  "",
			directives: map[string]bool{
				"#EXAMPLE": false,
			},
		},
		{
			name: "One item",
			env:  "#TEST",
			directives: map[string]bool{
				"#EXAMPLE":      false,
				"#TEST":         true,
				"#TEST:Foo":     true,
				"#TEST:\"Foo\"": true,
				"#TESTTEST":     true,
				"#ATEST":        false,
			},
		},
		{
			name: "Several items",
			env:  "#TEST,#ATEST",
			directives: map[string]bool{
				"#EXAMPLE":      false,
				"#TEST":         true,
				"#TEST:Foo":     true,
				"#TEST:\"Foo\"": true,
				"#TESTTEST":     true,
				"#ATEST":        true,
			},
		},
	} {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fn := globalDirectivePrefixFromEnv(tc.env)
			for directive, expect := range tc.directives {
				actual := fn(directive)
				if actual != expect {
					t.Fatalf("Expected %q to return %v for %q, got %v", tc.env, expect, directive, actual)
				}
			}
		})
	}
}
