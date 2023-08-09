package constant

type CommentActionParams struct {
	UserID      int64
	VideoId     int64  `form:"video_id" binding:"required"`
	ActionType  int64  `form:"action_type" binding:"required,oneof=1 2"`
	CommentText string `form:"comment_text"`
	CommentId   int64  `form:"comment_id"`
}

type FavoriteActionParams struct {
	UserID     int64
	VideoId    int64 `form:"video_id" binding:"required"`
	ActionType int64 `form:"action_type" binding:"required,oneof=1 2"`
}

type RelationActionParams struct {
	UserID     int64
	ToUserID   int64 `form:"to_user_id" binding:"required"`
	ActionType int64 `form:"action_type" binding:"required,oneof=1 2"`
}
