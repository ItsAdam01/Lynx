---
title: "Introduction"
weight: 1
---

# Lynx FIM: A Learning Journey in Host Intrusion Detection

Welcome to the documentation for **Lynx FIM**, a lightweight host-based intrusion detection agent (HIDS) I'm building in Go. This project is a centerpiece of my self-directed study in cybersecurity, specifically focused on file integrity and real-time system monitoring.

## Purpose

As I transition from general software development into the cybersecurity space, I've found that the best way to understand a system's defense is to build it yourself. Lynx FIM is my way of learning:

1.  **Integrity at Scale:** How we establish a "Source of Truth" (a baseline) for thousands of system files and ensure it remains untampered.
2.  **Kernel-level Monitoring:** How we can leverage OS-specific hooks (like Linux `inotify`) to detect unauthorized changes the moment they happen.
3.  **Modern Alerting:** How we can pipe these security events into a modern alerting infrastructure like Slack, Discord, or a centralized SOC.

## Documentation Structure

This documentation site is where I'm organizing my thoughts, technical research, and the long-term implementation plan for the project.

-   **Requirements:** What a modern FIM agent actually *needs* to do.
-   **Technical Specs:** The "how" behind the code—my research into crypto, Go, and system events.
-   **Implementation Plan:** My roadmap for Summer 2025, from a basic CLI to a fully functional HIDS.

> **Note:** This project is part of a 2-month intensive learning cycle from June to August 2025.
