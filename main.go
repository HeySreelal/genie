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
	var context string

	// Parse arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Printf("%s v%s\n", appName, version)
			return
		case "--help", "-h":
			printHelp()
			return
		default:
			// If first argument is not a flag, treat it as context
			// Join all arguments as context (in case of spaces)
			context = strings.Join(os.Args[1:], " ")
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
	if context != "" {
		fmt.Printf("📝 Using context: \"%s\"\n", context)
	}

	// Get git diff and determine what type of changes we're analyzing
	diff, changesType, err := getGitDiff()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error getting git diff: %v\n", err)
		os.Exit(1)
	}

	if strings.TrimSpace(diff) == "" {
		fmt.Println("✨ No changes detected. Nothing to commit!")
		return
	}

	// Show what we're analyzing
	switch changesType {
	case "staged":
		fmt.Println("📋 Analyzing staged changes (ready to commit)")
	case "unstaged":
		fmt.Println("📝 No staged changes found, analyzing unstaged changes")
		fmt.Println("💡 Tip: Run 'git add .' to stage all changes before committing")
	}

	// Get git status for context
	status, err := getGitStatus()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error getting git status: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🧠 Generating commit message with Gemini AI...")

	// Generate commit message
	commitMsg, err := generateCommitMessage(apiKey, diff, status, context, changesType)
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
}

func printHelp() {
	fmt.Printf(`%s v%s - AI-powered Git commit message generator

USAGE:
    %s [OPTIONS]
    %s [CONTEXT]

OPTIONS:
    -h, --help      Show this help message
    -v, --version   Show version information

ARGUMENTS:
    CONTEXT         Optional context to help generate better commit messages
                   (e.g., "changes from Bot API 9.0", "refactor for performance")

SETUP:
    1. Get your Gemini API key from: https://aistudio.google.com/apikey
    2. Set the environment variable: export GOOGLE_AI_TOKEN=your_api_key_here
    3. Run %s in any git repository with changes

DESCRIPTION:
    %s analyzes your git changes and generates perfect commit messages
    using Google's Gemini AI. It follows conventional commit standards,
    includes relevant emojis, and automatically copies the message to
    your clipboard for easy use.

    The tool prioritizes staged changes (files added with 'git add'), but
    if no staged changes are found, it will analyze all unstaged changes
    in your working directory.

    You can optionally provide context to help generate more accurate
    commit messages when you have many related changes.

EXAMPLES:
    %s                              # Generate commit message for changes
    %s "Bot API 9.0 migration"      # Generate with context
    %s "performance improvements"   # Generate with context
    %s --version                   # Show version
    %s --help                     # Show this help

`, appName, version, appName, appName, appName, appName, appName, appName, appName, appName, appName)
}

func isGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

func getGitDiff() (string, string, error) {
	// First, try to get staged changes
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", "", err
	}

	stagedDiff := strings.TrimSpace(string(output))

	// If we have staged changes, return them
	if stagedDiff != "" {
		return stagedDiff, "staged", nil
	}

	// If no staged changes, get unstaged changes
	cmd = exec.Command("git", "diff")
	output, err = cmd.Output()
	if err != nil {
		return "", "", err
	}

	unstagedDiff := strings.TrimSpace(string(output))
	if unstagedDiff != "" {
		return unstagedDiff, "unstaged", nil
	}

	// If still no diff, check for untracked files
	cmd = exec.Command("git", "ls-files", "--others", "--exclude-standard")
	output, err = cmd.Output()
	if err != nil {
		return "", "", err
	}

	untrackedFiles := strings.TrimSpace(string(output))
	if untrackedFiles != "" {
		// For untracked files, we can't get a proper diff, so we'll create a summary
		files := strings.Split(untrackedFiles, "\n")
		var summary strings.Builder
		summary.WriteString("New untracked files:\n")
		for _, file := range files {
			if file != "" {
				summary.WriteString(fmt.Sprintf("+ %s\n", file))
			}
		}
		return summary.String(), "untracked", nil
	}

	return "", "", nil
}

func getGitStatus() (string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func generateCommitMessage(apiKey, diff, status, context, changesType string) (string, error) {
	var prompt string

	changesDescription := ""
	switch changesType {
	case "staged":
		changesDescription = "staged changes (ready to commit)"
	case "unstaged":
		changesDescription = "unstaged changes (not yet staged for commit)"
	case "untracked":
		changesDescription = "untracked files (new files not yet added to git)"
	}

	basePrompt := `You are a senior software engineer tasked with writing the perfect git commit message.

Analyze the following git diff and status, then generate a concise, descriptive commit message that follows these guidelines:

1. Start with a relevant emoji that represents the type of change
2. Use conventional commit format after emoji: emoji type(scope): description
3. Types: feat, fix, docs, style, refactor, test, chore, perf, ci, build
4. Keep the first line under 50 characters if possible
5. Be specific about what changed, not just that something changed
6. Use imperative mood ("add" not "added" or "adds")
7. Don't include file names unless crucial to understanding
8. Focus on the "why" and "what" rather than "how"`

	contextSection := ""
	if context != "" {
		contextSection = fmt.Sprintf(`

**IMPORTANT: Take into account this context provided by the developer:**
"%s"

This context should guide your understanding of what these changes are about. Use this information to create a more accurate and meaningful commit message that reflects the broader purpose of these changes.`, context)
	}

	emojiGuidelines := `

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
- 🔥 remove: removing code/files`

	analysisNote := fmt.Sprintf("\n\n**Note:** You are analyzing %s.", changesDescription)

	prompt = basePrompt + contextSection + emojiGuidelines + analysisNote + fmt.Sprintf(`

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
