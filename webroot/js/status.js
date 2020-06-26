var playerRestartTimer;

async function getStatus() {
  const url = "/status";

  try {
    const response = await fetch(url);
    const status = await response.json();

    clearTimeout(playerRestartTimer);

    if (!app.isOnline && status.online) {
      // The stream was offline, but now it's online.  Force start of playback after an arbitrary
      // delay to make sure the stream has actual data ready to go.
      playerRestartTimer = setTimeout(function () {
        restartPlayer();
      }, 3000);
    }

    app.streamStatus = status.online
      ? "Stream is online."
      : "Stream is offline."

    const player = videojs('video');
    if (app.isOnline) {
      player.poster('/thumbnail.jpg');
    } else {
      // Change this to some kind of offline image.
      player.poster('/img/logo.png');
    }

    app.viewerCount = status.viewerCount;
    app.sessionMaxViewerCount = status.sessionMaxViewerCount;
    app.overallMaxViewerCount = status.overallMaxViewerCount;
    app.isOnline = status.online;

  } catch (e) {
    app.streamStatus = "Stream server is offline."
    app.viewerCount = 0
  }

}