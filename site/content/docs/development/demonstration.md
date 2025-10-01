---
title: "Demonstration & Proof of Concept"
weight: 6
---

# Proof of Concept: Lynx FIM in Action

This page demonstrates the Lynx FIM agent's ability to detect and report unauthorized file system changes in real-time. Below is a correlation between the terminal commands and the resulting Discord alerts.

## 🧪 Laboratory Test Scenario

In this audit, I used an isolated lab directory to verify the full lifecycle of the agent: initialization, baselining, and real-time detection.

### 1. Establishing the Source of Truth
First, I generated the signed cryptographic baseline for the test directory.

**Terminal Input:**
```bash
./lynx baseline -o lab_baseline.json
```

**Terminal Output:**
```text
Successfully created baseline: lab_baseline.json 
```

### 2. Real-time Monitoring and Alerting
Next, I started the monitoring agent and triggered several file system events (modifications, deletions, and additions).

**Terminal Input:**
```bash
./lynx start -b lab_baseline.json
```

**Live Event Log:**
```text
[CRITICAL] FILE_MODIFIED: ./test-dir/test2
[CRITICAL] FILE_DELETED: ./test-dir/test2.txt
[WARNING] FILE_CREATED: ./test-dir/testrename.txt
[WARNING] FILE_CREATED: ./test-dir/testadd
```

### 3. Visual Verification (Discord)
The following image shows exactly how these events were dispatched and rendered in the Discord security channel. Note the semantic labeling and emojis used to distinguish between warnings and critical breaches.

![Lynx FIM Discord Alerts Captured](/images/discord-alerts.png)

## Observations and Lessons

- **Precision:** The agent correctly distinguished between a file modification (CRITICAL) and a new file creation (WARNING).
- **Responsiveness:** Alerts appeared in the Discord channel within milliseconds of the file being touched in the lab.
- **Data Integrity:** The use of absolute paths in the final webhook payload ensures that the security analyst knows exactly where the event occurred on the host.

---

## 🗺️ Navigation
- **[Performance Analysis](../development/performance.md)**: Efficiency and scalability research.
- **[Back to Introduction](../../_index.md)**
