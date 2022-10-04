package calendar_test

import (
	"io/ioutil"
	"os"
	"path"
	"time"

	. "github.com/AP-Hunt/what-next/m/calendar"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CalendarCache", func() {
	var (
		cacheDir string
		cache    *CalendarCache
	)

	BeforeEach(func() {
		tempDir, err := os.MkdirTemp(os.TempDir(), "what-next_cache_unit_test_*")
		Expect(err).ToNot(HaveOccurred())
		cacheDir = tempDir

		cache = NewCalendarCache(cacheDir)
	})

	AfterEach(func() {
		os.RemoveAll(cacheDir)
	})

	Describe("Put", func() {
		It("writes the given content to a file named for the hash of the URL", func() {
			url := "https://example.com"
			content := "BEGIN:VCALENDAR END:VCALENDAR"

			expectedCacheKey := cache.CacheKey(url)
			expectedCachePath := path.Join(cacheDir, string(expectedCacheKey))

			err := cache.Put(url, []byte(content))
			Expect(err).ToNot(HaveOccurred())
			Expect(expectedCachePath).To(BeAnExistingFile())

			fileContent, err := ioutil.ReadFile(expectedCachePath)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(fileContent)).To(Equal(content))
		})
	})

	Describe("Get", func() {
		Context("on a cache hit", func() {
			It("returns the content of the cached calendar", func() {
				url := "https://example.com"
				content := "BEGIN:VCALENDAR END:VCALENDAR"
				key := cache.CacheKey(url)

				err := ioutil.WriteFile(
					path.Join(cacheDir, key),
					[]byte(content),
					0600,
				)
				Expect(err).ToNot(HaveOccurred())

				cacheContent, err := cache.Get(url)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(cacheContent)).To(Equal(content))
			})
		})

		Context("on a cache miss because the key doesn't exist", func() {
			It("will return an ErrCacheMiss with a reason of 'not found'", func() {
				_, err := cache.Get("unknown")
				Expect(err).To(BeAssignableToTypeOf(&ErrCacheMiss{}))

				cacheMiss := err.(*ErrCacheMiss)
				Expect(cacheMiss.Reason).To(Equal("not found"))
			})
		})

		Context("on a cache miss because the key expired", func() {
			It("will return an ErrCacheMiss with a reason of 'expired'", func() {
				url := "https://example.com"
				content := "BEGIN:VCALENDAR END:VCALENDAR"
				key := cache.CacheKey(url)

				cachePath := path.Join(cacheDir, key)
				err := ioutil.WriteFile(
					cachePath,
					[]byte(content),
					0600,
				)
				Expect(err).ToNot(HaveOccurred())

				// Move mod time back a few hours
				os.Chtimes(
					cachePath,
					time.Now(),
					time.Now().Add(-6*time.Hour),
				)

				_, err = cache.Get(url)
				Expect(err).To(BeAssignableToTypeOf(&ErrCacheMiss{}))

				cacheMiss := err.(*ErrCacheMiss)
				Expect(cacheMiss.Reason).To(Equal("expired"))
			})
		})
	})
})
