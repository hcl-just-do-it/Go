package model

// 拉取评论返回的信息

type Comment struct {
	VideoID    int64  `json:"video_id"`
	Content    string `json:"content"`                                    // 评论内容
	CreateDate string `json:"create_date"`                                // 评论发布日期，格式 mm-dd
	CommentID  int64  `gorm:"primaryKey autoIncrement" json:"comment_id"` // 评论id
	UserID     int64  `json:"user_id"`                                    // 评论用户信息
}

type CommentInfo struct {
	VideoID    int64    `json:"video_id"`
	Content    string   `json:"content"`     // 评论内容
	CreateDate string   `json:"create_date"` // 评论发布日期，格式 mm-dd
	CommentID  int64    `json:"comment_id"`  // 评论id
	User       UserInfo `json:"user"`        // 评论用户信息
}
