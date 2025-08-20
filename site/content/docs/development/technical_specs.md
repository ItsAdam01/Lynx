---
title: "Technical Specifications"
weight: 3
---

# Technical Architecture: The Inner Workings of Lynx FIM

This document details my technical research and the design choices I've made for the Lynx FIM architecture. I'm focusing on "best practice" implementations for integrity and performance as I learn the ropes of HIDS development.

## 1. Core Technology: Why Go?

I've chosen Go (Golang) as the foundation for this project for several reasons:
- **Low-Level Control with High-Level Productivity:** Go gives me the performance I need for file hashing while providing a safe and productive environment for building the alerting system.
- **Static Compilation:** A single binary with no external dependencies is a massive security benefit. I don't need to worry about the target server having the right version of Python or OpenSSL.
- **Concurrency:** Go's goroutines will let me handle real-time file events and webhook alerting concurrently without blocking the main monitoring loop.

## 2. Cryptographic Strategy

To ensure data integrity, I'm using a two-layered approach:

### SHA-256 for File Hashing
I'm using the `crypto/sha256` package from Go's standard library. 
- **Reasoning:** SHA-256 is the industry standard for integrity checks. It's collision-resistant and provides a 256-bit hash that's perfect for identifying even a single-bit change in a multi-gigabyte file.

### HMAC for Baseline Protection
I'm using `crypto/hmac` to protect the `baseline.json` file.
- **Problem:** If an attacker modifies a system file, they might try to also modify the baseline to match the new hash.
- **Solution:** When I generate the baseline, I'll also generate a signature of the entire file using an HMAC secret. The agent will re-calculate this signature every time it starts. Without the secret key, an attacker cannot generate a valid signature for a tampered baseline.

## 3. Real-Time Event Loop

I'm using the `fsnotify` library, which wraps the Linux `inotify` system call.

### How it Works:
1.  **Watch Initialization:** The agent recursively walks through the configured directories and adds an `inotify` watch to each one.
2.  **Kernel Hook:** When a file event occurs (e.g., `IN_MODIFY`), the kernel puts a message into a buffer for the agent.
3.  **Event Handler:** The agent reads from this buffer, identifies the file path, and triggers the anomaly detection logic.

## 4. Data Structures & Schemas

### Configuration File (`config.yaml`)
```yaml
agent_name: "dev-lab-01"
hmac_secret_env: "LYNX_HMAC_SECRET"
log_file: "/var/log/lynx.log"
webhook_url: "https://discord.com/api/webhooks/..."

directories_to_watch:
  - "/etc/ssh"
  - "/usr/local/bin"

files_to_watch:
  - "/etc/passwd"
  - "/etc/hosts"
```

### Baseline Storage (`baseline.json`)
```json
{
  "metadata": {
    "generated_at": "2025-06-20T14:00:00-04:00",
    "total_files": 128
  },
  "hashes": {
    "/etc/passwd": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "/etc/ssh/sshd_config": "8b1a9953c4611296a827abf8c47804d7e6c49c6b"
  },
  "signature": "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"
}
```

## 5. Security Constraints

- **Sensitive Data:** The `LYNX_HMAC_SECRET` must be provided as an environment variable. The agent will panic and exit if it's not present.
- **Privilege:** The agent will need to run with `sudo` (root) permissions to read sensitive system files and create the `inotify` watches.
