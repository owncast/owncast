package controllers

import (
	"net/http"

	"github.com/owncast/owncast/metrics"
)

func GetPlayInstructions(w http.ResponseWriter, r *http.Request) {
	instructions := metrics.GetPlayInstructions()

	WriteResponse(w, instructions)
}
