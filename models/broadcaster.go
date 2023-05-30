package models

import "time"

// Broadcaster represents the details around the inbound broadcasting connection.
type Broadcaster struct {
	Time          time.Time            `json:"time"`
	RemoteAddr    string               `json:"remoteAddr"`
	StreamDetails InboundStreamDetails `json:"streamDetails"`
}

// InboundStreamDetails represents an inbound broadcast stream.
type InboundStreamDetails struct {
	VideoCodec     string  `json:"videoCodec"`
	AudioCodec     string  `json:"audioCodec"`
	Encoder        string  `json:"encoder"`
	Width          int     `json:"width"`
	Height         int     `json:"height"`
	VideoBitrate   int     `json:"videoBitrate"`
	AudioBitrate   int     `json:"audioBitrate"`
	VideoFramerate float32 `json:"framerate"`
	VideoOnly      bool    `json:"-"`
}

// RTMPStreamMetadata is the raw metadata that comes in with a RTMP connection.
type RTMPStreamMetadata struct {
	VideoCodec     interface{} `json:"videocodecid"`
	AudioCodec     interface{} `json:"audiocodecid"`
	Encoder        string      `json:"encoder"`
	Width          int         `json:"width"`
	Height         int         `json:"height"`
	VideoBitrate   float32     `json:"videodatarate"`
	VideoFramerate float32     `json:"framerate"`
	AudioBitrate   float32     `json:"audiodatarate"`
}
