---
title: "Lynx FIM"
type: "docs"
---

# Lynx FIM: A Host-Based Intrusion Detection Agent

Welcome to my project site for **Lynx FIM**, a lightweight host-based intrusion detection agent built in Go. This is where I'm documenting my research, design decisions, and progress as I learn about cybersecurity and real-time system monitoring.

## 🛡️ Project Overview

Lynx FIM is designed to monitor critical system files for unauthorized changes. By establishing a cryptographic "baseline" and then listening for kernel-level file events, the agent can immediately detect and alert when something isn't right.

## 🗺️ Navigation

I've organized the documentation to follow my learning journey:

- **[Introduction]({{< relref "docs/" >}})**: A high-level overview of why I'm building this.
- **[Requirements]({{< relref "docs/requirements" >}})**: What a FIM agent actually needs to do.
- **[Technical Specifications]({{< relref "docs/technical_specs" >}})**: The architecture and research behind the tool.
- **[Implementation Plan]({{< relref "docs/implementation_plan" >}})**: My roadmap for Summer 2025.

---

> **Note:** This project is part of a 2-month intensive learning cycle from June to August 2025. I'm focusing on building a "best practice" implementation to truly understand the core concepts.
