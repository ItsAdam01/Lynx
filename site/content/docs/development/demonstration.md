---
title: "Demonstration & Proof of Concept"
weight: 6
---

# Proof of Concept: Lynx FIM in Action

This page demonstrates the Lynx FIM agent's ability to detect and report unauthorized file system changes in real-time. I've captured the actual CLI and JSON log output from my final August 2025 audit.

## Scenario: Manual Integrity Audit

In this scenario, I established a baseline of a test directory and then manually tampered with a "critical system file." When I ran the `verify` command, the agent immediately flagged the discrepancy.

### Manual Audit CLI Output
```text
Starting manual integrity audit for agent: default-agent
Comparing against baseline from: 2025-08-15 16:40:30
❌ 1 Anomaly(s) Detected:
  - CRITICAL: File modified: /home/adamatienza/Lynx/test_audit/critical.txt
```

![Actual Discord Alerts Captured](images/discord-alerts.png)

## Scenario: Real-time Monitoring and Alerting

Next, I started the agent in the background and triggered another file modification. The agent immediately captured the event and logged it as a structured JSON object, ready for ingestion by a SIEM system.

### Structured JSON Log Output (`lynx.log`)
```json
{
  "time": "2025-08-15T16:40:30.716Z",
  "level": "INFO",
  "msg": "Starting Lynx FIM agent",
  "agent_name": "default-agent",
  "total_baseline_files": 27
}
{
  "time": "2025-08-15T16:40:32.715Z",
  "level": "WARN",
  "msg": "Anomaly detected",
  "details": "CRITICAL: File modified: /home/adamatienza/Lynx/test_audit/secret.txt"
}
{
  "time": "2025-08-15T16:40:34.716Z",
  "level": "INFO",
  "msg": "Shutting down Lynx FIM",
  "signal": "terminated"
}
```

## Observations and Lessons

- **Precision:** The SHA-256 hashing correctly identifies even a single-character change to a file.
- **Structure:** The JSON logs are perfect for further analysis. Each event includes a timestamp, a severity level, and specific details about the anomaly.
- **Graceful Shutdown:** The agent handles signals like `Ctrl+C` cleanly, ensuring that all background processes (like the alert dispatcher) are properly closed.

---

> "Watching the logs turn red after I tampered with my own test files was a huge moment. It's the point where all the theory about SHA-256 and fsnotify became a real, working tool." - *Finalizing the Proof of Concept.*
