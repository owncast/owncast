// style hackings
window.VIDEOJS_NO_DYNAMIC_STYLE = true;

var waitingTimeoutTimer;

// Wait until the player is setup before we start polling status
const player = videojs('video');

player.on('ready', function () {
  getStatus();
  setInterval(getStatus, 5000);
  setupPlayerEventHandlers();
})

function setupPlayerEventHandlers() {
  const player = videojs('video');

  player.on('error', function (e) {
    console.log(e);
  })

  player.on('loadeddata', function (e) {
    console.log("loadeddata");
  })

  player.on('ended', function (e) {
    console.log("ended");
  })

  // player.on('abort', function (e) {
  //   console.log("abort");
  // })

  // player.on('durationchange', function (e) {
  //   console.log("durationchange");
  // })

  // player.on('stalled', function (e) {
  //   console.log("stalled");
  // })
  
  player.on('playing', function (e) {
    // console.log("playing");
    clearTimeout(waitingTimeoutTimer);
  })

  player.on('waiting', function (e) {
    // console.log("waiting");

    // waitingTimeoutTimer = setTimeout(function () {
    //   restartPlayer();
    // }, 3000)
  })
}

function restartPlayer() {
  try {
    const player = videojs('video');

    player.src(player.src()); // Reload the same video
    player.load();
    player.play();
  } catch (e) {
    console.log(e)
  }

}