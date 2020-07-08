// https://docs.videojs.com/player
class Player {
  constructor({ castApp }) {
    window.VIDEOJS_NO_DYNAMIC_STYLE = true; // style override

    this.castApp = castApp;

    console.log("======== ", castApp)

    const options = {
      liveui: true,
      sources: [{
        src: URL_STREAM,
        type: 'application/x-mpegURL',
      }],
    };
    this.player = videojs(VIDEO_ID, options);
    this.addAirplay();
    this.setupPlayerEventHandlers();

    this.player.ready(this.handleReady.bind(this));
  }
  resetPlayer() {
    this.player.reset();
    // this.player.src({ type: 'application/x-mpegURL', src: URL_STREAM });
    // setVideoPoster(app.isOnline); // need?
  }
  handleReady() {
    this.log('Ready')
    if (this.castApp.handlePlayerReady) {
      this.castApp.handlePlayerReady();
    }
    this.player.reset();
  }

  restartPlayer() {
    this.log('Start player');
    this.player.pause();
    // this.player.src(player.src()); // Reload the same video
    this.player.load();
    this.player.play();
  }

  setPoster(online) {
    const cachebuster = Math.round(new Date().getTime() / 1000);
    const poster = online ? POSTER_THUMB + "?okhi=" + cachebuster : POSTER_DEFAULT;
    this.player.poster(poster);
  }


  setupPlayerEventHandlers() {
    this.player.on('error', function (e) {
      this.log('Error: ', e);
    });

    this.player.on('ended', function (e) {
      this.log('ended');
      this.resetPlayer();
    });

    this.player.on('playing', (function (e) {
      if (this.castApp.handlePlayerPlaying) {
        this.castApp.handlePlayerPlaying();
      }
    }).bind(this));
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
    
        var concreteButtonInstance = this.player.controlBar.addChild(new concreteButtonClass());
        concreteButtonInstance.addClass("vjs-airplay");
      }
    });  
  }
  

}


  // // Create the player for the first time
  // const player = videojs(VIDEO_ID, null, function () {
  //   getStatus();
  //   setInterval(getStatus, 5000);
  //   setupPlayerEventHandlers();
  // })

  // player.ready(function () {
  //   console.log('Player ready.')
  //   resetPlayer(player);
  // });

  // function resetPlayer(player) {
  //   player.reset();
  //   player.src({ type: 'application/x-mpegURL', src: URL_STREAM });
  //   setVideoPoster(app.isOnline);
  // }

  // function setupPlayerEventHandlers() {
  //   const player = videojs(VIDEO_ID);

  //   player.on('error', function (e) {
  //     console.log("Player error: ", e);
  //   })

  //   // player.on('loadeddata', function (e) {
  //   //   console.log("loadeddata");
  //   // })

  //   player.on('ended', function (e) {
  //     console.log("ended");
  //     resetPlayer(player);
  //   })
  //   //
  //   // player.on('abort', function (e) {
  //   //   console.log("abort");
  //   // })
  //   //
  //   // player.on('durationchange', function (e) {
  //   //   console.log("durationchange");
  //   // })
  //   //
  //   // player.on('stalled', function (e) {
  //   //   console.log("stalled");
  //   // })
  //   //
  //   player.on('playing', function (e) {
  //     if (playerRestartTimer) {
  //       clearTimeout(playerRestartTimer);
  //     }
  //   })
  //   //
  //   // player.on('waiting', function (e) {
  //   //   // console.log("waiting");
  //   // })
  // }

  // restartPlayer() {
  //   try {
  //     const player = videojs(VIDEO_ID);
  //     player.pause();
  //     player.src(player.src()); // Reload the same video
  //     player.load();
  //     player.play();
  //   } catch (e) {
  //     console.log(e)
  //   }

  // }
// }