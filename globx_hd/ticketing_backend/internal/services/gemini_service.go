package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// ============================================================================
// GEMINI API SERVICE
// ============================================================================
// This service handles communication with Google's Gemini 2.5 Flash API
// for email parsing and entity extraction.
// ============================================================================

const (
	GeminiAPIURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent"
)

// GeminiService handles Gemini API interactions
type GeminiService struct {
	apiKey     string
	httpClient *http.Client
}

// GeminiRequest represents the request body for Gemini API
type GeminiRequest struct {
	Contents         []GeminiContent `json:"contents"`
	GenerationConfig GeminiGenConfig `json:"generationConfig,omitempty"`
}

// GeminiContent represents a content block
type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart represents a part of the content
type GeminiPart struct {
	Text string `json:"text"`
}

// GeminiGenConfig represents generation configuration
type GeminiGenConfig struct {
	Temperature     float64 `json:"temperature,omitempty"`
	MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
}

// GeminiResponse represents the response from Gemini API
type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
	Error      *GeminiError      `json:"error,omitempty"`
}

// GeminiCandidate represents a response candidate
type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

// GeminiError represents an API error
type GeminiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// ExtractedEmailData represents the structured data extracted by Gemini
type ExtractedEmailData struct {
	PhoneNumbers  []string `json:"phone_numbers"`
	PersonNames   []string `json:"person_names"`
	OrgNames      []string `json:"org_names"`
	ProductHints  []string `json:"product_hints"`
	PriorityHints []string `json:"priority_hints"`
}

// ParsedEmail represents the email components parsed from the prompt
type ParsedEmail struct {
	SenderEmail string
	SenderName  string
	Subject     string
	Body        string
}

