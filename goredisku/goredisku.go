package goredisku

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack"
)

type wgAndMtx func(
	wg *sync.WaitGroup,
	mtx *sync.Mutex,
)

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

func (grk *GoRedisKu) storeSync(
	key string,
	value interface{},
	wg *sync.WaitGroup,
	mtx *sync.Mutex,
) error {
	fmt.Println("Doing cache insertion")
	fatalErrors := make(chan error)
	mtx.Lock()
	go func() {
		// delete existing cache
		_, err := grk.Client.Del(ctx, key).Result()
		if err != nil {
			fatalErrors <- err
		}

		// create new cache
		result, err := msgpack.Marshal(value)
		if err != nil {
			fatalErrors <- err
		}
		_, err = grk.Client.SetNX(ctx, key, result, grk.GlobalExpire).Result()
		if err != nil {
			fatalErrors <- err
		}
	}()
	mtx.Unlock()
	wg.Done()
	close(fatalErrors)
	return <-fatalErrors
}

func (grk *GoRedisKu) store(
	key string,
	value interface{},
) error {
	fmt.Println("Doing cache insertion")
	// delete existing cache
	_, err := grk.Client.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	// create new cache
	result, err := msgpack.Marshal(value)
	if err != nil {
		return err
	}
	_, err = grk.Client.SetNX(ctx, key, result, grk.GlobalExpire).Result()
	if err != nil {
		return err
	}
	return nil
}

// WT is Write Through mechanism for storing cache
func (grk *GoRedisKu) WT(
	key string,
	value interface{},
	dbInteract event,
) error {
	var wg sync.WaitGroup
	var mtx sync.Mutex
	wg.Add(1)

	err := grk.storeSync(key, value, &wg, &mtx)
	if err != nil {
		return err
	}

	wg.Wait()
	dbInteract()
	return nil
}

// WB is Write Back mechanism for storing cache
func (grk *GoRedisKu) WB(
	key string,
	value interface{},
	dbInteract wgAndMtx,
) error {
	var wg sync.WaitGroup
	var mtx sync.Mutex
	wg.Add(1)

	go dbInteract(&wg, &mtx)

	wg.Wait()
	err := grk.store(key, value)
	if err != nil {
		return err
	}
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
