package export

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/waynezhang/foto/internal/cache"
	"github.com/waynezhang/foto/internal/config"
	"github.com/waynezhang/foto/internal/constants"
	"github.com/waynezhang/foto/internal/files"
	"github.com/waynezhang/foto/internal/indexer"
	mm "github.com/waynezhang/foto/internal/minimize"
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

// MockContext

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

func (m *MockContext) exportPhotos(sections []indexer.Section, outputPath string, cache cache.Cache, progressFn progressFunc) {
	m.Called(sections, outputPath, cache, nil)
}

func (m *MockContext) generateIndexHtml(cfg config.Config, templatePath string, sections []indexer.Section, path string, minimizer mm.Minimizer) {
	m.Called(cfg, templatePath, sections, path, minimizer)
}

func (m *MockContext) processOtherFolders(folders []string, outputPath string, minimizer mm.Minimizer, messageFunc func(src string, dst string)) {
	m.Called(folders, outputPath, minimizer, nil)
}

// MockFunc
type MockFunc struct {
	mock.Mock
}

func (m *MockFunc) progressFunc(path string) {
	m.Called(path)
}

func (m *MockFunc) messageFunc(src string, dst string) {
	m.Called(src, dst)
}

// MockMinimizer

type MockMinimizer struct {
	mock.Mock
}

func (m *MockMinimizer) Minimizable(path string) bool {
	return m.Called(path).Bool(0)
}

func (m *MockMinimizer) MinimizeFile(src string, dest string) error {
	arg := m.Called(src, dest).Get(0)
	if arg != nil {
		return arg.(error)
	}
	return nil
}

// Tests

func TestExport(t *testing.T) {
	tmp, cache := prepareTempDirAndCache(t)
	defer os.RemoveAll(tmp)

	var section1 indexer.Section
	var section2 indexer.Section
	_ = mapstructure.Decode(testdata.Collection1, &section1)
	_ = mapstructure.Decode(testdata.Collection1, &section2)
	sections := []indexer.Section{section1, section2}

	mockCtx := new(MockContext)
	mockCtx.On("cleanDirectory", mock.Anything).Return(nil)
	mockCtx.On("buildIndex", mock.Anything).Return(sections, nil)
	mockCtx.On("exportPhotos", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	mockCtx.On("generateIndexHtml", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	mockCtx.On("processOtherFolders", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	minimizer := mm.NoneMinimizer{}

	cfg := new(MockConfig)
	cfg.On("GetOtherFolders").Return([]string{"folder-1", "folder-2"})
	outputPath := "test-directory"
	export(cfg, outputPath, minimizer, cache, mockCtx)

	mockCtx.AssertCalled(t, "cleanDirectory", outputPath)
	mockCtx.AssertCalled(t, "buildIndex", cfg)
	mockCtx.AssertCalled(t, "exportPhotos", sections, filepath.Join(outputPath, "photos"), cache, nil)
	mockCtx.AssertCalled(t, "generateIndexHtml", cfg, constants.TemplateFilePath, sections, filepath.Join(outputPath, "index.html"), minimizer)
	mockCtx.AssertCalled(t, "processOtherFolders", []string{"folder-1", "folder-2"}, outputPath, minimizer, nil)
}

func TestCleanDirectory(t *testing.T) {
	tmp, _ := prepareTempDirAndCache(t)
	defer os.RemoveAll(tmp)

	ctx := defaultExportContext{}
	_ = ctx.cleanDirectory(tmp)
	assert.False(t, files.IsExisting(tmp))
}

func TestExportPhotos(t *testing.T) {
	tmp, cache := prepareTempDirAndCache(t)
	defer os.RemoveAll(tmp)

	var section1 indexer.Section
	var section2 indexer.Section
	_ = mapstructure.Decode(testdata.Collection1, &section1)
	_ = mapstructure.Decode(testdata.Collection1, &section2)
	sections := []indexer.Section{section1, section2}

	mockFunc := new(MockFunc)
	mockFunc.On("progressFunc", mock.Anything).Return()
	progressFunc := func(path string) {
		mockFunc.progressFunc(path)
	}

	ctx := defaultExportContext{}
	ctx.exportPhotos(sections, tmp, cache, progressFunc)

	for _, s := range sections {
		assert.True(t, files.IsExisting(filepath.Join(tmp, s.Slug)))
		for _, set := range s.ImageSets {
			expectedOriginalPath := filepath.Join(tmp, s.Slug, "original", set.FileName)
			assert.Truef(t, files.IsExisting(expectedOriginalPath), expectedOriginalPath)

			expectedThumbnailPath := filepath.Join(tmp, s.Slug, "thubmnail", set.FileName)
			assert.Truef(t, files.IsExisting(expectedOriginalPath), expectedThumbnailPath)
		}
	}
	mockFunc.AssertNumberOfCalls(t, "progressFunc", 6) // 6 files
}

func TestGenerateIndexHTML(t *testing.T) {
	tmp, _ := prepareTempDirAndCache(t)
	defer os.RemoveAll(tmp)

	sections := []indexer.Section{}
	path := filepath.Join(tmp, "index.html")

	cfg := MockConfig{}
	cfg.On("AllSettings").Return(map[string]interface{}{})

	mockMinimizer := new(MockMinimizer)
	mockMinimizer.On("MinimizeFile", mock.Anything, mock.Anything).Return(nil)

	ctx := defaultExportContext{}
	ctx.generateIndexHtml(&cfg, testdata.TestHtmlFile, sections, path, mockMinimizer)
	assert.True(t, files.IsExisting(path))
	cfg.AssertCalled(t, "AllSettings")

	mockMinimizer.AssertCalled(t, "MinimizeFile", mock.Anything, mock.Anything)
}

func TestProcessOtherFolders(t *testing.T) {
	tmp, _ := prepareTempDirAndCache(t)
	tmp = "/tmp/a"
	files.EnsureDirectory(tmp)
	defer os.RemoveAll(tmp)

	mockMinimizer := new(MockMinimizer)
	mockMinimizer.On("Minimizable", mock.Anything).Return(true)
	mockMinimizer.On("MinimizeFile", mock.Anything, mock.Anything).Return(nil)

	mockFunc := new(MockFunc)
	mockFunc.On("messageFunc", mock.Anything, mock.Anything).Return(nil)
	messageFunc := func(src string, dst string) {
		mockFunc.messageFunc(src, dst)
	}

	collection1Folder := testdata.Collection1["folder"].(string)
	collection2Folder := testdata.Collection2["folder"].(string)
	new(defaultExportContext).processOtherFolders([]string{
		collection1Folder,
		collection2Folder,
	}, tmp, mockMinimizer, messageFunc)

	file1 := filepath.Join(tmp, filepath.Base(collection1Folder), testdata.Collection1FileName1)
	file2 := filepath.Join(tmp, filepath.Base(collection2Folder), testdata.Collection2FileName1)
	assert.True(t, files.IsExisting(file1))
	assert.True(t, files.IsExisting(file2))

	mockFunc.AssertNumberOfCalls(t, "messageFunc", 2) // 2 folders

	mockMinimizer.AssertCalled(t, "Minimizable", mock.Anything)
	mockMinimizer.AssertCalled(t, "MinimizeFile", mock.Anything, mock.Anything)
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

func TestMinimizer(t *testing.T) {
	assert.Equal(t, reflect.TypeOf(mm.NoneMinimizer{}), reflect.TypeOf(minimizer(false)))
	assert.Equal(t, reflect.TypeOf(mm.MinifyMinimizer{}), reflect.TypeOf(minimizer(true)))
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
