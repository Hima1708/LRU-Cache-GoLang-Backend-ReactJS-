package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type LRUCache struct {
	Head *Node
	Tail *Node
	HT   map[int]*Node
	Cap  int
}

type Node struct {
	Key        int
	Val        int
	Prev       *Node
	Next       *Node
	Timestamp  int64
	Expiration int64
}

func Constructor(capacity int) LRUCache {
	fmt.Printf("Creating LRU Cache with capacity: %d\n", capacity)
	return LRUCache{HT: make(map[int]*Node), Cap: capacity}
}

func (cache *LRUCache) Get(key int) int {
	fmt.Printf("Get request for key: %d\n", key)
	node, ok := cache.HT[key]
	if ok {
		if time.Now().Unix() > node.Expiration {
			fmt.Printf("Key has expired. Removing key: %d\n", key)
			cache.Remove(node)
			delete(cache.HT, key)
			return -1
		}
		fmt.Printf("Key found. Value: %d\n", node.Val)
		node.Timestamp = time.Now().Unix() // Update timestamp on access
		cache.Remove(node)
		cache.Add(node)
		return node.Val
	}
	fmt.Println("Key not found.")
	return -1
}

func (cache *LRUCache) Put(key int, value int, expiration int64) {
	fmt.Printf("Put request - key: %d, value: %d, expiration: %d\n", key, value, expiration)
	// Remove expired items before adding a new one
	cache.removeExpired()

	node, ok := cache.HT[key]
	if ok {
		fmt.Printf("Key exists. Updating value to: %d\n", value)
		node.Val = value
		node.Expiration = time.Now().Unix() + expiration // Update expiration
		cache.Remove(node)
		cache.Add(node)
		return
	} else {
		fmt.Printf("Key does not exist. Adding new key.\n")
		node = &Node{Key: key, Val: value, Timestamp: time.Now().Unix(), Expiration: time.Now().Unix() + expiration}
		cache.HT[key] = node
		cache.Add(node)
	}

	// Remove items if cache exceeds capacity
	if len(cache.HT) > cache.Cap {
		cache.removeLeastRecentlyUsed()
	}
}

func (cache *LRUCache) removeLeastRecentlyUsed() {
	if cache.Tail == nil {
		return
	}
	// Remove the tail (least recently used item)
	oldest := cache.Tail
	cache.Remove(oldest)
	delete(cache.HT, oldest.Key)
	fmt.Printf("Removed least recently used key: %d\n", oldest.Key)
}

func (cache *LRUCache) removeExpired() {
	now := time.Now().Unix()
	for key, node := range cache.HT {
		if now > node.Expiration {
			fmt.Printf("Evicting expired key: %d\n", key)
			cache.Remove(node)
			delete(cache.HT, key)
		}
	}
}

func (cache *LRUCache) Add(node *Node) {
	fmt.Printf("Adding node - key: %d, value: %d\n", node.Key, node.Val)
	node.Prev = nil
	node.Next = cache.Head
	if cache.Head != nil {
		cache.Head.Prev = node
	}
	cache.Head = node
	if cache.Tail == nil {
		cache.Tail = node
	}
	fmt.Printf("Node added - key: %d is now the head\n", node.Key)
}

func (cache *LRUCache) Remove(node *Node) {
	fmt.Printf("Removing node - key: %d, value: %d\n", node.Key, node.Val)
	if node != cache.Head {
		node.Prev.Next = node.Next
	} else {
		cache.Head = node.Next
	}
	if node != cache.Tail {
		node.Next.Prev = node.Prev
	} else {
		cache.Tail = node.Prev
	}
	fmt.Printf("Node removed - key: %d\n", node.Key)
}

var lruCache LRUCache
var clients = make(map[*websocket.Conn]bool)

func getCacheValue(c *gin.Context) {
	key, err := strconv.Atoi(c.Param("key"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid key"})
		return
	}

	value := lruCache.Get(key)
	if value == -1 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Key not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"key": key, "value": value})
}

