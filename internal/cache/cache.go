package cache

import (
	"github.com/waynezhang/foto/internal/constants"
)

type Cache interface {
	AddImage(src string, width int, file string)
	CachedImage(src string, width int) *string
	Clear()
}

func Shared() Cache {
	once.Do(func() {
		instance = NewFolderCache(constants.CacheDirectoryName)
	})
	return instance
}
