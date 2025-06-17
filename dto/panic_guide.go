package dto

// PanicGuideRequest represents the JSON body for creating or updating a panic guide.
type PanicGuideRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=128"`
	Description string `json:"description" binding:"required"`
}

// BookmarkRequest represents the JSON body or query parameters for bookmarking a panic guide.
type BookmarkRequest struct {
	UserID       uint `json:"user_id" binding:"required"`
	PanicGuideID uint `json:"panic_guide_id" binding:"required"`
}

// ListBookmarksRequest represents the JSON body or query parameters for listing a user's bookmarks.
type ListBookmarksRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}
