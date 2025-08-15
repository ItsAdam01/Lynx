---
title: "Implementation Story"
weight: 5
---

# The Implementation Story: A Developer's Log

This page is where I'm documenting the "story" of building Lynx FIM. It's more than just a list of features, it's a record of the technical hurdles I've faced, the breakthroughs I've had, and what I'm learning along the way as I dive into cybersecurity.

## Milestone 1: The Foundation (June 10, 2025)

I started by laying the groundwork for the CLI. I'm using `cobra` because it feels like the industry standard for Go tools. Getting the project structure right was my first challenge—Go has specific conventions for `cmd/` and `internal/`, and I want to make sure I'm following them from day one.

**Hurdle:** I initially struggled with how to handle configuration. I want to keep the tool simple but flexible. I decided on `viper` for YAML support, which lets me define watch paths easily.

## Milestone 2: Cryptographic Wins (June 12, 2025)

Today, I implemented the core "integrity" logic: SHA-256 hashing and HMAC signing. This is the heart of any FIM agent.

**The Lesson:** I learned that even simple things like hashing a string can be tricky. During my first unit test, I had a hash mismatch. I spent an hour debugging only to realize that my `echo -n` command and my Go test content had a slight difference in how they handled newlines. 

**Breakthrough:** I successfully implemented HMAC-SHA256 for signing payloads. This was a big "aha!" moment for me. I now understand how a secret key can be used to prove that a file (like our baseline) hasn't been modified by an unauthorized party.

## Milestone 3: Walking the File System (June 14, 2025)

I've now implemented the "scanner" part of the agent. It recursively walks through the directories I've configured and gathers all the files it needs to hash.

**Hurdle:** One challenge I faced was handling relative paths. I want to make sure the agent is consistent, so I decided to resolve everything to absolute paths. This ensures that even if I run the agent from different directories, the "Source of Truth" remains stable.

## Milestone 4: Secure Baseline Storage (June 16, 2025)

This was a major milestone. I've combined the file scanner with my cryptographic logic to create the `Baseline` storage system.

**Breakthrough:** I've implemented a system where the `baseline.json` file is signed with an HMAC. This means that if anyone tries to tamper with my baseline file (like changing a hash to hide their tracks), the agent will detect it immediately when it tries to load the file. 

**The Lesson:** I learned about the importance of consistent JSON marshalling. To verify the signature, I have to make sure the data I'm re-hashing is exactly the same as what I signed originally. I found that resetting the `Signature` field to an empty string before marshalling for verification is a clean way to handle this.

## Milestone 5: Learning to Trust Tests (June 18, 2025)

As I move deeper into this project, I've decided to fully embrace Test-Driven Development (TDD). 

### Why TDD?
At first, it felt slow. Why write a test before the code? But as I worked on the hashing logic, I realized that in cybersecurity, "close enough" is a failure. If my integrity checks are slightly off, the whole system is useless. TDD forces me to define exactly what "success" looks like before I type a single line of implementation logic.

### How I'm applying it:
My workflow is now:
1. **Red:** Write a failing test for a small piece of functionality (e.g., "The init command should create a file with specific default values").
2. **Green:** Write just enough code to make the test pass.
3. **Refactor:** Clean up the code while ensuring the test stays green.

### Where it's helping:
I applied this to the `lynx init` command today. By writing the test first, I was forced to think about the default configuration file format and ensure that even if I add new features later, the core initialization logic will still work as expected.

## Milestone 6: Coordinating the Baseline (June 22, 2025)

Today I tied everything together with the `baseline` command. This was my first time coordinating multiple internal packages (`crypto`, `fs`, `config`) to perform a complex task.

**Breakthrough:** Using TDD to test the coordination logic in `internal/app` was a game changer. I could verify that the whole process—scanning files, hashing them, and saving a signed baseline—was working correctly without having to manually run the CLI every time.

**The Lesson:** I learned about the value of mocks and environment variables in testing. To test the `baseline` command, I had to ensure the `LYNX_HMAC_SECRET` was correctly handled. Setting and unsetting environment variables in my tests made them reliable and isolated.

## Milestone 7: Into the Kernel - Real-time Monitoring (July 1, 2025)

Phase 1 was about the "Source of Truth" (the baseline). Phase 2 is about "Active Defense." I'm starting to implement real-time monitoring using `fsnotify`.

**The Research:** I've been reading about how the Linux kernel handles file events. Instead of my agent constantly scanning the disk—which would be slow and resource-heavy—I can use `inotify` (via `fsnotify`) to have the kernel "ping" my agent the moment a file is touched. This is a huge step up in efficiency.

**The Goal for July:** By the end of this month, I want `lynx start` to be a long-running process that watches my configured paths and logs any change immediately. I'm excited but also a bit nervous about handling the complexity of recursive directory watching.

## Milestone 8: Recursive Watching and Real-time Detection (July 3, 2025)

I've successfully implemented recursive directory watching today. This was a big technical hurdle for me. When a user creates a new folder within a watched path, my agent now automatically adds that folder to its monitoring queue.

