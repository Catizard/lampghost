package dto

// Send a notification to frontend
type NotificationDto struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}
