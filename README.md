# 🧞‍♂️ genie

> *Your magical AI-powered Git commit message generator*

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey?style=for-the-badge)]()
[![AI Powered](https://img.shields.io/badge/AI-Gemini%20Powered-FF6F00?style=for-the-badge&logo=google)](https://ai.google.dev/)

<div align="center">

**Never write boring commit messages again!** ✨

`genie` analyzes your git changes and generates perfect, emoji-rich commit messages using Google's Gemini AI with its generous free tier.

[🚀 Quick Start](#-quick-start) • [📖 Features](#-features) • [⚙️ Installation](#%EF%B8%8F-installation) • [🎯 Usage](#-usage) • [🤝 Contributing](#-contributing)

---

</div>

### 🎬 Demo

```bash
$ genie

🔍 Analyzing git changes...
🧠 Generating commit message with Gemini AI...

✨ Generated commit message:
┌─────────────────────────────────────────────────────────────────
│ 📦 build: Add installation script
└─────────────────────────────────────────────────────────────────

📋 Copying to clipboard... ✅ Done!
```

## 🚀 Quick Start

```bash
# 1. Clone the repository
git clone https://github.com/heysreelal/genie.git
cd genie

# 2. Run the installer (builds and installs globally)
chmod +x install.sh
./install.sh

# 3. Get your Gemini API key
# Visit: https://aistudio.google.com/apikey

# 4. Set environment variable
export GOOGLE_AI_TOKEN=your_api_key_here

# 5. Use it anywhere!
genie
```

## ✨ Features

<table>
<tr>
<td width="50%">

### 🎨 **Smart Emoji Integration**
Automatically adds contextual emojis to your commit messages for better visual recognition.

### 🧠 **AI-Powered Analysis** 
Uses Google's Gemini AI to understand your changes and generate meaningful commit messages.

### 📋 **Clipboard Ready**
Instantly copies generated messages to your clipboard for seamless workflow integration.

### 💬 **Context Aware**
Provide optional context to generate more accurate commit messages for complex changes.

</td>
<td width="50%">

### 📏 **Conventional Commits**
Follows industry-standard conventional commit format with proper scoping and typing.

### 🌍 **Cross-Platform**
Works perfectly on macOS, Linux, and Windows with native clipboard support.

### ⚡ **Lightning Fast**
Quick analysis and generation - no more staring at blank commit message boxes!

### 💸 **Free Tier Friendly**
Built specifically to leverage Gemini's generous free tier limits.

</td>
</tr>
</table>

## 🎯 Commit Message Examples

| Type | Example |
|------|---------|
| **New Feature** | `✨ feat(auth): add OAuth2 login integration` |
| **Bug Fix** | `🐛 fix(api): resolve null pointer exception in user service` |
| **Documentation** | `📝 docs(readme): update installation instructions` |
| **Performance** | `⚡ perf(db): optimize user query with indexing` |
| **Refactor** | `♻️ refactor(components): extract reusable button component` |
| **Tests** | `✅ test(auth): add unit tests for login validation` |

## ⚙️ Installation

### Option 1: Automated Installation (Recommended)

The easiest way to install `genie` globally:

```bash
# Clone the repository
git clone https://github.com/heysreelal/genie.git
cd genie

# Run the installer
chmod +x install.sh
./install.sh
```

The installer will:
- ✅ Check all requirements (Go, Git)
- 🔨 Build the binary
- 📦 Install to `/usr/local/bin` (requires sudo)
- ✨ Verify installation
- 📋 Show setup instructions

### Option 2: Manual Installation

```bash
# Clone the repository
git clone https://github.com/heysreelal/genie.git
cd genie

# Build the binary
go build -o genie main.go

# Install globally (optional)
sudo mv genie /usr/local/bin/
```

### Option 3: Direct Download

Download the latest release from [GitHub Releases](https://github.com/heysreelal/genie/releases)

### Prerequisites

<details>
<summary><strong>📋 Clipboard Support (Linux users)</strong></summary>

Install one of these clipboard utilities:

```bash
# Ubuntu/Debian
sudo apt install xclip

# Arch Linux  
sudo pacman -S xclip

# Fedora
sudo dnf install xclip

# Or alternatives
sudo apt install xsel      # Alternative 1
sudo apt install wl-copy  # For Wayland users
```

</details>

<details>
<summary><strong>🔑 Gemini API Key Setup</strong></summary>

1. Visit [Google AI Studio](https://aistudio.google.com/apikey)
2. Create a new API key
3. Set the environment variable:

```bash
# Temporary (current session)
export GOOGLE_AI_TOKEN=your_api_key_here

# Permanent (add to ~/.bashrc or ~/.zshrc)
echo 'export GOOGLE_AI_TOKEN=your_api_key_here' >> ~/.bashrc
source ~/.bashrc
```

</details>

## 🎯 Usage

### Basic Usage

```bash
# Generate commit message for your changes
genie

# Generate with context for better accuracy
genie "refactoring user authentication"
genie "Bot API 9.0 migration"
genie "performance improvements"

# Show help
genie --help

# Show version  
genie --version
```

### Advanced Usage with Context

When you have complex changes spanning multiple areas, provide context to help the AI generate more accurate commit messages:

```bash
# For API updates
genie "upgrading to Bot API 9.0"

# For refactoring work
genie "extracting reusable components"

# For bug fixes
genie "fixing memory leak in background service"

# For feature work
genie "implementing dark mode theme"
```

### Workflow Integration

```bash
# 1. Make your changes
echo "console.log('Hello World');" > app.js

# 2. Stage your changes
git add .

# 3. Generate commit message (with or without context)
genie "adding welcome message"

# 4. Commit (message already in clipboard!)
git commit
# Or paste the message: git commit -m "✨ feat(app): add welcome message display"
```

### Smart Change Detection

`genie` intelligently handles different types of changes:

- **Staged changes** (preferred): Analyzes files ready for commit
- **Unstaged changes**: Analyzes modified files when nothing is staged
- **Untracked files**: Handles new files not yet added to git

```bash
# Example workflow
$ git status
On branch main
Changes not staged for commit:
  modified:   src/auth.js
  
Untracked files:
  src/newfeature.js

$ genie "adding OAuth integration"
📝 No staged changes found, analyzing unstaged changes
💡 Tip: Run 'git add .' to stage all changes before committing

✨ Generated commit message:
┌─────────────────────────────────────────────────────────────────
│ ✨ feat(auth): add OAuth integration with Google provider
└─────────────────────────────────────────────────────────────────
```

### Pro Tips 💡

- **Use context for complex changes**: `genie "migrating from REST to GraphQL"`
- **Stage specific files**: `git add file1.js file2.css` before running `genie`
- **Review changes first**: `git diff --cached` to see what will be analyzed
- **Combine with hooks**: Integrate `genie` into your git hooks for automated workflows

## 🎨 Emoji Guide

| Emoji | Type | Description |
|-------|------|-------------|
| ✨ | `feat` | New features |
| 🐛 | `fix` | Bug fixes |
| 📝 | `docs` | Documentation |
| 💄 | `style` | Code formatting |
| ♻️ | `refactor` | Code refactoring |
| ⚡ | `perf` | Performance improvements |
| ✅ | `test` | Tests |
| 🔧 | `chore` | Maintenance |
| 👷 | `ci` | CI/CD changes |
| 📦 | `build` | Build system |
| 🚀 | `deploy` | Deployment |
| 🔒 | `security` | Security fixes |
| 🎨 | `ui` | UI/UX improvements |
| 🗃️ | `database` | Database changes |
| 🔥 | `remove` | Code removal |

## 🛠️ Configuration

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `GOOGLE_AI_TOKEN` | ✅ | Your Google Gemini API key |

### Command Line Options

```bash
genie [CONTEXT]
genie [OPTIONS]

CONTEXT:
  Optional context to help generate better commit messages
  Examples: "Bot API 9.0 migration", "performance improvements"

OPTIONS:
  -h, --help     Show help information
  -v, --version  Show version information
```

## 🤝 Contributing

We love contributions! Here's how you can help make `genie` even better:

### 🐛 Found a Bug?

1. Check existing [issues](https://github.com/heysreelal/genie/issues)
2. Create a new issue with detailed information
3. Include your OS, Go version, and error messages

### 💡 Have an Idea?

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Add tests if applicable
5. Commit using `genie` (dogfooding! 🐕)
6. Push and create a Pull Request

### 🔧 Development Setup

```bash
# Clone your fork
git clone https://github.com/heysreelal/genie.git
cd genie

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Build and test locally
go build -o genie main.go
./genie --help
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Google Gemini AI** - For powering our intelligent commit message generation with generous free tier limits
- **Go Community** - For the amazing ecosystem and tools
- **Contributors** - Thank you to everyone who helps improve `genie`!

## 📈 Stats

<div align="center">

![GitHub stars](https://img.shields.io/github/stars/heysreelal/genie?style=social)
![GitHub forks](https://img.shields.io/github/forks/heysreelal/genie?style=social)
![GitHub issues](https://img.shields.io/github/issues/heysreelal/genie)
![GitHub pull requests](https://img.shields.io/github/issues-pr/heysreelal/genie)

**Made with ❤️ by developers, for developers**

[⭐ Star us on GitHub](https://github.com/heysreelal/genie) • [🐛 Report Issues](https://github.com/heysreelal/genie/issues) • [💬 Discussions](https://github.com/heysreelal/genie/discussions)

</div>

---

<div align="center">

### 🎭 Inspired By

This project was inspired by [Commity](https://github.com/muhd-ameen/Commity), a similar tool that uses OpenAI's API. We created `genie` to take advantage of Google Gemini's generous free tier, making AI-powered commit messages accessible to everyone! 🚀

<sub>🧞‍♂️ <strong>genie</strong> - Because great commit messages shouldn't require three wishes!</sub>

</div>