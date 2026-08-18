
### 3. `docs/requirements.md`

```markdown
# Requirements

## Goal

Build a Go-based URL checker that can determine the availability and basic HTTP characteristics of URLs.

The system will be developed incrementally as new requirements are identified.

## Current Requirements

### URL Checking

The system must be able to check a URL.

### HTTP Status

The system must report the HTTP status code when an HTTP response is received.

HTTP error status codes such as `404` and `500` are valid responses and should not automatically be treated as request errors.

### Request Errors

The system must report an error when the HTTP request itself fails.

Examples include:

- DNS failures
- Connection failures
- Request timeouts

### Duration

The system must measure how long the URL check takes.

### Timeout

The system must support a configurable HTTP request timeout.

The current default timeout is:

```text
7 seconds