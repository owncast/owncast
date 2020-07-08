
// function setupApp() {  
//   Vue.filter('plural', pluralize);

//   window.app = new Vue({
//     el: '#app-container',
//     data: {
//       streamStatus: MESSAGE_OFFLINE, // Default state.
//       viewerCount: 0,
//       sessionMaxViewerCount: 0,
//       overallMaxViewerCount: 0,
//       messages: [],
//       isOnline: false,
//       layout: 'desktop',

//       // from config
//       logo: null,
//       socialHandles: [],
//       streamerName: '',
//       summary: '',
//       tags: [],
//       title: '',
//       extraUserContent: '',
//       appVersion: '',
//     },
//     watch: {
//       messages: {
//         deep: true,
//         handler: function (newMessages, oldMessages) {
//           if (newMessages.length !== oldMessages.length) {
//             // jump to bottom
//             jumpToBottom(appMessaging.scrollableMessagesContainer);
//           }
//         },
//       },
//     },
//   });


//   // init messaging interactions
//   var appMessaging = new Messaging();
//   appMessaging.init();

//   const config = await new Config().init();
//   app.logo = config.logo.small;
//   app.socialHandles = config.socialHandles;
//   app.streamerName = config.name;
//   app.summary = config.summary && addNewlines(config.summary);
//   app.tags =  config.tags;
//   app.appVersion = config.version;
//   app.title = config.title;
//   window.document.title = config.title;

//   getExtraUserContent(`${URL_PREFIX}${config.extraUserInfoFileName}`);
// }

// var websocketReconnectTimer;
// function setupWebsocket() {
//   clearTimeout(websocketReconnectTimer);

//   var ws = new WebSocket(URL_WEBSOCKET);

//   ws.onmessage = (e) => {
//     const model = JSON.parse(e.data);

//     // Ignore non-chat messages (such as keepalive PINGs)
//     if (model.type !== SOCKET_MESSAGE_TYPES.CHAT) { return; }

//     const message = new Message(model);
    
//     const existing = this.app.messages.filter(function (item) {
//       return item.id === message.id;
//     })
    
//     if (existing.length === 0 || !existing) {
//       this.app.messages = [...this.app.messages, message];
//     }
//   }

//   ws.onclose = (e) => {
//     // connection closed, discard old websocket and create a new one in 5s
//     ws = null;
//     console.log('Websocket closed.');
//     websocketReconnectTimer = setTimeout(setupWebsocket, 5000);
//   }

//   // On ws error just close the socket and let it re-connect again for now.
//   ws.onerror = (e) => {
//     console.log('Websocket error: ', e);
//     ws.close();
//   }

//   window.ws = ws;
// }

// setupApp();

// setupWebsocket();
////////////////////////
class Owncast {
  constructor() {
    this.websocket = null;
    this.websocketReconnectTimer = null;
    this.player;
    this.playerRestartTimer = null;
    this.configData;
    this.vueApp;
    this.messagingInterface = null;

    this.streamIsOnline = false;
    Vue.filter('plural', pluralize);

    //bindings
    this.handleNetworkingError = this.handleNetworkingError.bind(this);
    this.handlePlayerReady = this.handlePlayerReady.bind(this);
    this.handlePlayerPlaying = this.handlePlayerPlaying.bind(this);
    this.handleOfflineMode = this.handleOfflineMode.bind(this);
  }

  init() {
    this.getConfig();
    this.player = new Player({ castApp: this });
  
    this.messagingInterface = new MessagingInterface();
    this.vueApp = new Vue({
      el: '#app-container',
      data: {
        isOnline: false,
        layout: hasTouchScreen() ? 'touch' : 'desktop',
        messages: [],
        overallMaxViewerCount: 0,
        sessionMaxViewerCount: 0,
        streamStatus: MESSAGE_OFFLINE, // Default state.
        viewerCount: 0,
  
        // from config
        appVersion: '',
        extraUserContent: '',
        logo: null,
        socialHandles: [],
        streamerName: '',
        summary: '',
        tags: [],
        title: '',
      },
      watch: {
        messages: {
          deep: true,
          handler: this.messagingInterface.onReceivedMessages,
        },
      },
    });

    this.setupWebsocket();
    this.messagingInterface.init(this.websocket);

  }

