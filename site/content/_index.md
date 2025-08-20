---
title: "Lynx FIM"
type: "docs"
---

# Lynx FIM: A Host-Based Intrusion Detection Agent

Welcome to the project site for **Lynx FIM**, a lightweight host-based intrusion detection agent built in Go. This is where I'm documenting my research, design decisions, and progress as I learn about cybersecurity and real-time system monitoring.

## 🛡️ Project Overview

Lynx FIM monitors critical system files for unauthorized changes. By establishing a cryptographic "baseline" and then listening for kernel-level file events, the agent can immediately detect and alert when something is wrong.

## 🗺️ Documentation Map

I've organized the documentation to serve two different purposes:

### [🚀 Usage Guide]({{< relref "docs/usage/installation" >}})
Learn how to install, configure, and run Lynx on your own servers.
- **[Installation & Setup]({{< relref "docs/usage/installation" >}})**: Binary builds and Discord integration.
- **[General Features]({{< relref "docs/usage/features" >}})**: What Lynx can do for your security.

### [💻 Development & Research]({{< relref "docs/development/technical_specs" >}})
Deep dive into the architecture, the learning journey, and how I built the tool.
- **[Technical Specifications]({{< relref "docs/development/technical_specs" >}})**: Go, Cryptography, and Kernel events.
- **[Implementation Story]({{< relref "docs/development/implementation_story" >}})**: My milestones and lessons learned during Summer 2025.
- **[Proof of Concept]({{< relref "docs/development/demonstration" >}})**: Seeing Lynx in action with actual log outputs.

---

> **Note:** This project is part of a 2-month intensive learning cycle from June to August 2025. I'm focusing on building a "best practice" implementation to understand core security concepts.
