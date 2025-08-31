---
title: "Command Reference"
weight: 3
---

# CLI Command Reference

Lynx FIM is operated through a simple but powerful command-line interface. All commands support the global `--config` flag to specify a custom configuration path.

## Global Flags
- `--config string`: Path to the configuration file (default is "config.yaml").
- `-h, --help`: Display help information for any command.

---

## 1. `lynx init`
Initializes the workspace for Lynx FIM.

### Usage
```bash
lynx init [flags]
```

### What it does:
Creates a boilerplate `config.yaml` file in the current directory with sensible defaults. This is always the first step when setting up a new agent.

---

## 2. `lynx baseline`
Establishes the "Source of Truth" for your monitored files.

### Usage
```bash
lynx baseline [flags]
```

### Flags
- `-o, --output string`: Path where the signed baseline JSON will be saved (default "baseline.json").

### What it does:
1. Reads your `config.yaml` to find the target paths.
2. Scans every file and calculates its SHA-256 hash.
3. Signs the resulting data with your `LYNX_HMAC_SECRET`.
4. Saves the results to disk. **Note:** You must run this whenever you intentionally change system files so the agent knows the new state is valid.

---

## 3. `lynx start`
Activates real-time intrusion detection.

### Usage
```bash
lynx start [flags]
```

### Flags
- `-b, --baseline string`: Path to the verified baseline file (default "baseline.json").

### What it does:
This is the main long-running process.
1. Loads the baseline and verifies its HMAC signature.
2. Initializes the `inotify` watcher for all configured paths.
3. Starts the background Alert Dispatcher.
4. Logs any file creation, modification, or deletion to the structured JSON log and sends a webhook alert if configured.

---

## 4. `lynx verify`
Performs a one-off manual integrity audit.

### Usage
```bash
lynx verify [flags]
```

### Flags
- `-b, --baseline string`: Path to the verified baseline file (default "baseline.json").

### What it does:
Use this for manual sweeps. It performs a full scan of the file system and compares it against the baseline immediately, printing a summary of discrepancies to the terminal. Unlike `start`, this command exits as soon as the audit is complete.

---

## 🗺️ Navigation
- **[Installation & Setup]({{< relref "installation.md" >}})**: How to get the binary.
- **[Isolated Lab Testing]({{< relref "isolated_testing.md" >}})**: Safe testing tutorial.
- **[Back to Introduction]({{< relref "../_index.md" >}})**
