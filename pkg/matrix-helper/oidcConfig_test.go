package matrixhelper

import (
	"fmt"
	"strings"
	"testing"
)

// realisticDoc is a realistic user-keycloak config document with the
// client_secret line in the %s position. It contains the sibling keys the
// sync must never touch (password auth, telemetry, provider settings).
const realisticDoc = `oidc:
  providers:
  - id: keycloak
    alias: Keycloak
    issuer: https://keycloak.matrix.buc.sh/realms/buc
    client_id: matrix
    client_secret: %s
    scopes:
    - openid
    - email
    discover: true
password:
  enabled: true
  local_auth_enabled: true
telemetry:
  enabled: false
`

func Test_ensureOidcClientSecret(t *testing.T) {
	tests := []struct {
		name       string
		config     string
		source     string
		wantConfig string
		wantChange bool
		wantErr    bool
	}{
		// Decision table: an unchanged value does not apply, a differing
		// value does (first start with the Helm placeholder, rotation, and
		// manual drift).
		{"unchanged_unquoted", "client_secret: abc123", "abc123", "client_secret: abc123", false, false},
		{"unchanged_quoted", `client_secret: "abc123"`, "abc123", `client_secret: "abc123"`, false, false},
		{"rotation_applies", "client_secret: oldsecret", "newsecret", `client_secret: "newsecret"`, true, false},
		{"manual_drift_is_corrected", "client_secret: manually-edited", "livesecret123", `client_secret: "livesecret123"`, true, false},
		{"first_start_with_helm_placeholder_applies", fmt.Sprintf(realisticDoc, "__MATRIX_HELPER_CLIENT_SECRET__"), "livesecret123", fmt.Sprintf(realisticDoc, `"livesecret123"`), true, false},

		// Current-value forms: every form the old extract matrix exercised
		// is detected as unchanged when the source matches the extracted
		// value.
		{"current_unquoted_single_character", "client_secret: a", "a", "client_secret: a", false, false},
		{"current_unquoted_with_indentation", "  client_secret: abc123", "abc123", "  client_secret: abc123", false, false},
		{"current_unquoted_with_deep_indentation", "        client_secret: abc123", "abc123", "        client_secret: abc123", false, false},
		{"current_unquoted_with_tab_indentation", "\tclient_secret: abc123", "abc123", "\tclient_secret: abc123", false, false},
		{"current_unquoted_with_mixed_indentation", " \t client_secret: abc123", "abc123", " \t client_secret: abc123", false, false},
		{"current_unquoted_tab_after_colon", "client_secret:\tabc123", "abc123", "client_secret:\tabc123", false, false},
		{"current_double_quoted_with_indentation", `  client_secret: "abc123"`, "abc123", `  client_secret: "abc123"`, false, false},
		{"current_double_quoted_with_hash", `client_secret: "a#b"`, "a#b", `client_secret: "a#b"`, false, false},
		{"current_trailing_comment_unquoted", "client_secret: abc 123 # rotated 2026-01-01", "abc", "client_secret: abc 123 # rotated 2026-01-01", false, false},
		{"current_trailing_comment_quoted", `client_secret: "abc123" # rotated 2026-01-01`, "abc123", `client_secret: "abc123" # rotated 2026-01-01`, false, false},
		{"current_crlf_unquoted", "client_secret: abc123\r\n", "abc123", "client_secret: abc123\r\n", false, false},
		{"current_crlf_quoted", "client_secret: \"abc123\"\r\n", "abc123", "client_secret: \"abc123\"\r\n", false, false},
		{"current_sibling_key_foo_client_secret", "foo_client_secret: other\nclient_secret: abc123", "abc123", "foo_client_secret: other\nclient_secret: abc123", false, false},
		{"current_sibling_key_client_secret_foo", "client_secret_foo: other\nclient_secret: abc123", "abc123", "client_secret_foo: other\nclient_secret: abc123", false, false},
		{"current_sibling_key_not_client_secret", "not_client_secret: other\nclient_secret: abc123", "abc123", "not_client_secret: other\nclient_secret: abc123", false, false},
		{"current_value_with_supported_special_chars", "client_secret: a-b_c.d~+/=@%$&*()[]{}", "a-b_c.d~+/=@%$&*()[]{}", "client_secret: a-b_c.d~+/=@%$&*()[]{}", false, false},
		{"current_whitespace_stops_the_value", "client_secret: abc 123", "abc", "client_secret: abc 123", false, false},
		{"current_hash_stops_the_value", "client_secret: abc#def", "abc", "client_secret: abc#def", false, false},

		// Replacement: the output form and byte preservation.
		{"replace_unquoted_becomes_double_quoted", "client_secret: old\nissuer: x", "newsecret", "client_secret: \"newsecret\"\nissuer: x", true, false},
		{"replace_quoted_value_for_value", `client_secret: "old"`, "newsecret", `client_secret: "newsecret"`, true, false},
		{"replace_trailing_comment_preserved", "client_secret: old # comment", "newsecret", `client_secret: "newsecret" # comment`, true, false},
		{"replace_indentation_preserved", "  client_secret: old", "newsecret", "  client_secret: \"newsecret\"", true, false},
		{"replace_tab_indentation_preserved", "\tclient_secret: old", "newsecret", "\tclient_secret: \"newsecret\"", true, false},
		{"replace_crlf_ending_preserved", "client_secret: old\r\nissuer: x", "newsecret", "client_secret: \"newsecret\"\r\nissuer: x", true, false},
		{"replace_sibling_keys_untouched", "foo_client_secret: keep\nclient_secret: old\nclient_secret_foo: keep", "newsecret", "foo_client_secret: keep\nclient_secret: \"newsecret\"\nclient_secret_foo: keep", true, false},
		{"replace_empty_quoted_value", `client_secret: ""`, "abc123", `client_secret: "abc123"`, true, false},
		// A "$" in the new secret must be written literally, not expanded
		// as a regex replacement reference.
		{"replace_dollar_sign_is_literal", "client_secret: old", "a$1b", `client_secret: "a$1b"`, true, false},
		{"replace_realistic_document_other_bytes_preserved", fmt.Sprintf(realisticDoc, "old"), "newsecret", fmt.Sprintf(realisticDoc, `"newsecret"`), true, false},
		{"replace_realistic_document_crlf_other_bytes_preserved", strings.ReplaceAll(fmt.Sprintf(realisticDoc, "old"), "\n", "\r\n"), "newsecret", strings.ReplaceAll(fmt.Sprintf(realisticDoc, `"newsecret"`), "\n", "\r\n"), true, false},

		// Config without exactly one client_secret line.
		{"missing_key", "issuer: https://keycloak\nclient_id: matrix", "abc123", "", false, true},
		{"missing_key_empty_document", "", "abc123", "", false, true},
		{"only_sibling_key", "foo_client_secret: other", "abc123", "", false, true},
		{"missing_value", "client_secret:", "abc123", "", false, true},
		{"missing_value_trailing_spaces", "client_secret:   ", "abc123", "", false, true},
		{"missing_value_comment_only", "client_secret: # no value yet", "abc123", "", false, true},
		{"unterminated_quote", `client_secret: "abc123`, "abc123", "", false, true},
		{"duplicated_keys", "client_secret: abc\nclient_secret: def", "abc123", "", false, true},
		{"duplicated_keys_quoted_and_unquoted", "client_secret: abc\nclient_secret: \"def\"", "abc123", "", false, true},

		// Invalid source: an error, not a panic — even when it equals the
		// current value.
		{"empty_source_rejected", "client_secret: old", "", "", false, true},
		{"source_with_space_rejected", "client_secret: old", "a b", "", false, true},
		{"source_with_space_rejected_even_when_current_matches", `client_secret: "a b"`, "a b", "", false, true},
		{"source_with_quote_rejected", "client_secret: old", `a"b`, "", false, true},
		{"source_with_newline_rejected", "client_secret: old", "a\nb", "", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConfig, gotChanged, err := ensureOidcClientSecret(tt.config, tt.source)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got config %q changed=%v", gotConfig, gotChanged)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if gotChanged != tt.wantChange {
				t.Fatalf("expected changed=%v, got %v", tt.wantChange, gotChanged)
			}
			if gotConfig != tt.wantConfig {
				t.Fatalf("expected config %q, got %q", tt.wantConfig, gotConfig)
			}
		})
	}
}

