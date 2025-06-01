# Lynx FIM: A Cybersecurity Learning Project

Lynx FIM is a host-based intrusion detection agent (HIDS) I'm building to understand the fundamentals of file integrity monitoring (FIM) and real-time system alerting. 

Built in Go, this project is part of my self-directed study into cybersecurity, specifically focusing on how we can detect unauthorized changes to critical system files.

## Project Goals
- **Understand Integrity:** Explore the use of cryptographic hashes (SHA-256) and HMACs to verify that files haven't been tampered with.
- **Real-time Monitoring:** Learn to use `fsnotify` to listen for kernel-level file events (create, write, delete, rename).
- **Secure Logging:** Practice structured logging and alerting to external webhooks (Slack/Discord) for immediate notification of potential breaches.

## Features I'm Implementing
- [ ] **Baseline Creation:** Scanning and hashing configured directories into a "known-good" baseline.
- [ ] **Real-time Event Handling:** Monitoring file changes as they happen.
- [ ] **HMAC Verification:** Ensuring the baseline itself hasn't been modified.
- [ ] **CLI Interface:** Building a simple, intuitive command-line tool.

## Current Progress (Summer 2025)
- **June:** Initial project setup, basic config parsing, and experimenting with Go's hashing libraries.
- **July:** Implementing the `fsnotify` loop and handling various OS-specific file events.
- **August:** Finalizing alerting logic and refining the baseline verification process.

## How to Run (Development)
1. Clone the repository.
2. Ensure you have Go 1.22+ installed.
3. Run `go mod tidy` to fetch dependencies.
4. Build the binary: `go build -o lynx`.

## Commands
- `lynx init`: Generates a default configuration and a secret key for HMAC.
- `lynx baseline`: Scans target directories and builds the initial integrity database.
- `lynx start`: Begins monitoring the system for changes in real-time.
- `lynx verify`: Manually checks the current file state against the last known baseline.
