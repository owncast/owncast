// https://docs.videojs.com/player

const VIDEO_ID = 'video';

const VIDEO_OPTIONS = {
  autoplay: false,
  liveui: true, // try this
  preload: 'auto',
  html5: {
    vhs: {
      // used to select the lowest bitrate playlist initially. This helps to decrease playback start time. This setting is false by default.
      enableLowInitialPlaylist: true,

    }
  },
  liveTracker: {
    trackingThreshold: 0,
  }
};


class OwncastPlayer {
  constructor() {
    window.VIDEOJS_NO_DYNAMIC_STYLE = true; // style override

    this.vjsPlayer = null;
    this.config = null;    

    this.appPlayerReadyCallback = null;
    this.appPlayerPlayingCallback = null;
    this.appPlayerEndedCallback = null;

    // bind all the things because safari
    this.startPlayer = this.startPlayer.bind(this);
    this.handleReady = this.handleReady.bind(this);
    this.handlePlaying = this.handlePlaying.bind(this);
    this.handleEnded = this.handleEnded.bind(this);
    this.handleError = this.handleError.bind(this);
  }

  init() {
    videojs.Hls.xhr.beforeRequest = function (options) {
      const cachebuster = Math.round(new Date().getTime() / 1000);
      options.uri = `${options.uri}?cachebust=${cachebuster}`;
      return options;
    };

    var options = VIDEO_OPTIONS

    // Add source from config
    sources = {sources: [this.getVideoSource()]}
    options = ({...VIDEO_OPTIONS, ...sources})

    this.vjsPlayer = videojs(VIDEO_ID, options);
    this.addAirplay();
    this.vjsPlayer.ready(this.handleReady);
  }

  setupPlayerCallbacks(callbacks) {
    const { onReady, onPlaying, onEnded, onError } = callbacks;

    this.appPlayerReadyCallback = onReady;
    this.appPlayerPlayingCallback = onPlaying;
    this.appPlayerEndedCallback = onEnded;
    this.appPlayerErrorCallback = onError;
  }

  setConfig(configData) {
    this.config = configData
  }

  getVideoSource() {
    configSource = {
      src: this.config["PrivateHLSPath"]+"/stream.m3u8",
      type: 'application/x-mpegURL',
    };

    return configSource
  }

  // play
  startPlayer() {
    this.log('Start playing');
    const source = { ...getVideoSource() }
    this.vjsPlayer.src(source);
    // this.vjsPlayer.play();
  };

  handleReady() {
    this.log('on Ready');
    this.vjsPlayer.on('error', this.handleError);
    this.vjsPlayer.on('playing', this.handlePlaying);
    this.vjsPlayer.on('ended', this.handleEnded);

    if (this.appPlayerReadyCallback) {
      // start polling
      this.appPlayerReadyCallback();
    }
  }

  handlePlaying() {
    this.log('on Playing');
    if (this.appPlayerPlayingCallback) {
      // start polling
      this.appPlayerPlayingCallback();
    }
  }

  handleEnded() {
    this.log('on Ended');
    if (this.appPlayerEndedCallback) {
      this.appPlayerEndedCallback();
    }
    this.setPoster();
  }

  handleError(e) {
    this.log(`on Error: ${JSON.stringify(e)}`);
    if (this.appPlayerEndedCallback) {
      this.appPlayerEndedCallback();
    }
  }

  setPoster() {
    const cachebuster = Math.round(new Date().getTime() / 1000);
    const poster = POSTER_THUMB + "?okhi=" + cachebuster;

    this.vjsPlayer.poster(poster);
  }

  log(message) {
    console.log(`>>> Player: ${message}`);
  }

  addAirplay() {
    videojs.hookOnce('setup', function (player) {
      if (window.WebKitPlaybackTargetAvailabilityEvent) {
        var videoJsButtonClass = videojs.getComponent('Button');
        var concreteButtonClass = videojs.extend(videoJsButtonClass, {

          // The `init()` method will also work for constructor logic here, but it is
          // deprecated. If you provide an `init()` method, it will override the
          // `constructor()` method!
          constructor: function () {
            videoJsButtonClass.call(this, player);
          }, // notice the comma

          handleClick: function () {
            const videoElement = document.getElementsByTagName('video')[0];
            videoElement.webkitShowPlaybackTargetPicker();
          }
        });

        var concreteButtonInstance = this.vjsPlayer.controlBar.addChild(new concreteButtonClass());
        concreteButtonInstance.addClass("vjs-airplay");
      }
    });
  }


}

export { OwncastPlayer };
