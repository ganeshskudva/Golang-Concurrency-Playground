# Rate Limiter for PubSub System

## Overview

This module extends the **Publisher-Subscriber (PubSub)** system by incorporating **rate limiting** to control the rate at which messages are published. The implementation uses a **token bucket algorithm** provided by the `golang.org/x/time/rate` library, ensuring that publishers cannot overwhelm the system, maintaining stability and fairness.

---

## Features

- **Rate Limiting**:
  - Limits the number of messages that can be published per second.
  - Supports a configurable **burst size** for short spikes in publishing.

- **Concurrent Safe**:
  - Handles multiple publishers and subscribers simultaneously without conflicts.

- **Flexible Configuration**:
  - Users can specify the desired rate (`messages/second`) and burst size when initializing the `PubSub` system.

---

## Directory Structure

```plaintext
.
├── cmd
│   ├── ratelimiter                 # Entry point for rate limiter example
│   │   └── main.go
├── pubsub
│   ├── ratelimiter
│   │   ├── rate_limiter.go         # Core implementation for rate limiting
│   │   └── rate_limiter_test.go    # Unit tests for rate limiting
├── Makefile                        # Build and run commands
