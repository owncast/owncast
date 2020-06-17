async function setupApp() {  
  Vue.filter('plural', function (string, count) {
    if (count === 1) {
      return string;
    } else {
      return string + "s";
    }
  })

  window.app = new Vue({
    el: "#app-container",
    data: {
      streamStatus: "",
      viewerCount: 0,
      sessionMaxViewerCount: 0,
      overallMaxViewerCount: 0,
      messages: [],
      description: "",
      title: "",
    },
  });

  // init messaging interactions
  var appMessagingMisc = new Messaging();
  appMessagingMisc.init();

  const config = await new Config().init();
  app.description = autoLink(config.description, { embed: false });
  app.title = config.title;
}

var websocketReconnectTimer;
function setupWebsocket() {
  clearTimeout(websocketReconnectTimer)

  // Uncomment to point to somewhere other than goth.land
  const protocol = location.protocol == "https:" ? "wss" : "ws"
  var ws = new WebSocket(protocol + "://" + location.host + "/entry")

  // var ws = new WebSocket("wss://goth.land/entry")

  ws.onmessage = (e) => {
    const model = JSON.parse(e.data)

    // Ignore non-chat messages (such as keepalive PINGs)
    if (model.type !== SocketMessageTypes.CHAT) { return; }

    const message = new Message(model)
    
    const existing = this.app.messages.filter(function (item) {
      return item.id === message.id
    })
    
    if (existing.length === 0 || !existing) {
      this.app.messages.push(message);
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

  window.ws = ws;
}

setupApp()

setupWebsocket()

