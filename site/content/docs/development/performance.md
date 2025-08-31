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