// NewGeminiService creates a new Gemini service instance
func NewGeminiService() *GeminiService {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		// Fallback for development
		apiKey = "YOUR_GEMINI_API_KEY"
	}

	return &GeminiService{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewGeminiServiceWithKey creates a new Gemini service with a specific API key
func NewGeminiServiceWithKey(apiKey string) *GeminiService {
	return &GeminiService{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CallGemini sends a prompt to Gemini API and returns the response text
func (s *GeminiService) CallGemini(prompt string) (string, error) {
	if s.apiKey == "" || s.apiKey == "YOUR_GEMINI_API_KEY" {
		return "", fmt.Errorf("GEMINI_API_KEY not configured")
	}

	// Build request
	reqBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
		GenerationConfig: GeminiGenConfig{
			Temperature:     0.1, // Low temperature for consistent extraction
			MaxOutputTokens: 2048,
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s?key=%s", GeminiAPIURL, s.apiKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Gemini API: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API error
	if geminiResp.Error != nil {
		return "", fmt.Errorf("gemini API error: %s (code: %d, status: %s)",
			geminiResp.Error.Message,
			geminiResp.Error.Code,
			geminiResp.Error.Status)
	}

	// Extract text from response
	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini API")
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}

// ExtractEmailData sends the prompt to Gemini and parses the structured response
func (s *GeminiService) ExtractEmailData(prompt string) (*ExtractedEmailData, error) {
	// Call Gemini
	responseText, err := s.CallGemini(prompt)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response from Gemini
	extracted, err := s.parseExtractedData(responseText)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Gemini response: %w", err)
	}

	return extracted, nil
}

// parseExtractedData parses Gemini's JSON response into structured data
func (s *GeminiService) parseExtractedData(responseText string) (*ExtractedEmailData, error) {
	// Clean up the response - remove markdown code blocks if present
	cleanedText := strings.TrimSpace(responseText)

	// Remove markdown code block markers (handles ```json, ```JSON, ``` with newlines)
	// First, handle opening ``` markers with optional language identifier
	if strings.HasPrefix(cleanedText, "```") {
		// Find the end of the first line (after ```json or ```)
		firstNewline := strings.Index(cleanedText, "\n")
		if firstNewline != -1 {
			cleanedText = cleanedText[firstNewline+1:]
		} else {
			// No newline, just remove the ``` prefix
			cleanedText = strings.TrimPrefix(cleanedText, "```json")
			cleanedText = strings.TrimPrefix(cleanedText, "```JSON")
			cleanedText = strings.TrimPrefix(cleanedText, "```")
		}
	}

	// Remove closing ``` marker (handles trailing whitespace/newlines)
	cleanedText = strings.TrimSpace(cleanedText)
	if strings.HasSuffix(cleanedText, "```") {
		cleanedText = strings.TrimSuffix(cleanedText, "```")
		cleanedText = strings.TrimSpace(cleanedText)
	}

	// Try to extract JSON object - find first { and last }
	jsonStart := strings.Index(cleanedText, "{")
	jsonEnd := strings.LastIndex(cleanedText, "}")
	if jsonStart != -1 && jsonEnd != -1 && jsonEnd > jsonStart {
		cleanedText = cleanedText[jsonStart : jsonEnd+1]
	}

	// Parse JSON
	var extracted ExtractedEmailData
	if err := json.Unmarshal([]byte(cleanedText), &extracted); err != nil {
		// Return empty data with error details
		return &ExtractedEmailData{
			PhoneNumbers:  []string{},
			PersonNames:   []string{},
			OrgNames:      []string{},
			ProductHints:  []string{},
			PriorityHints: []string{},
		}, fmt.Errorf("JSON parse error: %w, response was: %s", err, responseText[:minInt(200, len(responseText))])
	}

	// Ensure no nil slices
	if extracted.PhoneNumbers == nil {
		extracted.PhoneNumbers = []string{}
	}
	if extracted.PersonNames == nil {
		extracted.PersonNames = []string{}
	}
	if extracted.OrgNames == nil {
		extracted.OrgNames = []string{}
	}
	if extracted.ProductHints == nil {
		extracted.ProductHints = []string{}
	}
	if extracted.PriorityHints == nil {
		extracted.PriorityHints = []string{}
	}

	return &extracted, nil
}

// ParseEmailFromPrompt extracts email components (From, Subject, Body) from the prompt
func (s *GeminiService) ParseEmailFromPrompt(prompt string) *ParsedEmail {
	parsed := &ParsedEmail{}

	// Extract sender email using regex
	// Pattern: From: "Name" <email@domain.com> or From: email@domain.com
	fromPattern := regexp.MustCompile(`From:\s*(?:"([^"]*)")?\s*<?([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})>?`)
	if matches := fromPattern.FindStringSubmatch(prompt); len(matches) >= 3 {
		parsed.SenderName = strings.TrimSpace(matches[1])
		parsed.SenderEmail = strings.TrimSpace(matches[2])
	}

	// If no name found, try alternate patterns
	if parsed.SenderName == "" && parsed.SenderEmail != "" {
		// Try to extract name from email username
		parts := strings.Split(parsed.SenderEmail, "@")
		if len(parts) > 0 {
			// Convert username to potential name (e.g., john.doe -> John Doe)
			username := parts[0]
			username = strings.ReplaceAll(username, ".", " ")
			username = strings.ReplaceAll(username, "_", " ")
			username = strings.ReplaceAll(username, "-", " ")
			// Capitalize words
			words := strings.Fields(username)
			for i, w := range words {
				if len(w) > 0 {
					words[i] = strings.ToUpper(string(w[0])) + strings.ToLower(w[1:])
				}
			}
			parsed.SenderName = strings.Join(words, " ")
		}
	}

	// Extract subject
	subjectPattern := regexp.MustCompile(`Subject:\s*(.+?)(?:\s+Body:|$)`)
	if matches := subjectPattern.FindStringSubmatch(prompt); len(matches) >= 2 {
		parsed.Subject = strings.TrimSpace(matches[1])
	}

	// Extract body - everything after "Body:" until the extraction instructions
	bodyPattern := regexp.MustCompile(`Body:\s*([\s\S]+?)(?:\s+Extract and return|$)`)
	if matches := bodyPattern.FindStringSubmatch(prompt); len(matches) >= 2 {
		parsed.Body = strings.TrimSpace(matches[1])
	}

	// If no body found with the pattern, try to get everything after "Body:"
	if parsed.Body == "" {
		bodyIdx := strings.Index(prompt, "Body:")
		if bodyIdx != -1 {
			bodyContent := prompt[bodyIdx+5:]
			// Find where the extraction instructions start
			extractIdx := strings.Index(bodyContent, "Extract and return")
			if extractIdx != -1 {
				bodyContent = bodyContent[:extractIdx]
			}
			parsed.Body = strings.TrimSpace(bodyContent)
		}
	}

	return parsed
}

// Helper function for string truncation
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
