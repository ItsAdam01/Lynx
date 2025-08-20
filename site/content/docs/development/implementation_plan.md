---
title: "Implementation Plan"
weight: 4
---

# Roadmap: Building Lynx FIM

This is my long-term implementation plan for Summer 2025. I'm breaking this down into three distinct phases to manage my learning and ensure each part of the system is solid before moving on.

## Phase 1: Scaffolding and Core Logic (June 2025)

The goal of Phase 1 is to build the CLI and the baseline creation logic. This is where I'll get my first real hands-on experience with Go's `crypto` libraries.

1.  **Project Scaffolding:** Set up the Go project with `cobra` for the CLI.
2.  **Configuration System:** Build a YAML-based configuration system for defining watch paths.
3.  **Baseline Generation:**
    *   Walk the file tree recursively.
    *   Calculate SHA-256 hashes for all files.
    *   Implement the HMAC signing logic for the `baseline.json` file.
4.  **CLI Command - `lynx init` and `lynx baseline`:** Successfully create the initial configuration and "known-good" baseline.

## Phase 2: Real-time Monitoring and Detection (July 2025)

The goal of Phase 2 is to move from static hashing to active monitoring. This will be my introduction to the Linux kernel's event system.

1.  **`fsnotify` Integration:** Implement the core monitoring loop that listens for file events.
2.  **Recursive Watching:** Handle the complexity of automatically adding watches when new directories are created.
3.  **Anomaly Logic:** When an event occurs, compare the new file state with the memory-loaded baseline.
4.  **Structured Logging:** Implement JSON logging for all events (Create, Modify, Delete).
5.  **CLI Command - `lynx start`:** Begin active, real-time monitoring of the system.

## Phase 3: Alerting, Testing, and Hardening (August 2025)

The final phase is about making the system robust and useful in a real-world scenario.

1.  **Webhook Integration:** Build the HTTP client to send event payloads to Slack or Discord.
2.  **Alert Filtering:** Implement simple rules to prevent "alert fatigue" (e.g., don't alert on temp files).
3.  **Unit and Integration Testing:** Write tests to ensure the hashing and HMAC logic is 100% correct.
4.  **Cross-Compilation:** Build and test the binary on different Linux distributions (Ubuntu, CentOS).
5.  **Documentation Finalization:** Complete the technical documentation and write a "getting started" guide.

## Beyond Summer 2025: Future Ideas

-   **Centralized Management:** Move away from local JSON logging to a central server.
-   **Automated Remediation:** If a critical file is modified, automatically restore it from a secure backup.
-   **Kernel-Level Filtering:** Explore more advanced ways to filter events at the kernel level for even better performance.
