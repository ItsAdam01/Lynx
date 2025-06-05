---
title: "Isolated Lab Testing (Experimental)"
weight: 4
---

# Experimental: Isolated Lab Setup

This guide provides a step-by-step walkthrough for running a temporary, isolated Lynx FIM process. This is ideal for testing detection and alerting without modifying your system files or installing the binary globally.

## 🧪 The "Single Directory" Setup

Follow these steps to create a self-contained testing environment in your `/tmp` directory.

### 1. Prepare the Lab
```bash
# Create and enter a temporary workspace
mkdir -p /tmp/lynx-lab && cd /tmp/lynx-lab

# Build the latest binary from your project root
# (Assuming you are in the project root for the make command)
make build
cp bin/lynx /tmp/lynx-lab/
cd /tmp/lynx-lab
```

### 2. Create Dummy Data to Watch
```bash
mkdir watched_dirs
echo "secret info" > watched_dirs/top_secret.txt
```

### 3. Initialize and Configure Locally
```bash
./lynx init

# Update the config to watch our lab directory instead of system paths
# Also ensures the secret key is read from the correct environment variable
sed -i 's|/etc/ssh|./watched_dirs|g' config.yaml
sed -i 's|hmac_secret_env: "LYNX_HMAC_SECRET"|hmac_secret_env: "LYNX_HMAC_SECRET"|g' config.yaml
```

### 4. Set Secret and Establish Baseline
```bash
export LYNX_HMAC_SECRET="lab-secret-123"
./lynx baseline -o lab_baseline.json
```

### 5. Start Monitoring
Run this command to start the agent. Note that this will block the current terminal window as it listens for events.
```bash
./lynx start -b lab_baseline.json
```

## 🔍 Verifying Detection

While the process is running in **Terminal 1**, open a **second terminal window** and trigger a tampering event:

```bash
cd /tmp/lynx-lab
echo "tampered!" >> watched_dirs/top_secret.txt
```

### Expected Output (Terminal 1)
You should immediately see the alert in your first terminal:
`CRITICAL: File modified: /tmp/lynx-lab/watched_dirs/top_secret.txt`

---

## 🗺️ Navigation
- **[Installation & Setup]({{< relref "installation.md" >}})**: Permanent installation and service setup.
- **[Command Reference]({{< relref "commands.md" >}})**: Detailed syntax for all Lynx commands.
- **[Back to Introduction]({{< relref "../_index.md" >}})**
