# PubSub System with Slow Subscriber Handling

## Overview

This project demonstrates a **Publisher-Subscriber (PubSub)** system in Go with added functionality to handle **slow subscribers**. The system ensures that slow subscribers do not block message delivery to other subscribers or publishers by using buffered channels and gracefully dropping messages when a subscriber's buffer is full.

---

## Features

- **Non-Blocking Delivery**:
  - Messages are delivered concurrently to all subscribers.
  - Slow subscribers are handled without blocking publishers or other subscribers.

- **Graceful Degradation**:
  - Messages are dropped for slow subscribers when their channel buffer is full, ensuring system stability.

- **Logging**:
  - Logs dropped messages to help monitor slow subscriber behavior.

---

## How It Works

1. **Buffered Channels**:
   - Each subscriber is assigned a buffered channel to temporarily hold messages.

2. **Message Dropping**:
   - If a subscriber's buffer is full, new messages are dropped to prevent blocking.

3. **Transparency**:
   - A log message is printed whenever a message is dropped for a slow subscriber.

---

## Usage

### Initialize PubSub

Create a new PubSub instance for a specific message type:

```go
ps := pubsub.NewPubSub[string]()