**The Breakthrough:** I implemented a coordination loop that compares every file event against the memory-loaded baseline. If a file is modified, I re-hash it and compare the new signature with the old one. If it's a new file, I log a warning. If it's deleted, I log a critical alert.

**The Lesson:** I learned about the importance of handling OS events correctly. For example, when a new directory is created, `fsnotify` gives me a `Create` event. I have to immediately add that directory to the watcher so I don't miss any files created inside it a split second later.

## Milestone 9: Structured Logging and the Start Command (July 8, 2025)

Phase 2 is now officially complete. I've wired up the real-time monitor to the CLI with the new `lynx start` command, and I've implemented structured JSON logging.

**The Research:** I learned that in the enterprise security world, simple text logs aren't enough. Security Information and Event Management (SIEM) systems like Splunk or ELK need structured data. I decided to use Go's standard library `log/slog` package to output all events as JSON.

**The Breakthrough:** Bringing everything together in `cmd/start.go` was incredibly satisfying. The agent now loads the configuration, verifies the signed baseline, initializes the JSON logger, and blocks indefinitely while listening for file system events. It even handles termination signals (`SIGINT`, `SIGTERM`) gracefully.

## Milestone 10: Real-time Alerting and Manual Audits (August 3, 2025)

It is August 2025, and I'm entering the final phase of my initial learning roadmap. This month is about getting these alerts out of the log files and into a platform like Slack or Discord.

**The Breakthrough:** I successfully implemented the Webhook alerting pipeline today. Using TDD, I verified that my agent can now send a structured JSON payload to any configured webhook URL. This means I can get security alerts on my phone the moment a critical system file is touched.

**The Audit Tool:** I've also implemented the `lynx verify` command. This is useful for manual audits where I want to do a "clean sweep" and compare the entire system against the baseline without running a persistent agent. 

**The Lesson:** I learned about the power of `net/http` and `httptest` in Go. Writing tests for the webhook required me to mock a web server, which was a great exercise in understanding how HTTP requests are actually structured and sent.

## Milestone 11: Speed and Security - Asynchronous Alerting (August 5, 2025)

As I tested the agent, I noticed a problem: if the webhook server is slow, my whole monitoring loop blocks while it waits for a response. In a security tool, that's unacceptable. Every millisecond of delay is a window for an attacker.

**The Breakthrough:** I implemented an asynchronous `AlertDispatcher` today using Go's channels and goroutines. Now, when the agent detects an anomaly, it simply "drops" the alert into a channel and gets back to monitoring immediately. A separate background process picks up the alert and handles the network delivery.

**The Lesson:** This was my first real experience with Go's concurrency patterns in a production-like scenario. Learning how to use a `select` statement to handle both outgoing alerts and a "stop" signal was a major milestone for me. It makes the agent feel much more professional and robust.

## Milestone 12: The Final Connection (August 8, 2025)

Today I officially "closed the loop" by integrating the asynchronous alert dispatcher into the `lynx start` command. 

**The Breakthrough:** It was a moment of pure satisfaction to see all the pieces working together. The agent now initializes the monitor, starts the background dispatcher, and then sits in a non-blocking loop waiting for file events. When an anomaly is detected, it's logged to JSON and then immediately "fired off" to the webhook channel.

**The Lesson:** I learned about the importance of channel buffering. By giving my `anomalies` and `alertChan` channels a small buffer, I've made the system even more resilient to bursts of file system activity. It's a small detail, but in a security tool, it's the difference between catching every event and missing a critical breach.

## Milestone 13: Ready for Deployment - Build Automation (August 15, 2025)

As I wrap up this project, I've moved from writing code to thinking about how others will use it. I've implemented a `Makefile` to handle building, testing, and cross-compiling the Lynx FIM agent.

**The Breakthrough:** With one command, I can now run my entire test suite and build binaries for both `amd64` and `arm64` Linux servers. This is a major step toward making the agent "production-ready." It feels like I've built a real tool, not just a learning project.

**The Lesson:** I learned that automation is just as important as the code itself. By building the testing into my `Makefile`, I've ensured that I never accidentally ship a binary that hasn't passed all my integrity checks. 

### Technical Achievements (Phase 3 Complete):
- [x] Verified SHA-256 hashing for files.
- [x] Implemented constant-time HMAC comparison to prevent timing attacks.
- [x] Established a strict test-driven development (TDD) workflow.
- [x] Recursive file system traversal with absolute path resolution.
- [x] HMAC-signed JSON storage for the baseline with tamper detection.
- [x] CLI implementation for `lynx init`, `lynx baseline`, `lynx start`, and `lynx verify`.
- [x] Core real-time monitoring with `fsnotify`.
- [x] Recursive directory watching and anomaly detection logic.
- [x] Structured JSON logging and Asynchronous Webhook alerting.
- [x] Fully integrated, non-blocking alerting pipeline.
- [x] Automated build system and cross-compilation with `Makefile`.

> "A tool is only as good as the process that builds it. Automation is the final layer of defense." - *Completing the August 2025 roadmap.*

> "A security tool that only looks back is a historian. A security tool that looks at the present is a defender." - *Closing the loop on real-time defense.*
