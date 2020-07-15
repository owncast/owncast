// https://docs.videojs.com/player

class OwncastPlayer {
  constructor() {
    window.VIDEOJS_NO_DYNAMIC_STYLE = true; // style override

    this.vjsPlayer = null;

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
    this.vjsPlayer = videojs(VIDEO_ID, VIDEO_OPTIONS);
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

  // play
  startPlayer() {
    this.log('Start playing');

    // this.vjsPlayer.load(); // causes errors? works without?
    this.vjsPlayer.play();
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

    // only do this if video is paused, so no unnecessary img fetches
    if (this.vjsPlayer.paused()) {
      this.vjsPlayer.poster(poster);
    }
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

