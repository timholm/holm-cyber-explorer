package types

// Interfaces for dependency injection and testing

type RedisClient interface {
	Get(key string) (string, error)
	Set(key string, value interface{}) error
	Del(key string) error
	LPush(key string, values ...interface{}) error
	RPop(key string) (string, error)
	Close() error
}

type YouTubeClient interface {
	GetChannelInfo(url string) (*Channel, error)
	GetVideoList(channelID string) ([]Video, error)
}

type StorageManager interface {
	GetChannelPath(channelID string) string
	GetVideoPath(channelID, videoID string) string
	SaveChannelInfo(channel *Channel) error
	SaveVideoMetadata(channelID string, video *Video) error
	FileExists(path string) bool
}
