package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirectoryManipulating(t *testing.T) {
	tmp, err := os.MkdirTemp("", "foto-test")
	assert.Nil(t, err)
	defer os.RemoveAll(tmp)

	assert.True(t, IsExisting(tmp))

	dir1 := filepath.Join(tmp, "parent1", "child1")
	_ = EnsureDirectory(dir1)
	assert.True(t, IsExisting(dir1))

	// no failure to create existing directory
	_ = EnsureParentDirectory(dir1)
	assert.True(t, IsExisting(dir1))

	dir2Parent := filepath.Join(tmp, "parent2")
	dir2 := filepath.Join(dir2Parent, "child2")
	_ = EnsureParentDirectory(dir2)
	assert.True(t, IsExisting(dir2Parent))
	assert.False(t, IsExisting(dir2))

	_ = PruneDirectory(tmp)
	assert.False(t, IsExisting(tmp))
}

func TestPath(t *testing.T) {
	assert.Equal(t, "base_path/index.html", filepath.ToSlash(OutputIndexFilePath("base_path")))
	assert.Equal(t, "base_path/photos", filepath.ToSlash(OutputPhotosFilePath("base_path")))

	photoFilePath := filepath.Join("some-directory", "photo.jpg")
	originalPath := OutputPhotoOriginalFilePath("base_path", "a-slug", photoFilePath)
	assert.Equal(t, "base_path/a-slug/original/photo.jpg", filepath.ToSlash(originalPath))
	thumbnailPath := OutputPhotoThumbnailFilePath("base_path", "a-slug", photoFilePath)
	assert.Equal(t, "base_path/a-slug/thumbnail/photo.jpg", filepath.ToSlash(thumbnailPath))
}
