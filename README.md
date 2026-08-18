# Go URL Checker

A URL checker built in Go as an engineering exercise.

The project is being developed incrementally using:

- Test-driven development
- Small, focused changes
- Automated tests
- Git commits at meaningful checkpoints
- Explicit architecture and requirements documentation

## Current Features

The checker can:

- Check whether an HTTP request to a URL succeeds
- Return the HTTP status code
- Return request errors
- Measure request duration
- Support configurable request timeouts
- Use a default request timeout of 7 seconds

## Running Tests

Run the complete test suite:

```bash
go test ./...