---
title: "General Features"
weight: 2
---

# Lynx FIM: Complete Feature Guide

Lynx FIM is a focused, high-performance host-based intrusion detection agent. This page provides a comprehensive breakdown of every feature I've implemented during this project.

## 1. Cryptographic Integrity (The Baseline)
At the heart of Lynx is the "Source of Truth"—a cryptographic record of every file you choose to protect.

-   **SHA-256 Hashing:** Every monitored file is fingerprinted using the industry-standard SHA-256 algorithm. Even a single bit change in a file will result in a completely different hash.
-   **HMAC Protection:** The `baseline.json` file is signed with a Hash-based Message Authentication Code (HMAC). This ensures that even if an attacker modifies the baseline file, the agent will detect the tampering and refuse to trust it.
-   **Configuration Locking:** Lynx hashes your `config.yaml` and stores that hash in the baseline. This prevents unauthorized changes to your watch paths or ignore lists after the baseline has been established.

## 2. Real-time Monitoring
Lynx doesn't just scan files; it actively defends your system using kernel-level event hooks.

-   **fsnotify Integration:** Leverages the Linux `inotify` system to receive instant notifications from the kernel whenever a file is touched.
-   **Event Detection:** Specifically monitors for:
    -   **Creation:** New files appearing in watched directories.
    -   **Modification:** Changes to the content of existing monitored files.
    -   **Deletion:** Monitored files being removed.
    -   **Renaming:** Automatically detected and reported as a deletion of the original path.
-   **Recursive Watching:** When you watch a directory, Lynx automatically monitors every subdirectory within it. If you create a new folder, Lynx hooks into it immediately without requiring a restart.

## 3. Advanced Filtering & Noise Reduction
To be useful in production, a FIM must be quiet. I've implemented a robust "ignore" mechanism to reduce alert fatigue.

-   **Ignore Patterns:** Support for `.gitignore`-style patterns in your configuration. You can exclude noisy files (like `*.log`, `*.tmp`, or `.DS_Store`) using standard shell globbing.
-   **Global & Local Paths:** Ignore rules apply recursively to all monitored directories.

## 4. Professional Alerting Pipeline
Detecting a breach is only half the battle; the other half is making sure the right people know about it instantly.

-   **Structured JSON Logging:** Every security event is logged as a machine-readable JSON object. This is perfect for integration with SIEM platforms like Splunk, ELK, or Datadog.
- **Asynchronous Webhooks:** Alerts are dispatched to Discord or Slack in the background using a non-blocking queue. This ensures that network latency never slows down the core monitoring loop.
- **Agent Identification:** Every alert includes an `agent_name` (configured in `config.yaml`). This serves as a unique identifier for the host, allowing security teams to quickly pinpoint which server in a fleet is reporting the incident.
- **Semantic Labeling:** Webhook payloads use professional, text-based semantic labels (e.g., `[CRITICAL]`, `[WARNING]`) instead of excessive icons. This ensures that the output is clean, readable, and focused on the data.


## 5. Flexible Operation Modes
Lynx is designed to be both a persistent defender and a manual audit tool.

-   **Persistent Agent (`start`):** A long-running background process with graceful shutdown handling.
-   **Manual Audit (`verify`):** A one-off command that performs a comprehensive "clean sweep" comparison of your entire system against the baseline and prints a detailed report.
-   **Easy Scaffolding (`init`):** Instantly generate a boilerplate configuration to get started in seconds.

---

## 🗺️ Navigation
- **[Installation & Setup]({{< relref "installation.md" >}})**: How to get started.
- **[Command Reference]({{< relref "commands.md" >}})**: Detailed syntax.
- **[Back to Introduction]({{< relref "../_index.md" >}})**
