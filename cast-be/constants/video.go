package constants

const (
	VideoTypeVOD     = "vod"
	VideoTypeLive    = "live"
	ThumbnailRootDir = "thumbnail"
	ThumbnailDefault = "_default"
	ThumbnailWidth   = 720
	ThumbnailHeight  = 405

	VideoListTrending   = "trending"
	VideoListSubscribed = "subscribed"
	VideoListLiked      = "liked"
)

var (
	VideoResolutions = []string{"Processing", "240p", "360p", "480p", "720p", "1080p"}
)
