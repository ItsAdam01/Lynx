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

I also benchmarked the coordination logic: scanning the file system, calculating hashes, and signing the output.

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

*   **Time Complexity:** {{< katex >}}O(N \cdot S){{< /katex >}}
    *   {{< katex >}}N{{< /katex >}} = Number of files.
    *   {{< katex >}}S{{< /katex >}} = Average size of the files.
    *   The agent must walk the directory tree ({{< katex >}}O(N){{< /katex >}}) and then read every byte of every file to calculate the SHA-256 hash ({{< katex >}}O(S){{< /katex >}} per file).
*   **Space Complexity:** {{< katex >}}O(N \cdot P){{< /katex >}}
    *   {{< katex >}}P{{< /katex >}} = Average length of the file path string.
    *   The agent stores a map of file paths to their corresponding hashes. This map grows linearly with the number of files being monitored.

### Phase 2: Real-time Monitoring
This is the "Idle Defense" phase where the agent waits for events.

*   **Time Complexity:** {{< katex >}}O(S_{changed}){{< /katex >}}
    *   Detection is {{< katex >}}O(1){{< /katex >}} because the Linux kernel pushes events to the agent (no polling required).
    *   When an event occurs, the agent only re-hashes the **specific file** that changed. This takes {{< katex >}}O(S){{< /katex >}} where {{< katex >}}S{{< /katex >}} is the size of that modified file. Comparison with the in-memory baseline is a map lookup, which is {{< katex >}}O(1){{< /katex >}} on average.
*   **Space Complexity:** {{< katex >}}O(N \cdot P){{< /katex >}}
    *   The agent maintains the full baseline in memory to allow for instant comparisons when a file event is received.

## 5. Portability and Static Compilation

One of the biggest wins I discovered during this project is Go's approach to compilation.

### Targeting Distributions vs. Architectures
I initially wondered if I needed to build separate versions for Ubuntu, CentOS, and Debian. My research taught me that because Go produces **statically linked binaries** (especially with `CGO_ENABLED=0`), the agent includes all the libraries it needs to run. 

This means a single binary built for `linux-amd64` will run on almost any modern Linux distribution without requiring any dependencies (like Python or GLIBC) to be installed on the target server.

### CI/CD Automation
I've implemented a GitHub Actions workflow to automate this process. Every time I push code, the system:
1.  Runs all unit tests to ensure no regressions.
2.  Cross-compiles the binary for both **AMD64** (Standard Servers) and **ARM64** (AWS Graviton / Raspberry Pi).
3.  Uploads the finished binaries as build artifacts.