// Test_ensureOidcClientSecret_fixedpoint pins that a replaced config is a
// fixed point: re-running with the same source reports no change, so the
// sync converges after a single apply.
func Test_ensureOidcClientSecret_fixedpoint(t *testing.T) {
	for _, doc := range []string{
		fmt.Sprintf(realisticDoc, "old"),
		fmt.Sprintf(realisticDoc, `"old"`),
		strings.ReplaceAll(fmt.Sprintf(realisticDoc, "old"), "\n", "\r\n"),
	} {
		updated, changed, err := ensureOidcClientSecret(doc, "newsecret")
		if err != nil {
			t.Fatalf("unexpected error replacing: %s", err)
		}
		if !changed {
			t.Fatalf("expected the first run to change the config")
		}
		stable, changed, err := ensureOidcClientSecret(updated, "newsecret")
		if err != nil {
			t.Fatalf("unexpected error re-running: %s", err)
		}
		if changed {
			t.Fatalf("expected the second run to be a no-op")
		}
		if stable != updated {
			t.Fatalf("expected the config to be stable, got %q", stable)
		}
	}
}

func Test_validateClientSecret(t *testing.T) {
	tests := []struct {
		name    string
		secret  string
		wantErr bool
	}{
		{"valid_simple", "abc123", false},
		{"valid_single_character", "a", false},
		{"valid_special_chars", "a-b_c.d~+/=@%$&*()[]{}", false},
		{"valid_helm_placeholder", "__MATRIX_HELPER_CLIENT_SECRET__", false},
		{"empty_rejected", "", true},
		{"space_rejected", "a b", true},
		{"tab_rejected", "a\tb", true},
		{"newline_rejected", "a\nb", true},
		{"quote_rejected", `a"b`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateClientSecret(tt.secret)
			if tt.wantErr && err == nil {
				t.Fatalf("expected an error for %q", tt.secret)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error for %q: %s", tt.secret, err)
			}
		})
	}
}
