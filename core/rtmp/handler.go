package rtmp

import (
	"bytes"
	"errors"
	"io"
	"os"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/yutopp/go-flv"
	flvtag "github.com/yutopp/go-flv/tag"
	yutmp "github.com/yutopp/go-rtmp"
	rtmpmsg "github.com/yutopp/go-rtmp/message"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/core"
	"github.com/gabek/owncast/core/ffmpeg"
	"github.com/gabek/owncast/utils"
)

var _ yutmp.Handler = (*Handler)(nil)

// Handler An RTMP connection handler
type Handler struct {
	yutmp.DefaultHandler
	flvFile *os.File
	flvEnc  *flv.Encoder
}

//OnServe handles the "OnServe" of the rtmp service
func (h *Handler) OnServe(conn *yutmp.Conn) {
}

//OnConnect handles the "OnConnect" of the rtmp service
func (h *Handler) OnConnect(timestamp uint32, cmd *rtmpmsg.NetConnectionConnect) error {
	// log.Printf("OnConnect: %#v", cmd)

	return nil
}

//OnCreateStream handles the "OnCreateStream" of the rtmp service
func (h *Handler) OnCreateStream(timestamp uint32, cmd *rtmpmsg.NetConnectionCreateStream) error {
	// log.Printf("OnCreateStream: %#v", cmd)

	return nil
}

//OnPublish handles the "OnPublish" of the rtmp service
func (h *Handler) OnPublish(timestamp uint32, cmd *rtmpmsg.NetStreamPublish) error {
	// log.Printf("OnPublish: %#v", cmd)
	log.Println("Incoming stream connected.")

	if cmd.PublishingName != config.Config.VideoSettings.StreamingKey {
		return errors.New("invalid streaming key; rejecting incoming stream")
	}

	if _isConnected {
		return errors.New("stream already running; can not overtake an existing stream")
	}

	// Record streams as FLV
	p := utils.GetTemporaryPipePath()
	syscall.Mkfifo(p, 0666)

	f, err := os.OpenFile(p, os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		return err
	}
	h.flvFile = f

	enc, err := flv.NewEncoder(f, flv.FlagsAudio|flv.FlagsVideo)
	if err != nil {
		_ = f.Close()
		return err
	}
	h.flvEnc = enc

	transcoder := ffmpeg.NewTranscoder()
	go transcoder.Start()

	_isConnected = true
	core.SetStreamAsConnected()

	return nil
}

//OnSetDataFrame handles the setting of the data frame
func (h *Handler) OnSetDataFrame(timestamp uint32, data *rtmpmsg.NetStreamSetDataFrame) error {
	r := bytes.NewReader(data.Payload)

	var script flvtag.ScriptData
	if err := flvtag.DecodeScriptData(r, &script); err != nil {
		log.Printf("Failed to decode script data: Err = %+v", err)
		return nil // ignore
	}

	// log.Printf("SetDataFrame: Script = %#v", script)

	if err := h.flvEnc.Encode(&flvtag.FlvTag{
		TagType:   flvtag.TagTypeScriptData,
		Timestamp: timestamp,
		Data:      &script,
	}); err != nil {
		log.Printf("Failed to write script data: Err = %+v", err)
	}

	return nil
}

//OnAudio handles when we get audio from the rtmp service
func (h *Handler) OnAudio(timestamp uint32, payload io.Reader) error {
	var audio flvtag.AudioData
	if err := flvtag.DecodeAudioData(payload, &audio); err != nil {
		return err
	}

	flvBody := new(bytes.Buffer)
	if _, err := io.Copy(flvBody, audio.Data); err != nil {
		return err
	}
	audio.Data = flvBody

	// log.Printf("FLV Audio Data: Timestamp = %d, SoundFormat = %+v, SoundRate = %+v, SoundSize = %+v, SoundType = %+v, AACPacketType = %+v, Data length = %+v",
	// 	timestamp,
	// 	audio.SoundFormat,
	// 	audio.SoundRate,
	// 	audio.SoundSize,
	// 	audio.SoundType,
	// 	audio.AACPacketType,
	// 	len(flvBody.Bytes()),
	// )

	if err := h.flvEnc.Encode(&flvtag.FlvTag{
		TagType:   flvtag.TagTypeAudio,
		Timestamp: timestamp,
		Data:      &audio,
	}); err != nil {
		log.Printf("Failed to write audio: Err = %+v", err)
	}

	return nil
}

//OnVideo handles when we video from the rtmp service
func (h *Handler) OnVideo(timestamp uint32, payload io.Reader) error {
	var video flvtag.VideoData
	if err := flvtag.DecodeVideoData(payload, &video); err != nil {
		return err
	}

	flvBody := new(bytes.Buffer)
	if _, err := io.Copy(flvBody, video.Data); err != nil {
		return err
	}
	video.Data = flvBody

	// log.Printf("FLV Video Data: Timestamp = %d, FrameType = %+v, CodecID = %+v, AVCPacketType = %+v, CT = %+v, Data length = %+v",
	// 	timestamp,
	// 	video.FrameType,
	// 	video.CodecID,
	// 	video.AVCPacketType,
	// 	video.CompositionTime,
	// 	len(flvBody.Bytes()),
	// )

	if err := h.flvEnc.Encode(&flvtag.FlvTag{
		TagType:   flvtag.TagTypeVideo,
		Timestamp: timestamp,
		Data:      &video,
	}); err != nil {
		log.Printf("Failed to write video: Err = %+v", err)
	}

	return nil
}

//OnClose handles the closing of the rtmp connection
func (h *Handler) OnClose() {
	log.Printf("OnClose of the rtmp service")

	if h.flvFile != nil {
		_ = h.flvFile.Close()
	}

	_isConnected = false
	core.SetStreamAsDisconnected()
}
