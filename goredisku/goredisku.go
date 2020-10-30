package goredisku

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack"
)

type wgAndMtx func(wg *sync.WaitGroup, m *sync.Mutex)
type event func()

var ctx = context.Background()

// RedisCred struct for create new instance of goredisku
type RedisCred struct {
	Address  string
	Password string
	Debug    bool
	DB       int
}

// IGoRedisKu interface for functions
type IGoRedisKu interface {
	Set(key string, value interface{}) (interface{}, error)
	Get(key string, into interface{}) error
	Del(key string) (int64, error)
}

// GoRedisKu struct
type GoRedisKu struct {
	Client       *redis.Client
	GlobalExpire time.Duration
}

// NewGrkClient create client for interacting to redis db
func NewGrkClient(client *redis.Client, globalExpire time.Duration) IGoRedisKu {
	return &GoRedisKu{
		Client:       client,
		GlobalExpire: globalExpire,
	}
}

// Connect to redis server
func (rc RedisCred) Connect() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     rc.Address,
		Password: rc.Password,
		DB:       rc.DB,
	})
}

func (grk *GoRedisKu) store() {
	fmt.Println("Inserting to Cache")
}

func (grk *GoRedisKu) storeConcurrently(wg *sync.WaitGroup, mtx *sync.Mutex) {
	mtx.Lock()
	fmt.Println("Inserting to Cache")
	time.Sleep(1 * time.Second)
	mtx.Unlock()
	wg.Done()
}

// WT is Write Through mechanism for storing cache
func (grk *GoRedisKu) WT(dbInteract event) error {
	var wg sync.WaitGroup
	var mtx sync.Mutex
	wg.Add(1)
	go grk.storeConcurrently(&wg, &mtx)

	wg.Wait()
	dbInteract()
	return nil
}

// WB is Write Back mechanism for storing cache
func (grk *GoRedisKu) WB(dbInteract event) error {
	var wg sync.WaitGroup
	var mtx sync.Mutex
	wg.Add(1)
	mtx.Lock()
	go dbInteract()
	time.Sleep(1 * time.Second)
	mtx.Unlock()
	wg.Done()

	wg.Wait()
	grk.store()
	return nil
}

// Set key value to redis cache db
func (grk *GoRedisKu) Set(key string, value interface{}) (interface{}, error) {
	result, err := msgpack.Marshal(value)
	if err != nil {
		return nil, err
	}
	_, err = grk.Client.SetNX(ctx, key, result, grk.GlobalExpire).Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Get key value from redis cache db
func (grk *GoRedisKu) Get(key string, into interface{}) error {
	fromRedis, err := grk.Client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	err = msgpack.Unmarshal([]byte(fromRedis), &into)
	if err != nil {
		return err
	}
	return nil
}

// Del key value from redis cache db
func (grk *GoRedisKu) Del(key string) (int64, error) {
	return grk.Client.Del(ctx, key).Result()
}
