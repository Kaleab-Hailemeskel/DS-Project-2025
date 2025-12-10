package media

type IHLSChunker interface{
	CreateHLSSegments(inputFile, outputDir, segNameTemplate, playlistName string) error
}