---
title: "Performance Analysis"
weight: 7
---

# Performance and Efficiency Research

As part of my learning journey, I wanted to understand the performance profile of Lynx FIM. Security tools must be efficient; if an agent consumes too much CPU or takes too long to hash files, it won't be used in production.

## 1. Hashing Throughput (SHA-256)

I ran benchmarks using Go's standard `crypto/sha256` library to see how fast we can process file data. 

### Benchmark Results (August 2025)
| File Size | Time per Operation | Estimated Throughput |
|-----------|--------------------|----------------------|
| 1 MB      | ~0.64 ms           | ~1.5 GB/s            |
| 10 MB     | ~6.48 ms           | ~1.5 GB/s            |

**My Observation:** Go's implementation of SHA-256 is incredibly fast. This confirms that Lynx can handle even large configuration or binary files without significant latency.

## 2. Baseline Generation Speed

I also benchmarked the coordination logic—scanning the file system, calculating hashes, and signing the output.

### Result:
- **100 Files:** ~1.28 ms total.

**My Observation:** The overhead of the file system walker and the HMAC signing is negligible. For most standard server directories (like `/etc/`), establishing a baseline should take less than a second.

## 3. Real-time Monitoring Overhead

While harder to benchmark precisely without specialized tools, my manual testing with `fsnotify` showed:
- **CPU Usage:** Near 0% while idling.
- **Memory:** Minimal (under 20MB) even when watching several hundred files.

This efficiency is due to using Linux `inotify` hooks rather than polling the disk. It allows the agent to stay "asleep" until the kernel notifies it of an event.

## 4. Complexity Analysis (Big O)

To ensure Lynx can scale to larger systems, I've analyzed the theoretical complexity of its core algorithms.

### Phase 1: Initial Baselining
This is the "Heavy Lifting" phase where the agent established the Source of Truth.

*   **Time Complexity: $O(N \cdot S)$**
    *   $N$ = Number of files.
    *   $S$ = Average size of the files.
    *   The agent must walk the directory tree ($O(N)$) and then read every byte of every file to calculate the SHA-256 hash ($O(S)$ per file).
*   **Space Complexity: $O(N \cdot P)$**
    *   $P$ = Average length of the file path string.
    *   The agent stores a map of file paths to their corresponding hashes. This map grows linearly with the number of files being monitored.

### Phase 2: Real-time Monitoring
This is the "Idle Defense" phase where the agent waits for events.

*   **Time Complexity: $O(S_{changed})$**
    *   Detection is $O(1)$ because the Linux kernel pushes events to the agent (no polling required).
    *   When an event occurs, the agent only re-hashes the **specific file** that changed. This takes $O(S)$ where $S$ is the size of that modified file. Comparison with the in-memory baseline is a map lookup, which is $O(1)$ on average.
*   **Space Complexity: $O(N \cdot P)$**
    *   The agent maintains the full baseline in memory to allow for instant comparisons when a file event is received.
