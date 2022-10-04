package calendar

import (
	"encoding/base32"
	"errors"
	"fmt"
	"hash"
	"hash/fnv"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

type ErrCacheMiss struct {
	Key    string
	Reason string
}

func (err *ErrCacheMiss) Error() string {
	return fmt.Sprintf("cache miss. reason: %s, key: %s", err.Reason, err.Key)
}

//counterfeiter:generate -o fakes/ . CalendarCacheInterface
type CalendarCacheInterface interface {
	// Get finds a cached calendar entry if
	//
	// The first return value the content of the calendar, if found
	// The second return value is whether a cached entry was found
	// The third return value is any errors
	Get(url string) ([]byte, error)

	// Put stores a calendar content in the cache
	Put(url string, content []byte) error
}

type CalendarCache struct {
	dir string

	HashAlgo hash.Hash32
}

func NewCalendarCache(cacheDir string) *CalendarCache {
	return &CalendarCache{
		dir: cacheDir,
	}
}

// Get finds a cached calendar entry if
// The first return value the content of the calendar, if found
// The second return value is whether a cached entry was found
// The third return value is any errors
func (c *CalendarCache) Get(url string) ([]byte, error) {
	key := c.CacheKey(url)
	cachePath := path.Join(c.dir, key)
	oneHourAgo := time.Now().Add(-1 * time.Hour)

	if finfo, err := os.Stat(cachePath); err == nil {
		if finfo.ModTime().Before(oneHourAgo) {
			return []byte{}, &ErrCacheMiss{
				Key:    key,
				Reason: "expired",
			}
		}

		return ioutil.ReadFile(cachePath)

	} else if errors.Is(err, os.ErrNotExist) {
		return []byte{}, &ErrCacheMiss{
			Key:    key,
			Reason: "not found",
		}

	}

	return []byte{}, &ErrCacheMiss{
		Key:    key,
		Reason: "unknown",
	}
}

// Put stores a calendar content in the cache
func (c *CalendarCache) Put(url string, content []byte) error {
	key := c.CacheKey(url)

	return ioutil.WriteFile(
		path.Join(c.dir, key),
		content,
		0600,
	)
}

func (c *CalendarCache) CacheKey(url string) string {
	hasher := fnv.New64a()
	hasher.Write([]byte(url))
	out := hasher.Sum([]byte{})

	encoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	return strings.ToLower(encoder.EncodeToString(out))
}
