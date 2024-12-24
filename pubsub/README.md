# Basic PubSub System

## Overview

This project implements a basic **Publisher-Subscriber (PubSub)** system in Go using generics. It allows publishers to send messages to subscribers through specific topics. The design focuses on concurrency, scalability, and reliability.

---

## Features

- **Generic Implementation**:
  - Supports any data type for messages using Go generics.
  
- **Concurrency Safe**:
  - Manages multiple publishers and subscribers safely using `sync.RWMutex`.

- **High Throughput**:
  - Uses buffered channels to handle bursts of messages without blocking.

- **Graceful Shutdown**:
  - Cleans up resources and closes all active subscriber channels.

---

## How It Works

1. **Topics**:
   - Messages are organized by topics.
   - Each topic can have multiple subscribers.

2. **Publish**:
   - Publishers send messages to a topic.
   - Messages are delivered concurrently to all subscribers of the topic.

3. **Subscribe**:
   - Subscribers listen for messages on a specific topic.
   - Each subscriber receives messages through a buffered channel.

4. **Unsubscribe**:
   - Subscribers can unsubscribe from a topic, stopping message delivery and releasing resources.

5. **Shutdown**:
   - Closes all subscriber channels and cleans up the `PubSub` instance.

---

## Directory Structure

```plaintext
.
├── pubsub
│   └── pubsub.go              # Core implementation of the PubSub system
