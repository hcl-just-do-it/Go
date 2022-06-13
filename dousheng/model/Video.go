package model

type Video struct {
	Id            int64  `gorm:"primaryKey autoIncrement" json:"video_id,omitempty"`
	AuthorID      int64  `json:"author_id"`
	PlayUrl       string `gorm:"type:varchar(255)" json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `gorm:"type:varchar(255)" json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	Title         string `gorm:"type:varchar(50)" json:"title,omitempty"`
	PublishTime   int64  `json:"publish_time"`
}
