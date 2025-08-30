---
title: "Installation & Setup"
weight: 1
---

# Installation and Setup Guide

This guide covers everything you need to get Lynx FIM running on your system, from compiling the source code to choosing which directories to monitor.

## 1. Supported Platforms

Lynx FIM is built in Go and leverages the Linux kernel's `inotify` system for real-time monitoring.

- **Primary Support:** Linux (Ubuntu, Debian, CentOS, RHEL, etc.).
- **Architecture:** AMD64 and ARM64.
- **Other OS:** While it will compile on macOS or Windows, the real-time recursive monitoring features are optimized for Linux.

## 2. Build and Installation

### Prerequisites
- **Go 1.22+**: Required to compile the source.
- **Make**: Recommended for using the automated build script.

### Compilation Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/ItsAdam01/Lynx.git
   cd Lynx
   ```
2. Build for your current architecture:
   ```bash
   make build
   ```
3. (Optional) Install to your path:
   ```bash
   sudo cp bin/lynx /usr/local/bin/
   ```

## 3. Configuration Setup

### Step 1: Initialize
Generate a default configuration file in your current directory:
```bash
lynx init
```

### Step 2: The HMAC Secret
Lynx requires a secret key to protect the integrity of its baseline. You must set this as an environment variable. I learned that keeping secrets out of config files is a major security best practice.

```bash
export LYNX_HMAC_SECRET="use-a-very-long-random-string-here"
```

### Step 3: Configure `config.yaml`
Open the generated `config.yaml`. Here is a sample of how I configured mine:

```yaml
agent_name: "prod-web-server-01"
log_file: "./lynx.log"
webhook_url: "https://discord.com/api/webhooks/..." # Your Discord Webhook

# Directories to monitor recursively
paths_to_watch:
  - "/etc/ssh"
  - "/etc/pam.d"
  - "/usr/local/bin"

# Specific critical files
files_to_watch:
  - "/etc/passwd"
  - "/etc/shadow"
  - "/etc/hosts"
```

## 4. Directory Recommendations

Choosing what to watch is a balance between security and performance. Here is what I discovered during my research:

### ✅ Recommended to Watch
- **/etc/**: Contains almost all system configuration files.
- **/usr/bin/ or /usr/local/bin/**: Where system executables live. Watching these helps detect "binary replacement" attacks.
- **/root/.ssh/**: To monitor for unauthorized SSH keys being added.

### ❌ Not Recommended to Watch
- **/proc/ or /sys/**: These are virtual file systems managed by the kernel; they change constantly and will flood you with alerts.
- **/var/log/**: Logs change every second. Monitoring them with a FIM will cause a loop of alerts.
- **/tmp/**: High-noise directory used by many applications for temporary storage.
