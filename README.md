- [Cache Middleware Handler](#cache-middleware-handler)
  * [documentation](#documentation)
    + [storage](#storage)
      - [in memory](#in-memory)
      - [filesystem](#filesystem)
      - [Redis](#redis)
    + [middleware usage](#middleware-usage)
      - [Options](#options)

# Cache Middleware Handler
This package provide a cache middleware handler function for Golang that can be set before `http.HandlerFunc` functions.
The goal is to increase the performance and reduce load on third party e.g. databases by using a cache.
The caching use the *path* and optionally the *method*, *query parameter* values and *header* values as a key for the cache.
Currently the following stores are supported for caching:
- in-memory
- filesystem
- Redis

I tried to keep the configuration as simple as possible.
The following example shows a simple scenario for this package.
A more detailed documentation is placed below.
```go
// ...
import (
	"cache_handler/"
	"cache_handler/store"
)
// ...
func dbQueryHandler(w http.ResponseWriter, r *http.Request) { /* ... */ }
// ...
// Create a store for cached data
store := store.NewInMemoryStore(10 * time.Minute)
// Add the http handler and configuration to the cache middleware
dbQueryCacheMiddlewareHandler := cache_handler.NewMiddleware(
  dbQueryHandler, // http handler
  store, // store to use

  // Optionally add some bypass permissions to allow bypass when
  // - header contains specific field with value
  cache_handler.AllowBypassHeader{Key: "Cache-Status", Value: "bypass"},

  // Optionally define grouping for the cached data
  // e.g.
  // - each Authorization (user) has it's own cache
  cache_handler.UseHeaderKey{Key: "Authorization"},
  // - results are different depending query parameter
  cache_handler.UseQueryParamsKey{Key: "name"},
  cache_handler.UseQueryParamsKey{Key: "department"},
)

http.HandleFunc("/users", dbQueryCacheMiddlewareHandler)
```

## documentation
### storage
The storage is used to store cached data.
#### in memory
This storage type uses the RAM to store cached data.
```go
// NewInMemoryStore(expiration time.Duration)
store := store.NewInMemoryStore(
  // define how long cached data are valid
  1 * time.Second,
)
```
#### filesystem
This storage type uses the filesystem to store cached data.
```go
// NewFilesystem(basePath string, expiration time.Duration)
store := store.NewFilesystem(
  // path to the directory where cache files should stored
  "/tmp/cache",
  // define how long cached data are valid
  1*time.Second
)
```
#### Redis
This storage type uses the Redis to store cached data.
```go
// NewRedisStore(endpoint string, port int, username string, password string, expiration time.Duration)
store := store.NewRedisStore(
  // URL to the redis service
  redisURL,
  // port to the redis service (set to 0 to use default)
  port,
  // username and password if required, if not set an empty string
  "username",
  "password",
  // define how long cached data are valid
  1*time.Second
)
```
### middleware usage
```go
// NewMiddleware(next http.HandlerFunc, store store.Store, opts ...Options)
cacheMiddlewareHandler := cache_handler.NewMiddleware(
	// next http handler func that needs to be cached
  httpHandler,
  // storage to use for cached data
  store,
  // List of options e.g.
  UseHeaderKey{Key: "Authorization"},
  UseQueryParamsKey{Key: "name"},
)
// use it as handler func
http.HandleFunc("/", cacheMiddlewareHandler)
```
#### options
Options are optional configuration parameter.
There are two types of options:
- options to define how to generate the key for caching
- options to define when to bypass the cache

Lets assume an API endpoint `/todo` accept all kinds of HTTP methods and allow `filter` query parameter for *get* methods.
So the cache should only do his job on `get` methods with corresponding `filter`.
In addition, different users have different views of the list.
The users can be differentiated by their `Authorization` header.
It should be also possible to disable caching by setting `Cache-Status` to `bypass` or `dev` on the header.
How to use this package on this case?
1. key options

1.1. `UseQueryParamsKey{Key string}` can be used to set a query parameter as key for the caching.
By setting `filter` as `Key` - `filter=buy` and `filter=prepare` will have their own cached data.
```go
cache_handler.UseQueryParamsKey{Key: "filter"}
```
1.2. `UseHeaderKey{Key string}` can be used to set a header field as key for the caching.
By setting `Authorization` as `Key` - each user (defined by the header field) has their own cached data.
```go
cache_handler.UseHeaderKey{Key: "Authorization"}
```

2. bypass options
2.1. `TestAllowBypassMethod{Key string}` can be user to define methods that should be bypassed.
For our example `post`, `put` and `delete`.
This is not case-sensitive.
```go
cache_handler.TestAllowBypassMethod{Key: "post"}
cache_handler.TestAllowBypassMethod{Key: "PUT"}
cache_handler.TestAllowBypassMethod{Key: "dElEtE"}
```
2.2. `AllowBypassHeader{Key string, Value string}` can be used to define header fields that allows bypass. In our example we want to take a look at `Cache-Status` header field.
```go
cache_handler.AllowBypassHeader{Key: "Cache-Status", Value: "bypass"}
cache_handler.AllowBypassHeader{Key: "Cache-Status", Value: "dev"}
```
