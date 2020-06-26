const streamURL = '/hls/stream.m3u8';
// const streamURL = 'https://goth.land/hls/stream.m3u8'; // Uncomment me to point to remote video

// style hackings
window.VIDEOJS_NO_DYNAMIC_STYLE = true;

// Create the player for the first time
const player = videojs('video', null, function () {
  getStatus();
  setInterval(getStatus, 5000);
  setupPlayerEventHandlers();
})

player.ready(function () {
  console.log('Player ready.')
  resetPlayer(player);
});

function resetPlayer(player) {
  player.reset();
  player.src({ type: 'application/x-mpegURL', src: streamURL });
  if (app.isOnline) {
    player.poster('/thumbnail.jpg');
  } else {
    // Change this to some kind of offline image.
    player.poster('/img/logo.png');
  }
}

function setupPlayerEventHandlers() {
  const player = videojs('video');

  player.on('error', function (e) {
    console.log("Player error: ", e);
  })

  // player.on('loadeddata', function (e) {
  //   console.log("loadeddata");
  // })

  player.on('ended', function (e) {
    console.log("ended");
    resetPlayer(player);
  })
  //
  // player.on('abort', function (e) {
  //   console.log("abort");
  // })
  //
  // player.on('durationchange', function (e) {
  //   console.log("durationchange");
  // })
  //
  // player.on('stalled', function (e) {
  //   console.log("stalled");
  // })
  //
  player.on('playing', function (e) {
    if (playerRestartTimer) {
      clearTimeout(playerRestartTimer);
    }
  })
  //
  // player.on('waiting', function (e) {
  //   // console.log("waiting");
  // })
}

function restartPlayer() {
  try {
    console.log('restarting')
    const player = videojs('video');
    player.pause();
    player.src(player.src()); // Reload the same video
    player.load();
    player.play();
  } catch (e) {
    console.log(e)
  }

}