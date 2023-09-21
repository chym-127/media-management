package protocols

type ImportMediaReqProtocol struct {
	Medias []MediaItem `json:"medias"`
}

type MediaItem struct {
	ID          uint          `json:"id" bson:"id"`
	Title       string        `json:"title" bson:"title" binding:"required"`
	ReleaseDate int16         `json:"release_date" bson:"release_date" binding:"required"`
	Description string        `json:"description" bson:"description"`
	Score       float64       `json:"score" bson:"score"`
	Episodes    []EpisodeItem `json:"episodes" bson:"episodes"`
	PlayConfig  string        `json:"play_config" bson:"play_config"`
	PosterUrl   string        `json:"poster_url" bson:"poster_url"`
	FanartUrl   string        `json:"fanart_url" bson:"fanart_url"`
	Area        string        `json:"area" bson:"area"`
	Type        int8          `json:"type" bson:"type"`
}

type EpisodeItem struct {
	Index       uint   `json:"index" bson:"index"`
	Season      uint   `json:"season" bson:"season"`
	Title       string `json:"title" bson:"title"`
	ReleaseDate string `json:"release_date" bson:"release_date"`
	Description string `json:"description" bson:"description"`
	Expand      string `json:"expand" bson:"expand"`
	Url         string `json:"url" bson:"url"`
	LocalPath   string `json:"local_path" bson:"local_path"`
}

type Page struct {
	PageSize  int64 `json:"page_size" bson:"page_size"`
	PageLimit int64 `json:"page_limit" bson:"page_limit"`
}

type ListMediaReq struct {
	Page
	Keywords string `json:"keywords" bson:"keywords"`
}

type GetMediaReq struct {
	ID uint `json:"id"`
}

type MediaDownloadRecordItem struct {
	ID            uint   `json:"id" bson:"id"`
	Title         string `json:"title" bson:"title"`
	MediaID       uint   `json:"mediaId" bson:"media_id"`
	EpisodeCount  uint   `json:"episode_count" bson:"episode_count"`
	DownloadCount uint   `json:"download_count" bson:"download_count"`
	SuccessCount  uint   `json:"success_count" bson:"success_count"`
	FailedCount   uint   `json:"failed_count" bson:"failed_count"`
	Type          uint   `json:"type" bson:"type"` //1队列中 2下载中 3下载成功 4下载失败
}
