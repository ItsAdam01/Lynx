---
title: "Requirements"
weight: 2
---

# Requirements: Building a Modern FIM

As I research existing host-based intrusion detection systems (HIDS) like OSSEC or Wazuh, I'm identifying the core functionality that a modern File Integrity Monitor (FIM) actually needs. For Lynx FIM, I'm prioritizing the following requirements for my initial learning build.

## Functional Requirements

The FIM agent must be able to perform these core tasks:

1.  **Baseline Creation:**
    *   **Recursive Scanning:** Scan configured directories (like `/etc/`, `/usr/bin/`) and identify all files.
    *   **Hashing:** Generate a SHA-256 cryptographic hash for every file to create its unique "fingerprint."
    *   **Secure Storage:** Save the "known-good" baseline to a file for later comparison.
2.  **Real-time Monitoring:**
    *   **Event Detection:** Listen for kernel-level file events: `Create`, `Write`, `Delete`, and `Rename`.
    *   **Anomaly Comparison:** When an event occurs, immediately compare the new file state against the baseline to identify the change.
3.  **Alerting & Logging:**
    *   **Structured Output:** Log all events in a machine-readable JSON format for easy parsing and ingestion.
    *   **Immediate Alerting:** Send a POST request to a pre-configured webhook (Slack/Discord) for any high-priority change (e.g., a change to `/etc/passwd`).
4.  **Self-Protection:**
    *   **Integrity Verification:** The baseline file itself must be signed with an HMAC (Hash-based Message Authentication Code).
    *   **Tamper Detection:** The agent must refuse to start if the baseline's signature doesn't match its current content.

## Non-Functional Requirements

These are the "quality" goals I'm aiming for as a developer:

*   **Performance:** The agent must have a minimal CPU and memory footprint. It should watch thousands of files without impacting server performance.
*   **Zero Dependencies:** By using Go, I want a single, statically linked binary that I can drop onto any Linux server (Ubuntu, CentOS, etc.) and it "just works."
*   **Security:** The HMAC secret key must be handled securely (e.g., via an environment variable) and never hardcoded in the repository.

## My Learning Goals (Why I'm Doing This)

*   **Go Systems Programming:** Learn how to build long-running, concurrent services in Go.
*   **Linux Security Fundamentals:** Understand how the Linux kernel notifies userspace about file changes (`inotify`).
* **Cryptography Basics:** Get hands-on experience with SHA-256 and HMAC to learn how data integrity is actually enforced in the real world.

---

## 🗺️ Navigation
- **[Technical Specifications]({{< relref "technical_specs.md" >}})**: How these requirements are built.
- **[Implementation Plan]({{< relref "implementation_plan.md" >}})**: Timeline for the build.
- **[Back to Introduction]({{< relref "../_index.md" >}})**

