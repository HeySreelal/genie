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

// Color constants
const (
	colorReset  = "\033[0m"
	colorGray   = "\033[90m" // Dim gray for logs
	colorGreen  = "\033[32m" // Success messages
	colorBlue   = "\033[34m" // Info messages
	colorYellow = "\033[33m" // Warnings
	colorRed    = "\033[31m" // Errors
	colorBold   = "\033[1m"  // Bold text
	colorCyan   = "\033[36m" // Highlights
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

// Helper functions for colored output
func grayf(format string, args ...interface{}) {
	fmt.Printf(colorGray+format+colorReset, args...)
}

func redf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, colorRed+format+colorReset, args...)
}

func boldf(format string, args ...interface{}) {
	fmt.Printf(colorBold+format+colorReset, args...)
}

func cyanf(format string, args ...interface{}) {
	fmt.Printf(colorCyan+format+colorReset, args...)
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
		redf("❌ Error: Not a git repository\n")
		os.Exit(1)
	}

	// Get API key from environment
	apiKey := os.Getenv("GOOGLE_AI_TOKEN")
	if apiKey == "" {
		redf("❌ Error: GOOGLE_AI_TOKEN environment variable not set\n")
		redf("   Get your API key from: https://aistudio.google.com/apikey\n")
		redf("   Then run: export GOOGLE_AI_TOKEN=your_api_key_here\n")
		os.Exit(1)
	}

	if context != "" {
		grayf("📝 Context: \"%s\"\n", context)
	}

	// Get git diff and determine what type of changes we're analyzing
	diff, changesType, err := getGitDiff()
	if err != nil {
		redf("❌ Error getting git diff: %v\n", err)
		os.Exit(1)
	}

	if strings.TrimSpace(diff) == "" {
		fmt.Println("✨ No changes detected. Nothing to commit!")
		return
	}

	// Show what we're analyzing - more subtle
	switch changesType {
	case "staged":
		grayf("Analyzing staged changes...\n")
	case "unstaged":
		grayf("No staged changes, analyzing unstaged changes...\n")
		grayf("💡 Tip: Run 'git add .' to stage changes first\n")
	case "untracked":
		grayf("Analyzing untracked files...\n")
		grayf("💡 Tip: Run 'git add .' to stage files first\n")
	}

	// Get git status for context
	status, err := getGitStatus()
	if err != nil {
		redf("❌ Error getting git status: %v\n", err)
		os.Exit(1)
	}

	// Generate commit message
	commitMsg, err := generateCommitMessage(apiKey, diff, status, context, changesType)
	if err != nil {
		redf("❌ Error generating commit message: %v\n", err)
		os.Exit(1)
	}

	// Display the generated commit message - make this prominent
	fmt.Println()
	boldf("✨ Generated commit message:\n")
	cyanf("┌─────────────────────────────────────────────────────────────────\n")
	cyanf("│ %s\n", commitMsg)
	cyanf("└─────────────────────────────────────────────────────────────────\n")

	// Copy to clipboard
	err = copyToClipboard(commitMsg)
	if err != nil {
		grayf("📋 Could not copy to clipboard: %v\n", err)
	} else {
		grayf("📋 Copied to clipboard\n")
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
	changesDescription := ""
	switch changesType {
	case "staged":
		changesDescription = "staged changes (ready to commit)"
	case "unstaged":
		changesDescription = "unstaged changes (not yet staged for commit)"
	case "untracked":
		changesDescription = "untracked files (new files not yet added to git)"
	}

	// Build the enhanced prompt
	var promptBuilder strings.Builder

	promptBuilder.WriteString(`You are a world-class senior software engineer and git expert with years of experience writing perfect, professional commit messages. Your task is to analyze git changes and generate the ideal commit message.

ANALYSIS REQUIREMENTS:
- Carefully examine the git diff and status to understand what actually changed
- Identify the primary purpose and impact of the changes
- Consider the scope and complexity of modifications
- Determine if this is a feature, fix, refactor, or other type of change

COMMIT MESSAGE RULES:
1. 🎯 START WITH APPROPRIATE EMOJI: Choose the most relevant emoji that represents the change type
2. 📏 FORMAT: Use conventional commit format: "emoji type(scope): description"
3. 🔤 IMPERATIVE MOOD: Use imperative mood ("add" not "added", "fix" not "fixed")
4. 📐 LENGTH: Keep first line under 50 characters when possible, maximum 72
5. 🎯 BE SPECIFIC: Focus on WHAT changed and WHY, not HOW
6. 🚫 NO FILENAMES: Don't mention specific files unless absolutely crucial
7. 💡 CLARITY: Make it immediately clear what the commit accomplishes
8. 🏷️ SCOPE: Include scope in parentheses when it adds clarity (e.g., auth, api, ui)

CONVENTIONAL COMMIT TYPES:
- feat: New features or enhancements
- fix: Bug fixes and corrections  
- docs: Documentation changes
- style: Code formatting, whitespace, styling
- refactor: Code restructuring without functionality changes
- test: Adding or modifying tests
- chore: Maintenance, build process, dependencies
- perf: Performance improvements
- ci: CI/CD pipeline changes
- build: Build system, external dependencies
- revert: Reverting previous changes

EMOJI SELECTION GUIDE:
✨ feat: new features, enhancements
🐛 fix: bug fixes, error corrections
📝 docs: documentation, README updates
💄 style: formatting, code style, UI styling
♻️ refactor: code refactoring, restructuring
✅ test: adding/updating tests
🔧 chore: maintenance, config, build
⚡ perf: performance optimizations
👷 ci: CI/CD, workflows, automation
📦 build: build system, dependencies
🚀 deploy: deployment, releases
🔒 security: security fixes, improvements
🎨 ui: UI/UX improvements, design
🗃️ database: database changes, migrations
🔥 remove: removing code, files, features
🩹 hotfix: critical fixes
🚚 move: moving or renaming files
📱 responsive: mobile/responsive changes
🌐 i18n: internationalization, localization
🔊 logging: adding or updating logs
🔇 mute: removing logs
👥 contributor: adding contributors
🚸 accessibility: improving accessibility
💚 green: fixing CI, improving build
🔖 release: version tags, releases
🚨 warning: fixing warnings, linter issues
🚧 wip: work in progress
💥 breaking: breaking changes
📈 analytics: adding analytics, tracking
🔐 auth: authentication, authorization
🌍 global: global changes, configurations`)

	// Add context section if provided
	if context != "" {
		promptBuilder.WriteString(fmt.Sprintf(`

🎯 DEVELOPER CONTEXT:
The developer provided this context: "%s"

This context is CRITICAL - use it to understand the broader purpose and ensure your commit message accurately reflects the intended changes within this context. The context should guide your interpretation of what these technical changes accomplish at a higher level.`, context))
	}

	promptBuilder.WriteString(fmt.Sprintf(`

📊 CHANGE ANALYSIS:
You are analyzing: %s

Git Status Output:
%s

Git Diff/Changes:
%s

🎯 RESPONSE FORMAT:
Respond with ONLY the commit message including emoji. No explanations, quotes, or additional text.

EXAMPLES OF EXCELLENT COMMIT MESSAGES:
✨ feat(auth): add OAuth2 Google integration
🐛 fix(api): handle null response in user endpoint  
♻️ refactor(utils): simplify date formatting logic
📝 docs: update API authentication guide
🔧 chore(deps): update React to v18.2.0
⚡ perf(db): optimize user query with indexing
🎨 ui: improve button hover animations
🔒 security: sanitize user input in forms

Generate the perfect commit message now:`, changesDescription, status, diff))

	reqBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: promptBuilder.String()},
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
