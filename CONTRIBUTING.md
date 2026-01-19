# Contributing to amazon-cli

Thank you for your interest in contributing to amazon-cli! This document provides guidelines and instructions for contributing to the project.

## Development Setup

### Prerequisites

- Go 1.22 or later
- Git

### Install Go

If you don't have Go installed, download it from [golang.org](https://golang.org/dl/).

Verify your installation:

```bash
go version
```

### Clone the Repository

```bash
git clone https://github.com/zkwentz/amazon-cli.git
cd amazon-cli
```

### Install Dependencies

```bash
go mod download
```

### Install Development Tools

Install golangci-lint for code quality checks:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Build from Source

```bash
go build -o amazon-cli .
```

Or use the Makefile:

```bash
make build
```

You can also install directly from the repository:

```bash
go install github.com/zkwentz/amazon-cli@latest
```

## Running Tests

### Run All Tests

```bash
go test ./...
```

Or use the Makefile:

```bash
make test
```

### Run Tests with Verbose Output

```bash
go test -v ./...
```

### Run Tests with Race Detection

```bash
go test -v -race ./...
```

### Run Tests with Coverage

```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

Or use the Makefile:

```bash
make test-coverage
```

### Run Tests for a Specific Package

```bash
go test ./internal/amazon/...
go test ./pkg/models/...
```

### Run a Specific Test

```bash
go test -v -run TestFunctionName ./path/to/package
```

## Building

### Build for Current Platform

```bash
go build -o amazon-cli .
```

Or use the Makefile:

```bash
make build
```

### Build for All Platforms

```bash
make build-all
```

This creates binaries in the `dist/` directory for:
- macOS (amd64, arm64)
- Linux (amd64, arm64)
- Windows (amd64)

### Install Locally

Install to your `$GOPATH/bin`:

```bash
make install
```

### Clean Build Artifacts

```bash
make clean
```

## Code Style

### Formatting

All Go code must be formatted with `gofmt`. Format your code before committing:

```bash
go fmt ./...
```

Or use the Makefile:

```bash
make fmt
```

### Linting

We use `golangci-lint` to enforce code quality standards. Run the linter before submitting PRs:

```bash
golangci-lint run
```

Or use the Makefile:

```bash
make lint
```

Fix any issues reported by the linter. Common checks include:
- Code formatting
- Unused variables and imports
- Error handling
- Code complexity
- Security issues
- Best practices

### Code Organization

- Follow standard Go project layout
- Keep packages focused and cohesive
- Use meaningful variable and function names
- Add comments for exported functions and types
- Write clear error messages

### Testing Standards

- Write tests for new functionality
- Maintain or improve code coverage
- Use table-driven tests where appropriate
- Mock external dependencies
- Test error paths and edge cases

## Submitting Pull Requests

### Branch Naming

Use descriptive branch names that indicate the purpose of your changes:

- `feature/description` - New features
- `fix/description` - Bug fixes
- `docs/description` - Documentation updates
- `refactor/description` - Code refactoring
- `test/description` - Test additions or updates

Examples:
- `feature/add-wish-list-support`
- `fix/cart-quantity-validation`
- `docs/update-api-examples`
- `refactor/simplify-auth-flow`

### Commit Messages

Write clear, descriptive commit messages:

**Format:**
```
Short summary (50 chars or less)

More detailed explanation if needed. Wrap at 72 characters.
Explain what changed and why, not how.

- Bullet points are okay
- Use present tense ("Add feature" not "Added feature")
- Reference issues: "Fixes #123" or "Relates to #456"
```

**Examples:**

Good:
```
Add support for wish list management

Implement commands to view, add, and remove items from wish lists.
Includes JSON output for all operations and proper error handling.

Fixes #42
```

Bad:
```
Update code
```

### Pull Request Process

1. **Create a feature branch** from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**:
   - Write clear, maintainable code
   - Add tests for new functionality
   - Update documentation as needed

3. **Ensure quality checks pass**:
   ```bash
   make fmt
   make lint
   make test
   ```

4. **Commit your changes**:
   ```bash
   git add .
   git commit -m "Your descriptive commit message"
   ```

5. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Open a pull request**:
   - Go to the [repository](https://github.com/zkwentz/amazon-cli) on GitHub
   - Click "New Pull Request"
   - Select your branch
   - Fill out the PR template (see below)

### Pull Request Template

When opening a PR, include:

```markdown
## Description
Brief description of what this PR does.

## Type of Change
- [ ] Bug fix (non-breaking change that fixes an issue)
- [ ] New feature (non-breaking change that adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Refactoring (no functional changes)
- [ ] Test updates

## Testing
Describe how you tested your changes:
- [ ] All existing tests pass
- [ ] Added new tests for new functionality
- [ ] Manually tested the changes

## Checklist
- [ ] Code follows the project's style guidelines
- [ ] Code has been formatted with `go fmt`
- [ ] Code passes `golangci-lint` checks
- [ ] All tests pass (`go test ./...`)
- [ ] Documentation has been updated (if applicable)
- [ ] Commit messages are clear and descriptive

## Related Issues
Fixes #(issue number)
```

### Review Process

- All PRs require review before merging
- Address reviewer feedback promptly
- Keep PRs focused and reasonably sized
- Update your branch if `main` has changed
- Ensure CI checks pass

### Continuous Integration

Our CI pipeline runs automatically on all PRs:

- **Tests**: All tests must pass with race detection
- **Lint**: Code must pass golangci-lint checks
- **Build**: Code must build successfully for all target platforms

Check the Actions tab on GitHub to see CI results.

## Additional Guidelines

### Error Handling

- Always handle errors appropriately
- Provide meaningful error messages
- Use structured error types where applicable
- Return JSON-formatted errors for CLI commands

### Security

- Never commit credentials or sensitive data
- Review code for security vulnerabilities
- Use HTTPS for all external connections
- Validate user input appropriately

### Documentation

- Update README.md for new features
- Add inline comments for complex logic
- Keep documentation accurate and up-to-date
- Include examples in documentation

## Getting Help

- Open an issue for bugs or feature requests
- Check existing issues before creating new ones
- Ask questions in issue discussions
- Be respectful and professional

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help create a welcoming environment
- Follow the project's guidelines

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
