# Security Policy

## Supported Versions

We provide security updates for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| Latest  | :white_check_mark: |
| < Latest| :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability in Gophetch, please report it responsibly:

1. **Do not** open a public issue
2. Email security details to: [security@example.com]
3. Include the following information:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

## Response Timeline

- We will acknowledge receipt within 48 hours
- We will provide a detailed response within 7 days
- We will keep you informed of our progress

## Security Considerations

Gophetch is a terminal application that:
- Reads system information (CPU, memory, disk usage)
- Accesses environment variables for username detection
- Creates temporary files for permission testing
- Executes system commands (whoami, id, ps)

## Best Practices

- Run Gophetch in trusted environments
- Review any custom frame files before loading
- Keep your Go installation updated
- Use the latest version of Gophetch

## Scope

This security policy covers:
- The Gophetch application itself
- Official releases and documentation
- The main repository

Third-party integrations or custom modifications are not covered by this policy.
