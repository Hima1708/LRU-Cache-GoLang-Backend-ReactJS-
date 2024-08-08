# LRU-Cache-App
##LRU-Cache Backend
This Go application provides an LRU (Least Recently Used) cache backend. It includes features for managing cache entries and real-time updates via WebSocket.

## Key Features
LRU Cache: Automatically evicts the least recently used items when the cache is full.

## API Endpoints:
GET /cache/:key - Retrieve a cache entry by key.
GET /cache - Get all cache entries.
POST /cache - Add a new cache entry.
PUT /cache - Update an existing cache entry.
DELETE /cache/:key - Remove a cache entry.
PUT /cache/capacity - Change the cache capacity.
WebSocket: Real-time updates on cache changes.
