---
title: "General Features"
weight: 2
---

# Lynx FIM Features

Lynx is designed to be a lightweight, focused tool for host-based intrusion detection. Here is what I've implemented to ensure system integrity.

## 🛡️ Cryptographic Integrity
Lynx establishes a "Source of Truth" by creating a baseline of your system files.
- **SHA-256 Hashing:** Every file is fingerprinted with a unique cryptographic hash.
- **HMAC Protection:** The baseline file itself is signed with an HMAC secret. If an attacker modifies the baseline to hide their tracks, Lynx will detect the signature mismatch and refuse to start.

## 📡 Real-time Monitoring
Instead of constantly scanning the disk, which is slow, Lynx uses kernel-level hooks.
- **Kernel Events:** Lynx listens for `Create`, `Write`, `Delete`, and `Rename` events via the Linux `inotify` system.
- **Recursive Watching:** If you create a new directory inside a watched path, Lynx automatically adds it to the monitoring queue in real-time.

## 🚨 Professional Alerting
Detecting an event is only useful if someone is notified.
- **Structured JSON Logs:** All security events are logged in a machine-readable format, making it easy to integrate with SIEM platforms.
- **Asynchronous Webhooks:** Alerts are dispatched in the background. This ensures that a slow network connection never slows down the monitoring loop.

## 🛠️ Manual Audits
Sometimes you don't want a long-running process. Lynx includes a `verify` command for one-off manual comparisons against your stored baseline.
