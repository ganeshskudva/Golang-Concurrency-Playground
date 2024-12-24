# Deadlock Prevention in a PubSub System

## Overview

This repository demonstrates a **Publisher-Subscriber (PubSub)** system implemented in Go, with an emphasis on **deadlock prevention**. The design ensures safe and efficient handling of shared resources, such as the `subscribers` map, while enabling concurrent publishers and subscribers to interact without blocking or causing system-wide deadlocks.

---

## Features

- **Concurrent Safe PubSub**: Supports multiple publishers and subscribers interacting simultaneously on shared topics.
- **Deadlock Prevention**: Uses strategies like snapshotting subscribers and releasing locks early to avoid deadlocks.
- **Graceful Shutdown**: Closes all channels and releases resources cleanly, avoiding lingering goroutines.
- **High Throughput**: Buffered channels handle bursts of messages efficiently.

---

## Deadlock Prevention Techniques

1. **Minimized Lock Scope**:
   - Locks are held only for critical operations (e.g., accessing or modifying the `subscribers` map).
   - Message delivery happens outside locked regions.

2. **Snapshotting Subscribers**:
   - A snapshot of subscriber channels is taken before delivering messages, reducing lock contention.

3. **Avoid Blocking Inside Locks**:
   - Operations like `channel <- message` are performed outside of critical sections.

4. **No Nested Locks**:
   - Functions holding a lock do not invoke other functions that also acquire locks.

5. **Graceful Shutdown**:
   - Ensures all channels are closed properly and signals all goroutines to exit without deadlocks.

---

## Directory Structure

```plaintext
.
├── cmd
│   ├── pubsub                 # Entry point for the PubSub example
│   │   └── main.go
│   ├── deadlockprevention     # Entry point for Deadlock Prevention example
│   │   └── main.go
├── pubsub
│   ├── pubsub.go              # Core PubSub implementation
│   └── deadlockprevention
│       └── deadlock_prevention.go # Deadlock prevention focused implementation
├── Makefile                   # Build and run commands