  setConfigData(data) {

    this.vueApp.appVersion = data.version;
    this.vueApp.logo = data.logo.small;
    this.vueApp.socialHandles = data.socialHandles;
    this.vueApp.streamerName = data.name;
    this.vueApp.summary = data.summary && addNewlines(data.summary);
    this.vueApp.tags =  data.tags;
    this.vueApp.title = data.title;

    window.document.title = data.title;

    this.getExtraUserContent(`${URL_PREFIX}${data.extraUserInfoFileName}`);

    this.configData = data;

  }

  setupWebsocket() {
    clearTimeout(this.websocketReconnectTimer);
  
    var ws = new WebSocket(URL_WEBSOCKET);
  
    ws.onmessage = (e) => {
      const model = JSON.parse(e.data);
  
      // Ignore non-chat messages (such as keepalive PINGs)
      if (model.type !== SOCKET_MESSAGE_TYPES.CHAT) {
          return; 
        }
  
      const message = new Message(model);
      
      const existing = this.vueApp.messages.filter(function (item) {
        return item.id === message.id;
      })
      
      if (existing.length === 0 || !existing) {
        this.vueApp.messages = [...this.vueApp.messages, message];
      }
    }
  
    ws.onclose = (e) => {
      // connection closed, discard old websocket and create a new one in 5s
      ws = null;
      console.log('Websocket closed.');
      this.websocketReconnectTimer = setTimeout(this.setupWebsocket, 5000);
    }
  
    // On ws error just close the socket and let it re-connect again for now.
    ws.onerror = (e) => {
      console.log('Websocket error: ', e);
      ws.close();
    }
  
    this.websocket = ws;
  }

  getConfig() {
    fetch(URL_CONFIG)
    .then(response => {
      if (!response.ok) {
        throw new Error(`Network response was not ok ${response.ok}`);
      }
      return response.json();
    })
    .then(json => {
      this.setConfigData(json);
    })
    .catch(error => {
      this.handleNetworkingError(error);
    });
  }

  getStreamStatus() {
    fetch(URL_STATUS)
      .then(response => {
        if (!response.ok) {
          throw new Error(`Network response was not ok ${response.ok}`);
        }
        return response.json();
      })
      .then(json => {
        this.updateStreamStatus(json);
      })
      .catch(error => {
        console.log("this?", this)
        this.handleOfflineMode();
        this.handleNetworkingError(error);
      });
  }


  getExtraUserContent(path) {
    fetch(path)
      .then(response => {
        if (!response.ok) {
          throw new Error(`Network response was not ok ${response.ok}`);
        }
        return response.text();
      })
      .then(text => {
        const descriptionHTML = new showdown.Converter().makeHtml(text);
        this.vueApp.extraUserContent = descriptionHTML;
      })
      .catch(error => {
        this.handleNetworkingError(error);
      });
  }


  handleNetworkingError(error) {
    console.log(`>>> App Error: ${error}`)
  }

  handleOfflineMode() {
    this.player.setPoster(false);
    this.vueApp.streamStatus = MESSAGE_OFFLINE;
    this.vueApp.viewerCount = 0;
  }

  handlePlayerReady() {
    this.getStreamStatus();
    setInterval(this.getStreamStatus, 5000); // 
  }

  handlePlayerPlaying() {
    if (this.playerRestartTimer) {
      clearTimeout(this.playerRestartTimer);
    }
  }

  updateStreamStatus(status) {
    clearTimeout(this.playerRestartTimer);
    if (status.online && !this.streamIsOnline) {
      // The stream was offline, but now it's online.  Force start of playback after an arbitrary delay to make sure the stream has actual data ready to go.
      this.playerRestartTimer = setTimeout(this.player.restartPlayer, 3000);
    }
  
    this.vueApp.streamStatus = status.online ? MESSAGE_ONLINE : MESSAGE_OFFLINE;
  
    this.vueApp.viewerCount = status.viewerCount;
    this.vueApp.sessionMaxViewerCount = status.sessionMaxViewerCount;
    this.vueApp.overallMaxViewerCount = status.overallMaxViewerCount;

    this.streamIsOnline = status.online;

    this.vueApp.isOnline = this.streamIsOnline;

    this.player.setPoster(this.streamIsOnline);
  }
}