func getAllCacheValues(c *gin.Context) {
	cacheItems := make(map[int]int)
	for key, node := range lruCache.HT {
		if time.Now().Unix()-node.Timestamp <= 10 {
			cacheItems[key] = node.Val
		} else {
			fmt.Printf("Removing expired key: %d\n", key)
			lruCache.Remove(node)
			delete(lruCache.HT, key)
		}
	}
	c.IndentedJSON(http.StatusOK, cacheItems)
}

func postCacheValue(c *gin.Context) {
	var data struct {
		Key        int   `json:"key"`
		Value      int   `json:"value"`
		Expiration int64 `json:"expiration"` // Expiration in seconds
	}

	if err := c.BindJSON(&data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	lruCache.Put(data.Key, data.Value, data.Expiration)
	broadcastCacheState()
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Key-Value pair added", "key": data.Key, "value": data.Value, "expiration": data.Expiration})
}

func setCacheValue(c *gin.Context) {
	var data struct {
		Key        int   `json:"key"`
		Value      int   `json:"value"`
		Expiration int64 `json:"expiration"` // Expiration in seconds
	}

	if err := c.BindJSON(&data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	lruCache.Put(data.Key, data.Value, data.Expiration)
	broadcastCacheState()
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Key-Value pair updated", "key": data.Key, "value": data.Value, "expiration": data.Expiration})
}

func deleteCacheValue(c *gin.Context) {
	key, err := strconv.Atoi(c.Param("key"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid key"})
		return
	}

	value := lruCache.Get(key)
	if value == -1 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Key not found"})
		return
	}

	lruCache.Put(key, -1, 0) // Remove the key by setting its value to -1 and expiration to 0
	broadcastCacheState()
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Key-Value pair deleted", "key": key})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity; adjust as needed
	},
}

func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("Error upgrading connection: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to upgrade connection"})
		return
	}
	clients[conn] = true

	defer func() {
		delete(clients, conn)
		conn.Close()
	}()

	// Continuously read messages (if needed) or just keep the connection open
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func broadcastCacheState() {
	cacheItems := make(map[int]map[string]interface{})
	now := time.Now().Unix()

	// Collect non-expired cache items
	for key, node := range lruCache.HT {
		if now < node.Expiration {
			cacheItems[key] = map[string]interface{}{
				"value":      node.Val,
				"expiration": node.Expiration - now, // Send remaining expiration time
			}
		} else {
			fmt.Printf("Evicting expired key: %d\n", key)
			lruCache.Remove(node)
			delete(lruCache.HT, key)
		}
	}

	// Prepare message
	message, err := json.Marshal(cacheItems)
	if err != nil {
		fmt.Printf("Error marshalling cache data: %v\n", err)
		return
	}

	// Broadcast the message to all clients
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Printf("Error broadcasting to client: %v\n", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func setCacheCapacity(c *gin.Context) {
	var data struct {
		Capacity int `json:"capacity"`
	}

	if err := c.BindJSON(&data); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	lruCache.Cap = data.Capacity
	fmt.Printf("Cache capacity updated to: %d\n", data.Capacity)

	// Clear expired items after updating capacity
	lruCache.removeExpired()

	// Trim cache to the new capacity
	lruCache.trimToCapacity()

	// Broadcast the updated cache state to reflect the new capacity
	broadcastCacheState()

	c.IndentedJSON(http.StatusOK, gin.H{"capacity": data.Capacity})
}

func (cache *LRUCache) trimToCapacity() {
	fmt.Printf("Trimming cache to capacity: %d\n", cache.Cap)

	// Remove elements if cache exceeds capacity
	for len(cache.HT) > cache.Cap {
		cache.removeLeastRecentlyUsed()
	}
}

func main() {
	lruCache = Constructor(5)
	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Update to match your React app's URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Define Routes
	router.GET("/cache/:key", getCacheValue)
	router.GET("/cache", getAllCacheValues)
	router.POST("/cache", postCacheValue)
	router.PUT("/cache", setCacheValue)
	router.DELETE("/cache/:key", deleteCacheValue)
	router.GET("/ws", handleWebSocket)
	router.PUT("/cache/capacity", setCacheCapacity)

	// Start the server
	go func() {
		for {
			time.Sleep(1 * time.Second)
			broadcastCacheState()
		}
	}()

	router.Run("localhost:8080")
}
