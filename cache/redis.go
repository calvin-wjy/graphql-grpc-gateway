package cache

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	client *redis.Client
}

// // TODO buat struct
// // bikin struct baru bikin mapper
// type Student struct {
// 	name string
// 	nim  string
// 	id   int
// }

func NewRedis() (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:        os.Getenv("REDIS"),
		DB:          0,
		DialTimeout: 100 * time.Millisecond,
		ReadTimeout: 100 * time.Millisecond,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) Get(ctx context.Context, key string) (*redis.StringCmd, error) {
	cmd := c.client.Get(ctx, key)
	// cmdb, _ := cmd.Bytes()

	// b := bytes.NewReader(cmdb)

	// var res *redis.StringCmd
	// gob.NewDecoder(b).Decode(&res)

	// fmt.Println("from get:", &res)
	return cmd, nil
}

func (c *Client) Set(ctx context.Context, key string, value interface{}) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, p, 300*time.Second).Err()

	// var b bytes.Buffer
	// gob.NewEncoder(&b).Encode(value)

	// return c.client.Set(ctx, key, string(b.Bytes()), 300*time.Second).Err()

}
