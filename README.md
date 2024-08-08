# LRU-Cache-App
## LRU-Cache Backend
This Go application provides an LRU (Least Recently Used) cache backend. It includes features for managing cache entries and real-time updates via WebSocket.

## Key Features
LRU Cache: Automatically evicts the least recently used items when the cache is full.

## API Endpoints:
- GET /cache/:key - Retrieve a cache entry by key.
- GET /cache - Get all cache entries.
- POST /cache - Add a new cache entry.
- PUT /cache - Update an existing cache entry.
- DELETE /cache/:key - Remove a cache entry.
- PUT /cache/capacity - Change the cache capacity.
- WebSocket: Real-time updates on cache changes.

## LRU-Cache Frontend
This project provides a React-based frontend for interacting with a cache management system. It includes functionalities to view, add, and retrieve cache entries, and set the cache capacity.

## Components
### App.js:
- The main component that orchestrates the entire application.
- Manages state for cache data, capacity, and modal visibility.
- Establishes a WebSocket connection to receive real-time updates from the backend.

### CapacityModal.js:
- Modal for setting the cache capacity.
- Allows users to input a new cache capacity and submit the change.

### CacheTable.js:
- Displays the current cache data in a table format.
- Shows key-value pairs and their expiration times.

### AddEntryModal.js:
- Modal for adding new cache entries.
- Allows users to input key, value, and expiration time for new cache entries.

### GetEntryModal.js:
- Modal for retrieving a specific cache entry by key.
- Users can enter a key to fetch the associated cache value.
