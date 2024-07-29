package testdata

const (
	TestConfigFile   = "../../fs/static/foto.toml"
	TestConfigFileV1 = "../../testdata/foto_v1.toml"
	TestConfigFileV2 = "../../testdata/foto_v2.toml"

	TestHtmlFile       = "../../fs/static/templates/template.html"
	TestCssFile        = "../../fs/static/assets/style.css"
	TestTxtFile        = "../../testdata/test.txt"
	TestJavascriptFile = "../../testdata/test.js"
)

const (
	Testfile         = "../../testdata/collection-1/2022-06-29.jpg"
	TestfileWidth    = 1440
	TestfileHeight   = 1080
	ExpectedChecksum = "2786728c2c9eb5334df492e1853e24c72f976e063ebd513b45bc47476178cc23"

	OriginalWidth            = 2048
	OriginalHeight           = 1536
	OriginalFile             = "../../testdata/testfile-2048.jpg"
	ExpectedOriginalChecksum = "6511f6cde2282593ac0bd643ba061413ff720bf219c42c21745497d90a5da9c8"

	ThumbnailWidth            = 640
	ThumbnailHeight           = 480
	ThumbnailFile             = "../../testdata/testfile-640.jpg"
	ExpectedThubmnailChecksum = "1c8a6195eefb53be554d86df9de1ae7c5559fa71938be1db595c3bef6c063796"

	CompressQuality = 75

	CompressQualityHQ           = 100
	ExpectedThubmnailHQChecksum = "2ff919b7f866c18ca835f37c6f550f80b4d39397a172e6fffab44ad93c9023ca"
)

var (
	Collection1 = map[string]interface{}{
		"title":     "Section 1",
		"text":      "This is Section 1",
		"slug":      "slug-section-1",
		"folder":    "../../testdata/collection-1",
		"ascending": true,
		"imageSets": []map[string]interface{}{
			{
				"fileName":      "2022-06-29.jpg",
				"thumbnailSize": 640,
				"originalSize":  2048,
			},
			{
				"fileName":      "2022-07-01.jpg",
				"thumbnailSize": 640,
				"originalSize":  2048,
			},
			{
				"fileName":      "2022-07-19.jpg",
				"thumbnailSize": 640,
				"originalSize":  2048,
			},
		},
	}
	Collection1FileName1          = "2022-06-29.jpg"
	Collection1ThumbnailChecksum1 = "1c8a6195eefb53be554d86df9de1ae7c5559fa71938be1db595c3bef6c063796"
	Collection1OriginalChecksum1  = "6511f6cde2282593ac0bd643ba061413ff720bf219c42c21745497d90a5da9c8"

	Collection1FileName2          = "2022-07-01.jpg"
	Collection1ThumbnailChecksum2 = "c32cac109360865112b865100e3c47a02795efe7f1df32db7bbc9329db06173c"
	Collection1OriginalChecksum2  = "52a811492f2cd6f0463112f2dfb406bdd5b9f996eb8864005b1f46e8057d545f"

	Collection1FileName3          = "2022-07-19.jpg"
	Collection1ThumbnailChecksum3 = "d55eeeb2e0b1f91d98c8d22cefafaf061aba26f79402dc7c95c68869049ba643"
	Collection1OriginalChecksum3  = "10ea1324b019eeb47a8ba608c470c0e55df6aa95076c752ea413fcf9060180da"
)

var (
	Collection2 = map[string]interface{}{
		"title":     "Section 2",
		"text":      "This is Section 2",
		"slug":      "slug-section-2",
		"folder":    "../../testdata/collection-2",
		"ascending": false,
		"imageSets": []map[string]interface{}{
			{
				"fileName":      "2022-09-28.jpg",
				"thumbnailSize": 640,
				"originalSize":  2048,
			},
			{
				"fileName":      "2022-09-20.jpg",
				"thumbnailSize": 640,
				"originalSize":  2048,
			},
			{
				"fileName":      "2022-04-29.jpg",
				"thumbnailSize": 640,
				"originalSize":  2048,
			},
		},
	}

	Collection2FileName1 = "2023-09-28.jpg"
	Collection2FileName2 = "2022-09-20.jpg"
	Collection2FileName3 = "2023-04-29.jpg"
)

var (
	EmptyCollection = map[string]interface{}{
		"title":     "Empty Collection",
		"text":      "This is an empty collection",
		"slug":      "empty-slug-section",
		"folder":    "../../testdata/empty-section",
		"ascending": true,
		"imageSets": []map[string]interface{}{},
	}
)

var (
	RotatedImageFile            = "../../testdata/rotated.jpg"
	RotatedImageWidth           = 1440
	RotatedImageHeight          = 1080
	RotatedImageThumbnailWidth  = 640
	RotatedImageThumbnailHeight = 480
)

var (
	WebpTestFile       = "../../testdata/webp/test.webp"
	WebpTestfileWidth  = 8000
	WebpTestfileHeight = 6000

	WebpThumbnailWidth            = 640
	WebpThumbnailHeight           = 480
	WebpExpectedThubmnailChecksum = "5a5f8fcae2e37e504d6e062ee8adefac45ec9102a410faa16d4a0278b310b0c8"
)
