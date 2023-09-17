package protocols

type ImportMediaReqProtocol struct {
	Medias []MediaItem `json:"medias"`
}

type MediaItem struct {
	Title       string        `json:"title" bson:"title" binding:"required"`
	ReleaseDate int16         `json:"releaseDate" bson:"release_date" binding:"required"`
	Description string        `json:"description" bson:"description"`
	Score       int16         `json:"score" bson:"score"`
	Episodes    []EpisodeItem `json:"episodes" bson:"episodes"`
	PlayConfig  string        `json:"playConfig" bson:"play_config"`
	PosterUrl   string        `json:"posterUrl" bson:"poster_url"`
	FanartUrl   string        `json:"fanartUrl" bson:"fanart_url"`
	Area        string        `json:"area" bson:"area"`
	Type        int8          `json:"type" bson:"type"`
}

type EpisodeItem struct {
	Url   string `json:"url" bson:"url"`
	Title string `json:"title" bson:"title"`
	Index int8   `json:"index" bson:"index"`
}

type Page struct {
	PageSize  int64 `json:"page_size" bson:"page_size"`
	PageLimit int64 `json:"page_limit" bson:"page_limit"`
}

type ListMediaReq struct {
	Page
	Keywords string `json:"keywords" bson:"keywords"`
}
