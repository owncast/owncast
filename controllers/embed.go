package controllers

import (
	"net/http"
)

// GetChatEmbedreadwrite gets the embed for readwrite chat.
func GetChatEmbedreadwrite(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index-standalone-chat-readwrite.html", http.StatusMovedPermanently)
}

// GetChatEmbedreadonly gets the embed for readonly chat.
func GetChatEmbedreadonly(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index-standalone-chat-readonly.html", http.StatusMovedPermanently)
}

// GetVideoEmbed gets the embed for video.
func GetVideoEmbed(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index-video-only.html", http.StatusMovedPermanently)
}
