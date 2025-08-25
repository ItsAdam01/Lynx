# Lynx FIM: A Cybersecurity Learning Project

Lynx FIM is a host-based intrusion detection agent (HIDS) I built to understand the fundamentals of file integrity monitoring and real-time system alerting in Go. 

This repository contains the full source code and a detailed documentation site covering my learning journey from June to August 2025.

## 📖 Project Documentation

I have created a comprehensive documentation site using Hugo to organize my research, technical specs, and usage guides.

### How to access the docs:
1. Ensure you have Hugo installed.
2. Navigate to the `site/` directory.
3. Run the development server: `hugo server -D`
4. Visit **`http://localhost:1313/`** in your browser.

You can also browse the documentation source directly:
- **[Usage Guide](site/content/docs/usage/installation.md)**: How to build, configure, and run the agent.
- **[Implementation Story](site/content/docs/development/implementation_story.md)**: My step-by-step learning journey and milestones.

---

## 🚀 Quick Start (Usage)

If you want to jump straight into using the tool:

1. **Build:** `make build`
2. **Init:** `./bin/lynx init`
3. **Configure:** Edit `config.yaml` and set your `LYNX_HMAC_SECRET` env var.
4. **Baseline:** `./bin/lynx baseline`
5. **Monitor:** `./bin/lynx start`

For detailed instructions on Discord webhook setup and manual auditing, see the **[Usage Guide](site/content/docs/usage/installation.md)**.

## 🛠️ Features I've Implemented
- **SHA-256 Hashing:** Unique digital fingerprints for every file.
- **HMAC Protection:** Secure signing of the baseline to prevent tampering.
- **Real-time Monitoring:** Kernel-level event detection via `fsnotify`.
- **Structured Logging:** All security alerts are output as JSON.
- **Asynchronous Alerting:** Background webhook delivery for high performance.

## 💻 Development & Persona
This project is part of my self-directed study into system security. I've documented every hurdle and breakthrough in the **[Implementation Story](site/content/docs/development/implementation_story.md)**.

- **Identity:** Adam Atienza
- **Timeline:** June 2025 – August 2025
- **Goal:** Learn by building a professional-grade security tool from scratch.
