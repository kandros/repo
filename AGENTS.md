# AGENTS.md - Development Guidelines

## Build/Test/Lint Commands
- `make build` - Build the binary with version info
- `make run` - Run the application directly
- `go run .` - Alternative run command
- `go build` - Basic build without version info
- `go test ./...` - Run all tests (no tests currently exist)
- `go fmt ./...` - Format code
- `go vet ./...` - Static analysis

## Code Style Guidelines
- **Language**: Go 1.20+
- **Imports**: Standard library first, then third-party, then local packages
- **Naming**: Use camelCase for variables/functions, PascalCase for exported items
- **Error Handling**: Use panic() for unrecoverable errors, return errors for recoverable ones
- **Types**: Define custom types for domain objects (e.g., `type RepoOptions struct`)
- **Comments**: Minimal comments, focus on package-level documentation

## Project Structure
- `main.go` - Entry point with CLI argument handling
- `repos/` - GitHub API integration
- `ui/` - Bubble Tea TUI components
- Uses Charm libraries (bubbletea, bubbles, lipgloss) for terminal UI