package redisPool

import (
	"github.com/go-redis/redis"
	"log"
	"math/rand"
	"os"
	"strconv"
)

type Pool struct {
	host       string
	pool       []*redis.Client
	maxEntries int
}

var pool = Pool{
	host:       "",
	maxEntries: 10,
	pool:       []*redis.Client{},
}

func (rp *Pool) InitPool() {
	rp.host = os.Getenv("REDISDB")
	if rp.host == "" {
		rp.host = "127.0.0.1"
	}
	for i := 0; i < rp.maxEntries; i++ {
		rp.pool = append(rp.pool, rp.connect())
	}
}

func (rp *Pool) GetClient() (*redis.Client, int) {
	entry := rand.Intn(rp.maxEntries)
	return rp.pool[entry], entry
}

func (rp *Pool) connect() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     rp.host + ":6379", // domain is redis normally
		Password: "",
		DB:       0,
	})
	return rdb
}

func (rp *Pool) Reconnect(entry int) {
	istr := strconv.Itoa(entry)
	log.Println("Reconnecting client " + istr)
	rp.pool[entry] = rp.connect()
	// Todo: investigate race condition with the other retries.
	/*go func() {
		for i := 0; i < rp.maxEntries; i++ {
			ping := rp.pool[entry].Ping()
			if ping.Err() != nil {
				istr = strconv.Itoa(i)
				log.Println("Reconnecting client " + istr)
				rp.pool[entry] = rp.connect()
			}
		}
	}()*/
}

func GetPool() *Pool {
	return &pool
}
