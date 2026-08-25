# Go URL Checker / Monitoring Project

## Project Vision

Build a Go-based monitoring tool that can eventually monitor live
applications and infrastructure, detect failures and performance
problems, and notify developers when something goes wrong.

---

# Milestone 1 — CLI URL Checker

Status: COMPLETE

The first version of the project focused on building a reliable command-line
URL checking tool.

## Goals

- [x] Accept URLs from the command line
- [x] Validate that at least one URL is provided
- [x] Validate CLI URL arguments
- [x] Check multiple URLs
- [x] Check URLs concurrently
- [x] Report HTTP status codes
- [x] Report request duration
- [x] Report URL checking errors
- [x] Distinguish request errors from HTTP status failures
- [x] Check all URLs before returning an overall failure
- [x] Make the CLI checking logic testable
- [x] Test error scenarios without depending on real external services
- [x] Maintain a passing test suite

## Current CLI behaviour

A successful check reports:

    URL STATUS DURATION

For example:

    https://example.com 200 120ms

A request-level failure reports:

    URL ERROR: error details

For example:

    http://example.com ERROR: connection refused

An HTTP failure reports the HTTP status:

    URL 500 DURATION

---

# Milestone 2 — Continuous Monitoring

Status: NEXT

Turn the one-shot URL checker into a service that repeatedly checks
configured URLs.

## Goals

- [ ] Run checks continuously
- [ ] Configure check intervals
- [ ] Track whether a URL is UP or DOWN
- [ ] Detect transitions from UP → DOWN
- [ ] Detect transitions from DOWN → UP
- [ ] Track response latency
- [ ] Detect slow responses
- [ ] Record check history
- [ ] Avoid repeatedly alerting for the same outage

## Go concepts we expect to learn

- time.Ticker
- goroutines
- channels
- context
- graceful shutdown
- configuration
- structured logging

---

# Milestone 3 — Alerting

Status: PLANNED

Notify the developer when something important happens.

## Potential alerts

- [ ] Server goes down
- [ ] Server recovers
- [ ] Response becomes unusually slow
- [ ] HTTP 5xx responses
- [ ] Repeated failures
- [ ] SSL certificate nearing expiry

## Potential notification channels

- [ ] Email
- [ ] Slack
- [ ] SMS

---

# Milestone 4 — Monitoring API

Status: PLANNED

Expose monitoring data through an HTTP API.

## Goals

- [ ] Create monitoring endpoints
- [ ] Expose current service status
- [ ] Expose check history
- [ ] Expose latency data
- [ ] Expose incidents
- [ ] Add API tests

---

# Milestone 5 — Monitoring Dashboard

Status: PLANNED

Build a dashboard for viewing monitored systems.

## Potential features

- [ ] Service list
- [ ] UP/DOWN status
- [ ] Response-time graphs
- [ ] Incident history
- [ ] Uptime percentage
- [ ] Recent errors
- [ ] Alert history

---

# Milestone 6 — Application Monitoring

Status: FUTURE

Move beyond simple URL monitoring.

Potential capabilities:

- [ ] API health checks
- [ ] Custom health-check endpoints
- [ ] TCP checks
- [ ] DNS checks
- [ ] SSL certificate monitoring
- [ ] Database connectivity checks
- [ ] Application error monitoring
- [ ] Deployment health checks
- [ ] Resource monitoring

---

# Long-Term Vision

Eventually this project should become a small developer-focused monitoring
platform capable of answering:

1. Is my application running?
2. Is it responding correctly?
3. Is it becoming slower?
4. When did it go down?
5. What caused the failure?
6. When did it recover?
7. Did a deployment cause the problem?
8. Should I be alerted?

The project should remain a learning project as well as a usable tool.

Each major feature should introduce and reinforce a real Go concept.