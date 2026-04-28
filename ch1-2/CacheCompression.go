package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
	"github.com/spf13/pflag"
)

const (
	StatusActive int = iota
	StatusExpired
	StatusCompressed
)

const (
	MaxScanTokenSize = 64 * 1024
)

type CacheItem struct {
	key           string
	value         []byte
	expires_at    time.Time
	status        int
	is_compressed (*bool)
}

type Cache struct {
	items      []CacheItem
	maxSize    int
	defaultTTL (time.Duration)
}

type CacheStats struct {
	Hits         int
	Misses       int
	Evictions    int
	Compressions int
}

func aliasNormalizeFunc(f *pflag.FlagSet, n string) pflag.NormalizedName {
	switch n {
	case "pass":
		n = "password"
		break
	case "ps":
		n = "password"
		break
	}
	return pflag.NormalizedName(n)
}

func ConsoleReady(name, password string) {
	//flag := viper.BindPFlag()

}

func NewCache(max_Size int, defaultTTL time.Duration) *Cache {
	fmt.Printf("Created cache with maxSize=%d, defaultTTL=%s", max_Size, defaultTTL)
	return &Cache{
		items:      make([]CacheItem, 0, max_Size),
		maxSize:    max_Size,
		defaultTTL: defaultTTL,
	}
}

func (c *Cache) Set(key string, value []byte, ttl time.Duration) error {
	if key == "" || value == nil {
		return errors.New("Error, key or value is empty")
	}
	if ttl <= 0 {
		ttl = c.defaultTTL
	}
	expiresAt := time.Now().Add(ttl)
	for i := 0; i < len(c.items); i++ {
		if c.items[i].key == "" || c.items[i].key == key {
			c.items[i].value = value
			c.items[i].expires_at = expiresAt
			fmt.Printf("Updated: %s = %s\n", key, string(value))
			return nil
		}
	}
	newitem := CacheItem{
		key:        key,
		value:      value,
		expires_at: expiresAt,
	}
	c.items = append(c.items, newitem)
	fmt.Printf("Added %s = %s\n", key, string(value))
	return nil
}

func (c *Cache) Delete(key string) error {
	for i := 0; i < len(c.items); i++ {
		if c.items[i].key == key {
			c.items = append(c.items[i:], c.items[:i+1]...)
			fmt.Printf("Deleted: %s\n", key)
			return nil
		}
	}
	return fmt.Errorf("key '%s' not found", key)
}

func (c *Cache) Get(key string) ([]byte, error) {
	for i := 0; i < len(c.items); i++ {
		if c.items[i].key == key {
			if time.Now().After(c.items[i].expires_at) == true {
				c.items = append(c.items[:i], c.items[i+1:]...)
				return nil, fmt.Errorf("key '%s' expired", key)
			}
			fmt.Printf("Found key: %s = %s\n", key, string(c.items[i].value))
			return c.items[i].value, nil
		}
	}
	return nil, fmt.Errorf("key '%s' not found", key)
}

func (c *Cache) Dump() {
	fmt.Println("\n=== Current cache ===")
	for _, item := range c.items {
		fmt.Printf("  %s = %s (expires at %s)\n",
			item.key,
			string(item.value),
			item.expires_at.Format("15:04:05"))
	}
	fmt.Println("=====================\n")
}

type CachedTime struct {
	userName        string
	compressed_pswd string
	time            time.Time
}

func (ct *CachedTime) ReadFromBigFile(filepath string) ([][]string, error) {
	re_line := regexp.MustCompile(`^\d*`)
	var lines [][]string
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Errorf("Error occured while opening file", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		match := re_line.FindString(string(text))
		if match != "" {
			// Сохраняем как строку и оригинал (например)
			lines = append(lines, []string{match, text})
		} else {
			// Если не совпало, сохраняем пустую строку
			lines = append(lines, []string{"", text})
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines, nil
}

func main() {
	/*
	cache := NewCache(10, 15)

	// СОЗДАЁМ КЛЮЧИ (пользователь сам придумывает)
	cache.Set("user:alice", []byte("Alice's data"), 5*time.Second)
	cache.Set("user:bob", []byte("Bob's data"), 10*time.Second)
	cache.Set("session:123", []byte("session token"), 3*time.Second)

	cache.Dump()

	// ПОЛУЧАЕМ ПО КЛЮЧАМ
	cache.Get("user:alice")
	cache.Get("user:bob")

	// Ждём пока протухнет session:123
	fmt.Println("\nWaiting 4 seconds...")
	time.Sleep(4 * time.Second)

	cache.Get("session:123") // уже протух, удалится сам

	cache.Dump() // session:123 пропал

	// УДАЛЯЕМ ключ
	cache.Delete("user:bob")
	cache.Dump()

	PORT := ":8001"
	arguments := os.Args
	if len(arguments) != 1 {
		PORT = ":" + arguments[1]
	}
	fmt.Println("Using port number: ", PORT)

	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/", myHandler)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	*/

	http.HandleFunc("/gps", handleGPS)
    
    // Показываем где будет файл
    currentDir, _ := os.Getwd()
    log.Printf("🚀 Сервер запущен на :8080")
    log.Printf("📁 CSV файл будет создан в: %s/track_data.csv", currentDir)
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
