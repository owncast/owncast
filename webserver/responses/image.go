package responses

import (
	"log"
	"net/http"
	"strconv"
)

func WriteBytesAsImage(data []byte, contentType string, w http.ResponseWriter, cacheSeconds int) {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(cacheSeconds))

	if _, err := w.Write(data); err != nil {
		log.Println("unable to write image.")
	}
}
