async function setupApp() {  
  Vue.filter('plural', pluralize);

  window.app = new Vue({
    el: "#app-container",
    data: {
      streamStatus: MESSAGE_OFFLINE, // Default state.
      viewerCount: 0,
      sessionMaxViewerCount: 0,
      overallMaxViewerCount: 0,
      messages: [],
      extraUserContent: "",
      isOnline: false,
      layout: "desktop",

      // from config
      logo: null,
      socialHandles: [],
      streamerName: "",
      summary: "",
      tags: [],
      title: "",
      appVersion: "",
    },
    watch: {
      messages: {
        deep: true,
        handler: function (newMessages, oldMessages) {
          if (newMessages.length !== oldMessages.length) {
            // jump to bottom
            jumpToBottom(appMessaging.scrollableMessagesContainer);
          }
        },
      },
    },
  });


  // init messaging interactions
  var appMessaging = new Messaging();
  appMessaging.init();

  const config = await new Config().init();
  app.logo = config.logo.small;
  app.socialHandles = config.socialHandles;
  app.streamerName = config.name;
  app.summary = config.summary && addNewlines(config.summary);
  app.tags =  config.tags;
  app.appVersion = config.version;
  app.title = config.title;
  window.document.title = config.title;

  getExtraUserContent(`${URL_PREFIX}${config.extraUserInfoFileName}`);
}

var websocketReconnectTimer;
function setupWebsocket() {
  clearTimeout(websocketReconnectTimer);

  var ws = new WebSocket(URL_WEBSOCKET);

  ws.onmessage = (e) => {
    const model = JSON.parse(e.data);

    // Ignore non-chat messages (such as keepalive PINGs)
    if (model.type !== SOCKET_MESSAGE_TYPES.CHAT) { return; }

    const message = new Message(model);
    
    const existing = this.app.messages.filter(function (item) {
      return item.id === message.id;
    })
    
    if (existing.length === 0 || !existing) {
      this.app.messages = [...this.app.messages, message];
    }
  }

  ws.onclose = (e) => {
    // connection closed, discard old websocket and create a new one in 5s
    ws = null;
    console.log("Websocket closed.")
    websocketReconnectTimer = setTimeout(setupWebsocket, 5000);
  }

  // On ws error just close the socket and let it re-connect again for now.
  ws.onerror = (e) => {
    console.log("Websocket error: ", e);
    ws.close();
  }

  window.ws = ws;
}

setupApp();

setupWebsocket();

