package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

// ─────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────

// makeBannerFile writes a minimal banner file to a temp path and returns it.
// Each character block is 8 lines of the character letter repeated, separated by a blank line.
func makeBannerFile(t *testing.T) string {
	t.Helper()
	var sb strings.Builder
	// Build entries for space (32) through tilde (126) — 95 characters.
	for i := 0; i < 95; i++ {
		ch := rune(32 + i)
		for row := 0; row < 8; row++ {
			// Each row of the character is just the character repeated 5 times.
			sb.WriteString(strings.Repeat(string(ch), 5))
			sb.WriteString("\n")
		}
		if i < 94 {
			sb.WriteString("\n") // blank line separator between characters
		}
	}
	f, err := os.CreateTemp(t.TempDir(), "banner-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(sb.String())
	f.Close()
	return f.Name()
}

// ─────────────────────────────────────────────
// LoadBanner tests
// ─────────────────────────────────────────────

func TestLoadBanner_ValidFile(t *testing.T) {
	path := makeBannerFile(t)
	banner, err := LoadBanner(path)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(banner) == 0 {
		t.Fatal("expected banner map to have entries, got empty map")
	}
}

func TestLoadBanner_FileNotFound(t *testing.T) {
	_, err := LoadBanner("banners/doesnotexist.txt")
	if err == nil {
		t.Fatal("expected an error for missing file, got nil")
	}
}

func TestLoadBanner_SpaceCharacterPresent(t *testing.T) {
	path := makeBannerFile(t)
	banner, err := LoadBanner(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := banner[' ']; !ok {
		t.Error("expected space character to be in banner map")
	}
}

func TestLoadBanner_AllPrintableASCII(t *testing.T) {
	path := makeBannerFile(t)
	banner, err := LoadBanner(path)
	if err != nil {
		t.Fatal(err)
	}
	for ch := rune(32); ch <= 126; ch++ {
		if _, ok := banner[ch]; !ok {
			t.Errorf("expected character %q (code %d) to be in banner map", ch, ch)
		}
	}
}

func TestLoadBanner_EachCharHas8Rows(t *testing.T) {
	path := makeBannerFile(t)
	banner, err := LoadBanner(path)
	if err != nil {
		t.Fatal(err)
	}
	for ch, lines := range banner {
		if len(lines) != 8 {
			t.Errorf("character %q: expected 8 rows, got %d", ch, len(lines))
		}
	}
}

func TestLoadBanner_WindowsLineEndings(t *testing.T) {
	// Write a file with \r\n line endings and verify it loads correctly.
	dir := t.TempDir()
	path := dir + "/banner.txt"
	var sb strings.Builder
	for i := 0; i < 95; i++ {
		ch := rune(32 + i)
		for row := 0; row < 8; row++ {
			sb.WriteString(strings.Repeat(string(ch), 5))
			sb.WriteString("\r\n")
		}
		if i < 94 {
			sb.WriteString("\r\n")
		}
	}
	os.WriteFile(path, []byte(sb.String()), 0644)

	banner, err := LoadBanner(path)
	if err != nil {
		t.Fatalf("expected no error with CRLF file, got: %v", err)
	}
	if len(banner) == 0 {
		t.Fatal("expected non-empty banner map from CRLF file")
	}
}

// ─────────────────────────────────────────────
// RenderLines tests
// ─────────────────────────────────────────────

func TestRenderLines_ReturnsEightRows(t *testing.T) {
	path := makeBannerFile(t)
	banner, _ := LoadBanner(path)
	rows := RenderLines("A", banner)
	if len(rows) != 8 {
		t.Errorf("expected 8 rows, got %d", len(rows))
	}
}

func TestRenderLines_EmptyInput(t *testing.T) {
	path := makeBannerFile(t)
	banner, _ := LoadBanner(path)
	rows := RenderLines("", banner)
	if len(rows) != 8 {
		t.Errorf("expected 8 rows for empty input, got %d", len(rows))
	}
	for i, row := range rows {
		if row != "" {
			t.Errorf("row %d: expected empty string, got %q", i, row)
		}
	}
}

func TestRenderLines_UnknownCharSkipped(t *testing.T) {
	path := makeBannerFile(t)
	banner, _ := LoadBanner(path)
	// emoji is not in the banner map — should be silently skipped, not panic
	rows := RenderLines("😀", banner)
	if len(rows) != 8 {
		t.Errorf("expected 8 rows even with unknown char, got %d", len(rows))
	}
}

func TestRenderLines_MultipleCharacters(t *testing.T) {
	path := makeBannerFile(t)
	banner, _ := LoadBanner(path)
	rowsSingle := RenderLines("A", banner)
	rowsDouble := RenderLines("AB", banner)
	// Each row of "AB" should be longer than each row of "A"
	for i := 0; i < 8; i++ {
		if len(rowsDouble[i]) <= len(rowsSingle[i]) {
			t.Errorf("row %d: expected longer output for two chars than one", i)
		}
	}
}

// ─────────────────────────────────────────────
// GenerateArt tests
// ─────────────────────────────────────────────

func TestGenerateArt_SingleLine(t *testing.T) {
	path := makeBannerFile(t)
	banner, _ := LoadBanner(path)
	result := GenerateArt("Hi", banner)
	if result == "" {
		t.Error("expected non-empty result")
	}
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(lines) != 8 {
		t.Errorf("expected 8 output lines for single-line input, got %d", len(lines))
	}
}

func TestGenerateArt_MultiLine(t *testing.T) {
	path := makeBannerFile(t)
	banner, _ := LoadBanner(path)
	// Two words separated by a newline → 8 rows per word = 16 rows total
	result := GenerateArt("Hi\nHo", banner)
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(lines) != 8 {
		t.Errorf("expected 16 output lines for two-line input, got %d", len(lines))
	}
}

func TestGenerateArt_EmptyInput(t *testing.T) {
	path := makeBannerFile(t)
	banner, _ := LoadBanner(path)
	result := GenerateArt("", banner)
	// Empty input → single empty line written by the loop
	if strings.TrimSpace(result) != "" {
		t.Errorf("expected effectively empty result for empty input, got %q", result)
	}
}

// ─────────────────────────────────────────────
// handleHome tests
// ─────────────────────────────────────────────

func TestHandleHome_GetRootReturns200(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handleHome(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestHandleHome_WrongPathReturns404(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	w := httptest.NewRecorder()
	handleHome(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandleHome_BodyContainsForm(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handleHome(w, req)
	body := w.Body.String()
	if !strings.Contains(body, "<form") {
		t.Error("expected response body to contain a form element")
	}
}

func TestHandleHome_NoPreTagOnHomePage(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handleHome(w, req)
	body := w.Body.String()
	if strings.Contains(body, "<pre>") {
		t.Error("expected no <pre> tag on the home page before any submission")
	}
}

// ─────────────────────────────────────────────
// handleAsciiArt tests
// ─────────────────────────────────────────────

func TestHandleAsciiArt_EmptyTextReturns400(t *testing.T) {
	form := url.Values{}
	form.Set("text", "")
	form.Set("banner", "standard")
	req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handleAsciiArt(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty text, got %d", w.Code)
	}
}

func TestHandleAsciiArt_ValidInputReturns200(t *testing.T) {
	form := url.Values{}
	form.Set("text", "Hi")
	form.Set("banner", "standard")
	req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handleAsciiArt(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for valid input, got %d", w.Code)
	}
}

func TestHandleAsciiArt_ResponseContainsPreTag(t *testing.T) {
	form := url.Values{}
	form.Set("text", "Hi")
	form.Set("banner", "standard")
	req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handleAsciiArt(w, req)
	if !strings.Contains(w.Body.String(), "<pre>") {
		t.Error("expected <pre> tag in response body after art generation")
	}
}

func TestHandleAsciiArt_TextPreservedInResponse(t *testing.T) {
	form := url.Values{}
	form.Set("text", "Hello")
	form.Set("banner", "standard")
	req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handleAsciiArt(w, req)
	// The switch links should embed the original text in their href
	if !strings.Contains(w.Body.String(), "text=Hello") {
		t.Error("expected original text to be preserved in banner switch links")
	}
}

func TestHandleAsciiArt_InvalidBannerReturns500(t *testing.T) {
	form := url.Values{}
	form.Set("text", "Hi")
	form.Set("banner", "nonexistent")
	req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handleAsciiArt(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 for invalid banner, got %d", w.Code)
	}
}

func TestHandleAsciiArt_DefaultsBannerToStandard(t *testing.T) {
	form := url.Values{}
	form.Set("text", "Hi")
	// banner intentionally omitted
	req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handleAsciiArt(w, req)
	// Without a banner field, handler defaults to "standard".
	// If standard.txt exists this returns 200, otherwise 500.
	// Either way it should NOT be 400 (bad request).
	if w.Code == http.StatusBadRequest {
		t.Error("expected handler to default banner, not return 400")
	}
}

func TestHandleAsciiArt_NonASCIIReturns400(t *testing.T) {
	form := url.Values{}
	form.Set("text", "héllo")
	form.Set("banner", "standard")
	req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handleAsciiArt(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for non-ASCII input, got %d", w.Code)
	}
}

// ─────────────────────────────────────────────
// handleSwitch tests
// ─────────────────────────────────────────────

func TestHandleSwitch_EmptyTextReturns400(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ascii-art-switch?text=&banner=standard", nil)
	w := httptest.NewRecorder()
	handleSwitch(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty text, got %d", w.Code)
	}
}

func TestHandleSwitch_ValidQueryReturns200(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ascii-art-switch?text=Hi&banner=standard", nil)
	w := httptest.NewRecorder()
	handleSwitch(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for valid query, got %d", w.Code)
	}
}

func TestHandleSwitch_ResponseContainsPreTag(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ascii-art-switch?text=Hi&banner=standard", nil)
	w := httptest.NewRecorder()
	handleSwitch(w, req)
	if !strings.Contains(w.Body.String(), "<pre>") {
		t.Error("expected <pre> tag in response body")
	}
}

func TestHandleSwitch_InvalidBannerReturns500(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ascii-art-switch?text=Hi&banner=fakebanner", nil)
	w := httptest.NewRecorder()
	handleSwitch(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 for invalid banner, got %d", w.Code)
	}
}

func TestHandleSwitch_DefaultsBannerToStandard(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ascii-art-switch?text=Hi", nil)
	w := httptest.NewRecorder()
	handleSwitch(w, req)
	// Should not be a bad request — handler must default to standard
	if w.Code == http.StatusBadRequest {
		t.Error("expected handler to default banner, not return 400")
	}
}

func TestHandleSwitch_TextPreservedInSwitchLinks(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ascii-art-switch?text=World&banner=standard", nil)
	w := httptest.NewRecorder()
	handleSwitch(w, req)
	if !strings.Contains(w.Body.String(), "text=World") {
		t.Error("expected original text to be preserved in switch links")
	}
}

func TestHandleSwitch_NonASCIIReturns400(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ascii-art-switch?text=héllo&banner=standard", nil)
	w := httptest.NewRecorder()
	handleSwitch(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for non-ASCII input, got %d", w.Code)
	}
}
