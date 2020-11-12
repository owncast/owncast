package controllers

import (
	"net/http"
)

// GetChatEmbed gets the embed for chat.
func GetChatEmbed(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index-standalone-chat.html", http.StatusMovedPermanently)
}

// GetVideoEmbed gets the embed for video.
func GetVideoEmbed(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index-video-only.html", http.StatusMovedPermanently)
}
