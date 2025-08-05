package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	appName   = "genie"
	version   = "1.0.0"
	geminiURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent"
)

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
	Error      *ErrorInfo  `json:"error,omitempty"`
}

type Candidate struct {
	Content ContentResponse `json:"content"`
}

type ContentResponse struct {
	Parts []PartResponse `json:"parts"`
}

type PartResponse struct {
	Text string `json:"text"`
}

type ErrorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Printf("%s v%s\n", appName, version)
			return
		case "--help", "-h":
			printHelp()
			return
		}
	}

	// Check if we're in a git repository
	if !isGitRepo() {
		fmt.Fprintf(os.Stderr, "❌ Error: Not a git repository\n")
		os.Exit(1)
	}

	// Get API key from environment
	apiKey := os.Getenv("GOOGLE_AI_TOKEN")
	if apiKey == "" {
		fmt.Fprintf(os.Stderr, "❌ Error: GOOGLE_AI_TOKEN environment variable not set\n")
		fmt.Fprintf(os.Stderr, "   Get your API key from: https://makersuite.google.com/app/apikey\n")
		fmt.Fprintf(os.Stderr, "   Then run: export GOOGLE_AI_TOKEN=your_api_key_here\n")
		os.Exit(1)
	}

	fmt.Println("🔍 Analyzing git changes...")

	// Get git diff
	diff, err := getGitDiff()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error getting git diff: %v\n", err)
		os.Exit(1)
	}

	if strings.TrimSpace(diff) == "" {
		fmt.Println("✨ No changes detected. Nothing to commit!")
		return
	}

	// Get git status for context
	status, err := getGitStatus()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error getting git status: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🧠 Generating commit message with Gemini AI...")

	// Generate commit message
	commitMsg, err := generateCommitMessage(apiKey, diff, status)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error generating commit message: %v\n", err)
		os.Exit(1)
	}

	// Display the generated commit message
	fmt.Println("\n✨ Generated commit message:")
	fmt.Println("┌─────────────────────────────────────────────────────────────────")
	fmt.Printf("│ %s\n", commitMsg)
	fmt.Println("└─────────────────────────────────────────────────────────────────")

	// Copy to clipboard
	fmt.Print("\n📋 Copying to clipboard...")
	err = copyToClipboard(commitMsg)
	if err != nil {
		fmt.Printf(" ❌ Failed\n")
		fmt.Fprintf(os.Stderr, "   Could not copy to clipboard: %v\n", err)
		fmt.Printf("   You can copy manually: %s\n", commitMsg)
	} else {
		fmt.Printf(" ✅ Done!\n")
	}

	fmt.Printf("\n🚀 Ready to commit! Run: git commit -m \"%s\"\n", commitMsg)
	fmt.Println("   Or simply: git commit (message is in clipboard)")

}

func printHelp() {
	fmt.Printf(`%s v%s - AI-powered Git commit message generator

USAGE:
    %s [OPTIONS]

OPTIONS:
    -h, --help      Show this help message
    -v, --version   Show version information

SETUP:
    1. Get your Gemini API key from: https://makersuite.google.com/app/apikey
    2. Set the environment variable: export GEMINI_API_KEY=your_api_key_here
    3. Run %s in any git repository with staged changes

DESCRIPTION:
    %s analyzes your git changes and generates perfect commit messages
    using Google's Gemini AI. It follows conventional commit standards,
    includes relevant emojis, and automatically copies the message to
    your clipboard for easy use.

EXAMPLES:
    %s                    # Generate commit message and copy to clipboard
    %s --version         # Show version
    %s --help           # Show this help

`, appName, version, appName, appName, appName, appName, appName, appName)
}

func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

func getGitDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// If no staged changes, get unstaged changes
	if strings.TrimSpace(string(output)) == "" {
		cmd = exec.Command("git", "diff")
		output, err = cmd.Output()
		if err != nil {
			return "", err
		}
	}

	return string(output), nil
}

func getGitStatus() (string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func generateCommitMessage(apiKey, diff, status string) (string, error) {
	prompt := fmt.Sprintf(`You are a senior software engineer tasked with writing the perfect git commit message.

Analyze the following git diff and status, then generate a concise, descriptive commit message that follows these guidelines:

1. Start with a relevant emoji that represents the type of change
2. Use conventional commit format after emoji: emoji type(scope): description
3. Types: feat, fix, docs, style, refactor, test, chore, perf, ci, build
4. Keep the first line under 50 characters if possible
5. Be specific about what changed, not just that something changed
6. Use imperative mood ("add" not "added" or "adds")
7. Don't include file names unless crucial to understanding
8. Focus on the "why" and "what" rather than "how"

Emoji guidelines:
- ✨ feat: new features
- 🐛 fix: bug fixes
- 📝 docs: documentation
- 💄 style: formatting, styling
- ♻️ refactor: code refactoring
- ✅ test: adding/updating tests
- 🔧 chore: maintenance tasks
- ⚡ perf: performance improvements
- 👷 ci: CI/CD changes
- 📦 build: build system changes
- 🚀 deploy: deployment related
- 🔒 security: security improvements
- 🎨 ui: UI/UX improvements
- 🗃️ database: database changes
- 🔥 remove: removing code/files

Git Status:
%s

Git Diff:
%s

Respond with ONLY the commit message including the emoji, no explanation or additional text.`, status, diff)

	reqBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
				},
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", geminiURL+"?key="+apiKey, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", err
	}

	if geminiResp.Error != nil {
		return "", fmt.Errorf("API error: %s", geminiResp.Error.Message)
	}

	if len(geminiResp.Candidates) == 0 {
		return "", fmt.Errorf("no response from Gemini API")
	}

	if len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini API")
	}

	commitMsg := strings.TrimSpace(geminiResp.Candidates[0].Content.Parts[0].Text)

	// Clean up the response (remove quotes if present)
	commitMsg = strings.Trim(commitMsg, "\"'")

	return commitMsg, nil
}

func copyToClipboard(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("pbcopy")
	case "linux":
		// Try different clipboard utilities
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else if _, err := exec.LookPath("xsel"); err == nil {
			cmd = exec.Command("xsel", "--clipboard", "--input")
		} else if _, err := exec.LookPath("wl-copy"); err == nil {
			cmd = exec.Command("wl-copy") // Wayland
		} else {
			return fmt.Errorf("no clipboard utility found (install xclip, xsel, or wl-copy)")
		}
	case "windows":
		cmd = exec.Command("cmd", "/c", "clip")
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
