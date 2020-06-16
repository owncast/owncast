function setupApp() {
  Vue.filter('plural', function (string, count) {
    if (count === 1) {
      return string;
    } else {
      return string + "s";
    }
  })

  window.app = new Vue({
    el: "#stream-info",
    data: {
      streamStatus: "",
      viewerCount: 0,
      sessionMaxViewerCount: 0,
      overallMaxViewerCount: 0,
    },
  });

  window.messagesContainer = new Vue({
    el: "#messages-container",
    data: {
      messages: []
    }
  })


  // style hackings
  window.VIDEOJS_NO_DYNAMIC_STYLE = true;
  if (hasTouchScreen()) {
    setVHvar();
    window.onorientationchange = handleOrientationChange;
  }


  // init messaging interactions
  var appMessagingMisc = new Messaging();
  appMessagingMisc.init();

  const config = new Config();
}

async function getStatus() {
  let url = "/status";

  try {
    const response = await fetch(url);
    const status = await response.json(); // read response body and parse as JSON

    if (!app.isOnline && status.online) {
      // The stream was offline, but now it's online.  Force start of playback after an arbitrary
      // delay to make sure the stream has actual data ready to go.
      setTimeout(function () {
        var player = videojs('video');
        player.pause()
        player.src(player.src()); // Reload the same video
        player.load();
        player.play();
      }, 3000)
      
    }
    
    app.streamStatus = status.online
      ? "Stream is online."
      : "Stream is offline."
    
    app.viewerCount = status.viewerCount;
    app.sessionMaxViewerCount = status.sessionMaxViewerCount;
    app.overallMaxViewerCount = status.overallMaxViewerCount;
    app.isOnline = status.online;
    
  } catch (e) {
    app.streamStatus = "Stream server is offline."
    app.viewerCount = 0
  }

}

var websocketReconnectTimer;
function setupWebsocket() {
  clearTimeout(websocketReconnectTimer)

  // Uncomment to point to somewhere other than goth.land
  // const protocol = location.protocol == "https:" ? "wss" : "ws"
  // var ws = new WebSocket(protocol + "://" + location.host + "/entry")

  var ws = new WebSocket("wss://goth.land/entry")

  ws.onmessage = (e) => {
    const model = JSON.parse(e.data)

    // Ignore non-chat messages (such as keepalive PINGs)
    if (model.type !== SocketMessageTypes.CHAT) { return; }

    const message = new Message(model)

    const existing = this.messagesContainer.messages.filter(function (item) {
      return item.id === message.id
    })
    
    if (existing.length === 0 || !existing) {
      this.messagesContainer.messages.push(message);
      setTimeout(() => { jumpToBottom("#messages-container"); } , 50); // could be better. is there a sort of Vue "componentDidUpdate" we can do this on?
    }
  }

  ws.onclose = (e) => {
    // connection closed, discard old websocket and create a new one in 5s
    ws = null
    console.log("Websocket closed.")
    websocketReconnectTimer = setTimeout(setupWebsocket, 5000)
  }

  // On ws error just close the socket and let it re-connect again for now.
  ws.onerror = (e) => {
    console.log("Websocket error: ", e)
    ws.close()
  }

  window.ws = ws
}

setupApp()
getStatus()
setupWebsocket()
setInterval(getStatus, 5000)

