package media

type IMediaProcessor interface {
	CreateHLSSegments(inputFile, outputDir, segNameTemplate, playlistName string) error
	ExtractAlbumArt(inputFile string) ([]byte, error)
}
