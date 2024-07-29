package cache

import (
	"sync"

	"github.com/waynezhang/foto/internal/constants"
)

type Cache interface {
	Migrate()
	AddImage(src string, width int, compressQuality int, file string)
	CachedImage(src string, compressQuality int, width int) *string
	Clear()
}

var (
	once     sync.Once
	instance Cache
)

func Shared() Cache {
	once.Do(func() {
		instance = NewFolderCache(constants.CacheDirectoryName)
		instance.Migrate()
	})
	return instance
}
