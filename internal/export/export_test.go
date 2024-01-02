package export

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/waynezhang/foto/internal/cache"
	"github.com/waynezhang/foto/internal/config"
	"github.com/waynezhang/foto/internal/files"
	"github.com/waynezhang/foto/internal/indexer"
	"github.com/waynezhang/foto/internal/testdata"
)

// MockCache

type MockCache struct {
	mock.Mock
}

func (m *MockCache) AddImage(src string, width int, file string) {
	m.Called(src, width, file)
}

func (m *MockCache) CachedImage(src string, width int) *string {
	arg := m.Called(src, width).Get(0)
	if arg == nil {
		return nil
	}
	return arg.(*string)
}

func (m *MockCache) Clear() {
	m.Called()
}

// MockConfig

type MockConfig struct {
	mock.Mock
}

func (m *MockConfig) GetSectionMetadata() []config.SectionMetadata {
	return m.Called().Get(0).([]config.SectionMetadata)
}

func (m *MockConfig) GetOtherFolders() []string {
	return m.Called().Get(0).([]string)
}
func (m *MockConfig) GetExtractOption() config.ExtractOption {
	return m.Called().Get(0).(config.ExtractOption)
}
func (m *MockConfig) AllSettings() map[string]interface{} {
	return m.Called().Get(0).(map[string]interface{})
}

type MockContext struct {
	mock.Mock
}

func (m *MockContext) cleanDirectory(outputPath string) error {
	return m.Called(outputPath).Error(0)
}

func (m *MockContext) buildIndex(cfg config.Config) ([]indexer.Section, error) {
	args := m.Called(cfg)
	var sections []indexer.Section
	var err error
	if args.Get(0) != nil {
		sections = args.Get(0).([]indexer.Section)
	}
	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}
	return sections, err
}

func (m *MockContext) exportPhotos(sections []indexer.Section, outputPath string, cache cache.Cache, progressFunc ProgressFunc) {
	m.Called(sections, outputPath, cache, nil)
}

func (m *MockContext) generateIndexHtml(cfg config.Config, sections []indexer.Section, path string, minimize bool) {
	m.Called(cfg, sections, path, minimize)
}

func (m *MockContext) processOtherFolders(folders []string, outputPath string, minimize bool, messageFunc func(src string, dst string)) {
	m.Called(folders, outputPath, minimize, nil)
}

func TestExportPhotos(t *testing.T) {
	tmp, cache := prepareTempDirAndCache(t)
	defer os.RemoveAll(tmp)

	sections := []indexer.Section{
		{
			Title:     "Section 1",
			Text:      "A description",
			Slug:      "slug-1",
			Folder:    "folder-1",
			Ascending: true,
			ImageSets: []indexer.ImageSet{
				{
					FileName:      "filename-1",
					ThumbnailSize: indexer.ImageSize{Width: 100, Height: 200},
					OriginalSize:  indexer.ImageSize{Width: 300, Height: 400},
				},
				{
					FileName:      "filename-2",
					ThumbnailSize: indexer.ImageSize{Width: 500, Height: 600},
					OriginalSize:  indexer.ImageSize{Width: 700, Height: 800},
				},
			},
		},

		{
			Title:     "Section 2",
			Text:      "A description",
			Slug:      "slug-2",
			Folder:    "folder-2",
			Ascending: true,
			ImageSets: []indexer.ImageSet{
				{
					FileName:      "filename-3",
					ThumbnailSize: indexer.ImageSize{Width: 100, Height: 200},
					OriginalSize:  indexer.ImageSize{Width: 300, Height: 400},
				},
				{
					FileName:      "filename-4",
					ThumbnailSize: indexer.ImageSize{Width: 500, Height: 600},
					OriginalSize:  indexer.ImageSize{Width: 700, Height: 800},
				},
			},
		},
	}

	mockCtx := new(MockContext)
	mockCtx.On("cleanDirectory", mock.Anything).Return(nil)
	mockCtx.On("buildIndex", mock.Anything).Return(sections, nil)
	mockCtx.On("exportPhotos", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	mockCtx.On("generateIndexHtml", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	mockCtx.On("processOtherFolders", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	cfg := new(MockConfig)
	cfg.On("GetOtherFolders").Return([]string{"folder-1", "folder-2"})
	outputPath := "test-directory"
	export(cfg, outputPath, true, cache, mockCtx)

	mockCtx.AssertCalled(t, "cleanDirectory", outputPath)
	mockCtx.AssertCalled(t, "buildIndex", cfg)
	mockCtx.AssertCalled(t, "exportPhotos", sections, filepath.Join(outputPath, "photos"), cache, nil)
	mockCtx.AssertCalled(t, "generateIndexHtml", cfg, sections, filepath.Join(outputPath, "index.html"), true)
	mockCtx.AssertCalled(t, "processOtherFolders", []string{"folder-1", "folder-2"}, outputPath, true, nil)
}

func TestCleanDirectory(t *testing.T) {
	tmp, _ := prepareTempDirAndCache(t)
	defer os.RemoveAll(tmp)

	ctx := DefaultExportContext{}
	ctx.cleanDirectory(tmp)
	assert.False(t, files.IsExisting(tmp))
}

func TestProcessOtherFolders(t *testing.T) {
	tmp, _ := prepareTempDirAndCache(t)
	defer os.RemoveAll(tmp)

	collection1Folder := testdata.Collection1["folder"].(string)
	collection2Folder := testdata.Collection2["folder"].(string)
	new(DefaultExportContext).processOtherFolders([]string{
		collection1Folder,
		collection2Folder,
	}, tmp, true, nil)

	file1 := filepath.Join(tmp, collection1Folder, testdata.Collection1FileName1)
	file2 := filepath.Join(tmp, collection2Folder, testdata.Collection2FileName1)
	assert.True(t, true, files.IsExisting(file1))
	assert.True(t, true, files.IsExisting(file2))
}

func TestResizeImageCache(t *testing.T) {
	tmp, _ := os.MkdirTemp("", "foto-test")
	defer os.RemoveAll(tmp)

	src := testdata.Testfile
	dst := filepath.Join(tmp, "resized.jpg")
	width := testdata.ThumbnailWidth
	cachedFile := testdata.ThumbnailFile

	// non cached
	cache1 := new(MockCache)

	cache1.On("CachedImage", src, width).Return(nil)
	cache1.On("AddImage", src, width, dst).Return(nil)

	err := resizeImageAndCache(src, dst, width, cache1)
	assert.Nil(t, err)
	cache1.AssertCalled(t, "CachedImage", src, width)
	cache1.AssertCalled(t, "AddImage", src, width, dst)

	// cached
	cache2 := new(MockCache)

	cache2.On("CachedImage", src, width).Return(&cachedFile)
	cache2.On("AddImage", src, width, dst).Unset()

	err = resizeImageAndCache(src, dst, width, cache2)
	assert.Nil(t, err)
	cache2.AssertCalled(t, "CachedImage", src, width)
	cache2.AssertNotCalled(t, "AddImage", src, width, dst)
}

// helper func
func prepareTempDirAndCache(t *testing.T) (string, cache.Cache) {
	tmp, err := os.MkdirTemp("", "foto-test")
	assert.Nil(t, err)

	cachePath := filepath.Join(tmp, "cache")
	cache := cache.NewFolderCache(cachePath)
	assert.NotNil(t, cache)

	return tmp, cache
}
