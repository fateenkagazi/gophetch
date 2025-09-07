# Contributing to Gophetch

Thank you for your interest in contributing to Gophetch. This document provides guidelines for contributing to the project.

## Getting Started

1. Fork the repository
2. Clone your fork locally
3. Create a new branch for your changes
4. Make your changes
5. Test your changes thoroughly
6. Submit a pull request

## Development Setup

```bash
git clone https://github.com/your-username/gophetch.git
cd gophetch
go mod tidy
go build
```

## Code Style

- Follow standard Go formatting with `gofmt`
- Use meaningful variable and function names
- Add comments for public functions and complex logic
- Keep functions focused and reasonably sized

## Testing

- Test your changes on multiple platforms (Windows, Linux, macOS, Android/Termux)
- Verify that the application starts and displays correctly
- Test custom frame file loading if your changes affect that functionality
- Ensure graceful shutdown with Ctrl+C

## Pull Request Guidelines

- Provide a clear description of your changes
- Reference any related issues
- Keep changes focused and atomic
- Update documentation if necessary
- Test on at least two different platforms

## Areas for Contribution

- Cross-platform compatibility improvements
- Performance optimizations
- Additional system information metrics
- Enhanced ASCII animation features
- Documentation improvements
- Bug fixes

## Reporting Issues

When reporting issues, please include:
- Operating system and version
- Go version
- Steps to reproduce
- Expected vs actual behavior
- Any error messages or logs

## Questions

If you have questions about contributing, please open an issue for discussion.
