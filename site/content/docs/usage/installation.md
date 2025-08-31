---
title: "Installation & Setup"
weight: 1
---

# Installation and Setup Guide

This guide covers everything you need to get Lynx FIM running on your system, from compiling the source code to ensuring it stays running after a reboot.

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
log_file: "/var/log/lynx.log"
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

## 4. Running the Agent

### Permissions
Because Lynx needs to read sensitive system files (like `/etc/shadow`) and hook into kernel events, it **must be run with root privileges**.

```bash
sudo LYNX_HMAC_SECRET="your-secret" ./bin/lynx start
```

## 5. Persistence (Running as a Service)

To ensure Lynx FIM starts automatically when the server boots and restarts if it fails, you should use `systemd`. This is how professional HIDS agents are deployed.

### Step 1: Create the Service File
Create a file at `/etc/systemd/system/lynx.service`:

```ini
[Unit]
Description=Lynx File Integrity Monitor
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/lynx
# Pass the secret via environment variable
Environment=LYNX_HMAC_SECRET=your-super-long-secret-key
ExecStart=/usr/local/bin/lynx start --config /opt/lynx/config.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

### Step 2: Enable and Start
Run the following commands to activate the service:

```bash
# Reload systemd to recognize the new service
sudo systemctl daemon-reload

# Enable the service to start on boot
sudo systemctl enable lynx

# Start the service immediately
sudo systemctl start lynx

# Check the status
sudo systemctl status lynx
```

## 6. Directory Recommendations

### ✅ Recommended to Watch
- **/etc/**: System configuration files.
- **/usr/bin/ or /usr/local/bin/**: System executables.
- **/root/.ssh/**: Unauthorized SSH keys.

### ❌ Not Recommended to Watch
- **/proc/ or /sys/**: Virtual kernel file systems.
- **/var/log/**: High-frequency log writes.
- **/tmp/**: High-noise temporary storage.

---

## 🗺️ Navigation
- **[Command Reference]({{< relref "commands.md" >}})**: How to operate the agent.
- **[General Features]({{< relref "features.md" >}})**: What Lynx can do for you.
- **[Back to Introduction]({{< relref "../_index.md" >}})**
