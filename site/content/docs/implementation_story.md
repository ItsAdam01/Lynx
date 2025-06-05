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

### Technical Achievements:
- [x] Verified SHA-256 hashing for files.
- [x] Implemented constant-time HMAC comparison to prevent timing attacks.
- [x] Established a strict test-driven development (TDD) workflow.
- [x] Recursive file system traversal with absolute path resolution.
- [x] HMAC-signed JSON storage for the baseline with tamper detection.
- [x] CLI implementation for `lynx init` and `lynx baseline`.
- [x] Core real-time monitoring with `fsnotify`.
- [x] Recursive directory watching and anomaly detection logic.

> "A security tool that only looks back is a historian. A security tool that looks at the present is a defender." - *Transitioning to Real-time defense.*
