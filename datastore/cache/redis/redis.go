package redis

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	redigo "github.com/gomodule/redigo/redis"
	rejson "github.com/nitishm/go-rejson/v4"
	redis "github.com/redis/go-redis/v9"

	cache "github.com/gnanasuryateja/golib/datastore/cache"
)

type RedisStoreConfig struct {
	Addr     string
	Port     string
	Username string
	Password string
	CA       *string
	CRT      *string
	Key      *string
	DB       *int
}

// validates the input params
func (rsc RedisStoreConfig) validate() error {
	if rsc.Addr == "" || rsc.Port == "" || rsc.Username == "" || rsc.Password == "" {
		return fmt.Errorf("RedisStoreConfig cannot have empty values")
	}
	if rsc.CA != nil && *rsc.CA == "" {
		return fmt.Errorf("CA cannot be empty value")
	}
	if rsc.CRT != nil && *rsc.CRT == "" {
		return fmt.Errorf("CRT cannot be empty value")
	}
	if rsc.Key != nil && *rsc.Key == "" {
		return fmt.Errorf("key cannot be empty value")
	}
	return nil
}

const (
	redis_ping_str                    = "ping: PONG"
	redis_default_db                  = 0
	redis_add_success_acknowledgement = "Sucessfully added to redis...:)"
)

type redisStore struct {
	client        *redis.Client
	rejsonHandler *rejson.Handler
}

// creates a new redisStore client
func NewRedisStoreClient(ctx context.Context, redisStoreConfig RedisStoreConfig) (cache.Cache, error) {

	// validate the redisStoreConfig
	err := redisStoreConfig.validate()
	if err != nil {
		return nil, err
	}

	redisUri := fmt.Sprintf("%v:%v", redisStoreConfig.Addr, redisStoreConfig.Port)
	redisOptions := redis.Options{
		Addr:     redisUri,
		Username: redisStoreConfig.Username,
		Password: redisStoreConfig.Password,
	}
	if redisStoreConfig.DB != nil {
		redisOptions.DB = *redisStoreConfig.DB
	} else {
		redisOptions.DB = redis_default_db
	}

	if redisStoreConfig.CA != nil && *redisStoreConfig.CA != "" && redisStoreConfig.CRT != nil && *redisStoreConfig.CRT != "" && redisStoreConfig.Key != nil && *redisStoreConfig.Key != "" {

		// Load client certificate and key
		cert, err := tls.LoadX509KeyPair(*redisStoreConfig.CRT, *redisStoreConfig.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate and key: %v", err)
		}

		// create a tls config
		tlsConfig := tls.Config{
			Certificates: []tls.Certificate{
				cert,
			},
		}

		// read the CA certificate
		caCert, err := os.ReadFile(*redisStoreConfig.CA)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
		redisOptions.TLSConfig = &tlsConfig
	}

	// create a new redis client
	client := redis.NewClient(&redisOptions)

	// check if the connection is successful
	ping := client.Ping(ctx)
	if ping.String() != redis_ping_str {
		return nil, ping.Err()
	}
	rjhandler := rejson.NewReJSONHandler()
	rjhandler.SetGoRedisClientWithContext(ctx, client)
	return &redisStore{
		rejsonHandler: rjhandler,
		client:        client,
	}, nil
}

// checks the connection to cache and returns error if any
func (rs *redisStore) HealthCheck(ctx context.Context) error {
	// check the connection
	ping := rs.client.Ping(ctx)
	if ping.String() == redis_ping_str {
		return nil
	}
	return fmt.Errorf("HealthCheck failed for Redis: %v", ping.Err())
}

// inserts data into cache
func (rs *redisStore) AddData(ctx context.Context, args ...any) (string, error) {

	// validate the passed args
	if len(args) < 2 {
		return "", fmt.Errorf("collection key or value is(are) missing")
	}
	if len(args) > 2 {
		return "", fmt.Errorf("more params are passed than expected")
	}

	// extract the key from args
	key, ok := args[0].(string)
	if !ok {
		return "", fmt.Errorf("invalid key is passed (not a string)")
	}

	// extract the value from args
	value := args[1]

	// add the key value to redis
	_, err := rs.rejsonHandler.JSONSet(key, ".", value)
	if err != nil {
		return "", err
	}

	// return the acknowledgement
	return redis_add_success_acknowledgement, nil
}

// gets the data from cache
func (rs *redisStore) GetData(ctx context.Context, args ...any) (any, error) {

	// validate the passed args
	if len(args) < 1 {
		return "", fmt.Errorf("collection key or value is(are) missing")
	}
	if len(args) > 1 {
		return "", fmt.Errorf("more params are passed than expected")
	}

	// extract the key from args
	key, ok := args[0].(string)
	if !ok {
		return "", fmt.Errorf("invalid key is passed (not a string)")
	}

	// return the data
	return redigo.Bytes(rs.rejsonHandler.JSONGet(key, "."))
}

// gets all the keys from cache
func (rs *redisStore) GetKeys(ctx context.Context, pattern string) ([]string, error) {
	if pattern == "" {
		var keys []string
		var cursor uint64 = 0
		for {
			var result []string
			var err error
			result, cursor, err = rs.client.Scan(ctx, cursor, "*", 10).Result()
			if err != nil {
				fmt.Println("error scanning keys:", err)
				return nil, err
			}
			keys = append(keys, result...)
			if cursor == 0 {
				break
			}
		}
		return keys, nil
	}
	return rs.client.Keys(ctx, pattern).Result()
}

// deletes the data from cache
func (rs *redisStore) DeleteData(ctx context.Context, args ...any) (any, error) {

	// validate the passed args
	if len(args) < 1 {
		return "", fmt.Errorf("collection key or value is(are) missing")
	}
	if len(args) > 1 {
		return "", fmt.Errorf("more params are passed than expected")
	}

	// extract the key from args
	key, ok := args[0].(string)
	if !ok {
		return "", fmt.Errorf("invalid key is passed (not a string)")
	}

	// return the response
	return rs.rejsonHandler.JSONDel(key, ".")
}
