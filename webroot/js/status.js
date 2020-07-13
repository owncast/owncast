var playerRestartTimer;


function handleStatus(status) {
  clearTimeout(playerRestartTimer);
  if (!app.isOnline && status.online) {
    // The stream was offline, but now it's online.  Force start of playback after an arbitrary delay to make sure the stream has actual data ready to go.
    playerRestartTimer = setTimeout(restartPlayer, 3000);
  }

  app.streamStatus = status.online ? MESSAGE_ONLINE : MESSAGE_OFFLINE;

  app.viewerCount = status.viewerCount;
  app.sessionMaxViewerCount = status.sessionMaxViewerCount;
  app.overallMaxViewerCount = status.overallMaxViewerCount;
  app.isOnline = status.online;
  // setVideoPoster(app.isOnline);
}

function handleOffline() {
  const player = videojs(VIDEO_ID);
  player.poster(POSTER_DEFAULT);
  app.streamStatus = MESSAGE_OFFLINE;
  app.viewerCount = 0;
}

function getStatus() {
  fetch(URL_STATUS)
    .then(response => {
      if (!response.ok) {
        throw new Error(`Network response was not ok ${response.ok}`);
      }
      return response.json();
    })
    .then(json => {
      handleStatus(json);
    })
    .catch(error => {
      handleOffline();
    });
}
