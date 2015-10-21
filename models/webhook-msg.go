package models

// WebhookMsg defines schema for message recieved from Slack webhook
type WebhookMsg struct {
	Token       string  `form:"token" binding:"required"`
	TeamID      string  `form:"team_id" binding:"required"`
	TeamDomain  string  `form:"team_domain"`
	ChannelID   string  `form:"channel_id" binding:"required"`
	ChannelName string  `form:"channel_name" binding:"required"`
	Timestamp   float64 `form:"timestamp"`
	UserID      string  `form:"user_id" binding:"required"`
	UserName    string  `form:"user_name" binding:"required"`
	Text        string  `form:"text"`
	TriggerWord string  `form:"trigger_word"`
	ServiceID   string  `form:"service_id"`
}
