// Package specpolicy ports the r19 prose/history assertions, separate from wire semantics.
package specpolicy

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Python's Unicode whitespace includes the four information separators.
func pythonSpace(r rune) bool { return unicode.IsSpace(r) || r >= 0x1c && r <= 0x1f }
func matches(pattern, s string) bool {
	pattern = strings.ReplaceAll(pattern, `\s`, `[\t\n\v\f\r\x1c-\x1f\x85\p{Zs}\x{2028}\x{2029}]`)
	return regexp.MustCompile(pattern).MatchString(s)
}
func flat(s string) string { return strings.Join(strings.FieldsFunc(s, pythonSpace), " ") }
func decimal(s string) (int, error) {
	// Python's Unicode \d and int() recognize decimal digits beyond ASCII.
	var b strings.Builder
	for _, r := range s {
		digit := -1
		for _, rr := range unicode.Digit.R16 {
			if uint32(r) >= uint32(rr.Lo) && uint32(r) <= uint32(rr.Hi) && (uint32(r)-uint32(rr.Lo))%uint32(rr.Stride) == 0 {
				digit = int((uint32(r)-uint32(rr.Lo))/uint32(rr.Stride)) % 10
				break
			}
		}
		if digit < 0 {
			for _, rr := range unicode.Digit.R32 {
				if uint32(r) >= rr.Lo && uint32(r) <= rr.Hi && (uint32(r)-rr.Lo)%rr.Stride == 0 {
					digit = int((uint32(r)-rr.Lo)/rr.Stride) % 10
					break
				}
			}
		}
		if digit < 0 {
			return 0, fmt.Errorf("invalid decimal digit")
		}
		b.WriteByte(byte('0' + digit))
	}
	return strconv.Atoi(b.String())
}
func section(s string, n int) string {
	p := strings.SplitN(s, fmt.Sprintf("## %d.", n), 2)
	if len(p) != 2 {
		return ""
	}
	return strings.SplitN(p[1], fmt.Sprintf("## %d.", n+1), 2)[0]
}
func Check(root string) error {
	b, e := os.ReadFile(filepath.Join(root, "spec/totipo-vault-format-v1.md"))
	if e != nil {
		return e
	}
	a, e := os.ReadFile(filepath.Join(root, "review/V1_PRE_R16_REVISION_HISTORY.md"))
	if e != nil {
		return e
	}
	return Validate(string(b), string(a))
}
func Validate(s, archive string) error {
	if !utf8.ValidString(s) || !utf8.ValidString(archive) {
		return fmt.Errorf("spec policy: invalid UTF-8")
	}
	// pathlib.read_text uses universal newline translation.
	normalize := func(s string) string { return strings.ReplaceAll(strings.ReplaceAll(s, "\r\n", "\n"), "\r", "\n") }
	s, archive = normalize(s), normalize(archive)
	parts := strings.SplitN(s, "## 21. Revision history", 2)
	if len(parts) != 2 {
		return fmt.Errorf("spec policy: revision history split")
	}
	norm, history := parts[0], parts[1]
	sections := regexp.MustCompile(`(?m)^## (\p{Nd}+)\.`).FindAllStringSubmatch(s, -1)
	if len(sections) != 21 {
		return fmt.Errorf("spec policy: section count")
	}
	for i, m := range sections {
		n, e := decimal(m[1])
		if e != nil || n != i+1 {
			return fmt.Errorf("spec policy: section order")
		}
	}

	if !(strings.Contains(norm, "**Revision:** r19") && strings.Contains(norm, "**Protocol version:** 1")) {
		return fmt.Errorf("spec policy: %s", "'**Revision:** r19' in norm and '**Protocol version:** 1' in norm")
	}
	if !(strings.Contains(history, "../review/V1_PRE_R16_REVISION_HISTORY.md")) {
		return fmt.Errorf("spec policy: %s", "'../review/V1_PRE_R16_REVISION_HISTORY.md' in history")
	}
	for n := 1; n < 20; n++ {
		where := history
		if n < 16 {
			where = archive
		}
		if !strings.Contains(where, fmt.Sprintf("### v1/r%d\n", n)) {
			return fmt.Errorf("spec policy: revision r%d", n)
		}
	}
	for _, term := range []string{"DEVICE", "DEVICE_ID", "P-256", "ECDSA", "DER", "SIGNATURE", "AUTHOR_DEVICE_ID", "AUTHOR_TIME", "PROVENANCE", "OPAQUE_ROUTABLE", "OPAQUE_UNSCOPED", "OBJECT_TYPE", "OBJECT_VERSION", "VAULT_BINDING", "READY", "PROCESSING_INCOMPLETE", "HISTORY_REGRESSION", "HISTORY_MEMORY_LOST", "candidate-use", "K_signature_context"} {
		re := regexp.MustCompile(regexp.QuoteMeta(term))
		for _, loc := range re.FindAllStringIndex(norm, -1) {
			boundary := func(b byte) bool { return b >= 'A' && b <= 'Z' || b == '_' }
			if (loc[0] == 0 || !boundary(norm[loc[0]-1])) && (loc[1] == len(norm) || !boundary(norm[loc[1]])) {
				return fmt.Errorf("spec policy: retired term %s", term)
			}
		}
	}
	if strings.Contains(strings.ToLower(norm), "provenance") {
		return fmt.Errorf("spec policy: %s", "'provenance' not in norm.lower()")
	}
	if strings.Contains(strings.ToLower(norm), "opaque") {
		return fmt.Errorf("spec policy: %s", "'opaque' not in norm.lower()")
	}
	for _, anchor := range []string{"TOKEN is the sole semantic object grammar", "MAX_PARENTS = 4", "0 <= PARENT_COUNT <= 4", "1005", "1006", "1024", "VAULT_ID = SHA-256(exact canonical VAULT representation)", "A ≡c B iff A reaches B and B reaches A", "Totipo v1 protocol state is create-only after vault initialization", "Synchronization is\noptional and external", "MUST NOT require exhaustive enumeration", "Parent availability is not a publication prerequisite"} {
		if !(strings.Contains(norm, anchor)) {
			return fmt.Errorf("spec policy: %s", anchor)
		}
	}
	registry := regexp.MustCompile(`(?m)^\| (0x[0-9a-f]{4}) \| ([A-Z_]+) \|`).FindAllStringSubmatch(norm, -1)
	names := []string{"TOKEN_ID", "PARENT_COUNT", "PARENT_ID", "STATUS", "ISSUER", "ACCOUNT", "ALGORITHM", "DIGITS", "PERIOD", "SECRET_BYTES", "CLIENT_NAME", "CLIENT_TIME"}
	if len(registry) != len(names) {
		return fmt.Errorf("spec policy: tag registry count")
	}
	for i, m := range registry {
		if m[1] != fmt.Sprintf("0x%04x", i+1) || m[2] != names[i] {
			return fmt.Errorf("spec policy: tag registry order")
		}
	}
	for _, m := range regexp.MustCompile(`Section (\p{Nd}+)`).FindAllStringSubmatch(norm, -1) {
		n, e := decimal(m[1])
		if e != nil || n < 1 || n > 20 {
			return fmt.Errorf("spec policy: section reference %s", m[1])
		}
	}
	layoutMatch := regexp.MustCompile("(?s)```text\nconfigured Totipo store\n(.*?)```").FindStringSubmatch(norm)
	if len(layoutMatch) != 2 {
		return fmt.Errorf("spec policy: layout block")
	}
	layout := layoutMatch[1]
	if !(strings.TrimFunc(strings.SplitN(layout, "\n", 2)[0], pythonSpace) == "vault") {
		return fmt.Errorf("spec policy: %s", "layout bootstrap path")
	}
	if matches("(?m)^\\s*VAULT\\s*$", layout) {
		return fmt.Errorf("spec policy: %s", "not re.search(r'^\\s*VAULT\\s*$', layout, re.M)")
	}
	if matches("`VAULT`[^.\\n]*(?:pathname|path is)", norm) {
		return fmt.Errorf("spec policy: %s", "not re.search(r'`VAULT`[^.\\n]*(?:pathname|path is)', norm)")
	}
	fold := section(norm, 17)
	preserved := map[string]string{}
	for _, m := range regexp.MustCompile("(?m)^\\| (TokenValue|CLIENT_NAME|CLIENT_TIME) \\| ([^|]+) \\|$").FindAllStringSubmatch(fold, -1) {
		preserved[m[1]] = m[2]
	}
	if !(len(preserved) == 3 && preserved["TokenValue"] != "" && preserved["CLIENT_NAME"] != "" && preserved["CLIENT_TIME"] != "") {
		return fmt.Errorf("spec policy: %s", "fold table")
	}
	for _, field := range []string{"CLIENT_NAME", "CLIENT_TIME"} {
		if !strings.Contains(preserved[field], "same exact presence") {
			return fmt.Errorf("spec policy: fold metadata %s", field)
		}
	}
	values := section(norm, 16)
	tupleMatch := regexp.MustCompile("(?s)TokenValue is the complete tuple `([^`]+)`").FindStringSubmatch(values)
	if len(tupleMatch) != 2 {
		return fmt.Errorf("spec policy: TokenValue tuple")
	}
	valueTuple := tupleMatch[1]
	if strings.Contains(valueTuple, "CLIENT_NAME") && !strings.Contains(valueTuple, "CLIENT_TIME") {
		return fmt.Errorf("spec policy: %s", "'CLIENT_NAME' not in valueTuple and 'CLIENT_TIME' not in valueTuple")
	}
	heads := section(norm, 15)
	for _, field := range []string{"CLIENT_NAME", "CLIENT_TIME", "TokenValue"} {
		if !strings.Contains(heads, field) {
			return fmt.Errorf("spec policy: heads %s", field)
		}
	}
	if !(matches("MUST NOT erase metadata from head objects", heads)) {
		return fmt.Errorf("spec policy: %s", "re.search(r'MUST NOT erase metadata from head objects', heads)")
	}
	store := section(norm, 3)
	entryTypes := map[string]string{}
	for _, m := range regexp.MustCompile("(?m)^\\| `(vault|objects-v1)` \\| ([^|]+) \\|").FindAllStringSubmatch(store, -1) {
		entryTypes[m[1]] = m[2]
	}
	if !(len(entryTypes) == 2 && entryTypes["vault"] == "regular file only" && entryTypes["objects-v1"] == "directory only") {
		return fmt.Errorf("spec policy: %s", "entry types")
	}
	flatStore := flat(store)
	if !(matches("symlink.*?MUST NOT be deliberately followed or interpreted", flatStore)) {
		return fmt.Errorf("spec policy: %s", "re.search(r'symlink.*?MUST NOT be deliberately followed or interpreted', flatStore)")
	}
	if !(matches("symlink.*?MUST NOT be deliberately traversed", flatStore)) {
		return fmt.Errorf("spec policy: %s", "re.search(r'symlink.*?MUST NOT be deliberately traversed', flatStore)")
	}
	for _, syscall := range []string{"openat2", "O_NOFOLLOW", "O_PATH", "statx", "SecureDirectoryStream"} {
		if strings.Contains(norm, syscall) {
			return fmt.Errorf("spec policy: %s", syscall)
		}
	}
	if !(strings.Contains(norm, "optional implementation hardening")) {
		return fmt.Errorf("spec policy: %s", "'optional implementation hardening' in norm")
	}
	metadata := flat(section(norm, 12))
	if !(matches("For as long as.*?represented, returned, exposed, or retained.*?MUST be preserved exactly", metadata)) {
		return fmt.Errorf("spec policy: %s", "re.search(r'For as long as.*?represented, returned, exposed, or retained.*?MUST be preserved exactly', metadata)")
	}
	if !(strings.Contains(metadata, "does not require persistent retention")) {
		return fmt.Errorf("spec policy: %s", "'does not require persistent retention' in metadata")
	}
	if !(strings.Contains(metadata, "no advisory-history, remembered-head, protocol-cache, or cross-run persistence requirement")) {
		return fmt.Errorf("spec policy: %s", "'no advisory-history, remembered-head, protocol-cache, or cross-run persistence requirement' in metadata")
	}
	writing := flat(fold)
	if !(matches("Unavailability of history alone MUST NOT prohibit authorship when the TOKEN_ID and complete desired TokenValue.*?known or supplied", writing)) {
		return fmt.Errorf("spec policy: %s", "re.search(r'Unavailability of history alone MUST NOT prohibit authorship when the TOKEN_ID and complete desired TokenValue.*?known or supplied', writing)")
	}
	if !(strings.Contains(writing, "generate a new CSPRNG TOKEN_ID") && strings.Contains(writing, "empty parent set")) {
		return fmt.Errorf("spec policy: %s", "'generate a new CSPRNG TOKEN_ID' in writing and 'empty parent set' in writing")
	}
	password := flat(section(norm, 5))
	if !(strings.Contains(password, "interactive application MUST require explicit confirmation")) {
		return fmt.Errorf("spec policy: %s", "'interactive application MUST require explicit confirmation' in password")
	}
	if !(strings.Contains(password, "before creating a vault with an empty password")) {
		return fmt.Errorf("spec policy: %s", "'before creating a vault with an empty password' in password")
	}
	creation := flat(section(norm, 7))
	for _, anchor := range []string{"ordinary vault creation MUST NOT silently initialize", "MUST require an explicit recovery, reconfiguration, or new-location workflow", "MUST NOT be classified as garbage", "Creation MUST NOT delete them", "Reopen canonical `vault` as a regular file", "NO-REPLACE", "Report success only after all required publication, persistence, and revalidation work succeeds"} {
		if !(strings.Contains(creation, anchor)) {
			return fmt.Errorf("spec policy: %s", anchor)
		}
	}
	if !(matches("MUST NOT describe tombstoning or deleting a token as securely erasing", flat(section(norm, 16)))) {
		return fmt.Errorf("spec policy: %s", "re.search(r'MUST NOT describe tombstoning or deleting a token as securely erasing', section_text(16))")
	}
	changes := flat(section(norm, 8))
	if !(strings.Contains(changes, "v1 has no in-place password change and defines no password-rewrap operation")) {
		return fmt.Errorf("spec policy: %s", "'v1 has no in-place password change and defines no password-rewrap operation' in changes")
	}
	if !(strings.Contains(changes, "Reusing K_root for migration is forbidden")) {
		return fmt.Errorf("spec policy: %s", "'Reusing K_root for migration is forbidden' in changes")
	}
	if !(strings.Contains(changes, "INFORMATIVE")) {
		return fmt.Errorf("spec policy: %s", "'INFORMATIVE' in changes")
	}
	if !(strings.Contains(changes, "may omit not-yet- synchronized data")) {
		return fmt.Errorf("spec policy: %s", "'may omit not-yet- synchronized data' in changes")
	}
	for _, retired := range []string{"VAULT_FINGERPRINT", "totipo/v1/vault-fingerprint", "CURRENT != BASE", "stageReplacement", "compare-before-replace", "latest password wrapper"} {
		if strings.Contains(norm, retired) {
			return fmt.Errorf("spec policy: %s", retired)
		}
	}
	root := flat(section(norm, 6))
	for _, anchor := range []string{"K_root = CSPRNG(32 bytes)", "immutable for the complete lifetime", "MUST NOT be rotated in-place", "BOOTSTRAP_VERSION = 0x01", "NEVER a newer version", "same-root alternate legitimate bootstrap representations"} {
		if !(strings.Contains(root, anchor)) {
			return fmt.Errorf("spec policy: %s", anchor)
		}
	}
	if !(strings.Contains(creation, "MUST NOT replace it, even if it is exact-identical")) {
		return fmt.Errorf("spec policy: %s", "'MUST NOT replace it, even if it is exact-identical' in creation")
	}
	scopes := flat(section(norm, 20))
	for _, scope := range []string{"Core protocol", "Store/writer", "Application"} {
		if !strings.Contains(scopes, "| "+scope+" |") {
			return fmt.Errorf("spec policy: scope %s", scope)
		}
	}
	if !(matches("\\| Core protocol \\|[^|]*metadata preservation", scopes)) {
		return fmt.Errorf("spec policy: %s", "re.search(r'\\| Core protocol \\|[^|]*metadata preservation', scopes)")
	}
	for _, anchor := range []string{"per-object metadata rule in Section 12", "Section 17 additionally governs operation-wide fold metadata", "MUST satisfy the applicable core requirements", "core library need not guarantee the behavior of a GUI", "portable corpus success alone does not establish application conformance", "No stronger atomic primitive is imposed by this scope"} {
		if !(strings.Contains(scopes, anchor)) {
			return fmt.Errorf("spec policy: %s", anchor)
		}
	}
	return nil
}
