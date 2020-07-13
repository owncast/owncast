class Owncast {
  constructor() {
    this.player;    

    this.websocket = null;
    this.configData;
    this.vueApp;
    this.messagingInterface = null;

    // timers
    this.websocketReconnectTimer = null;
    this.playerRestartTimer = null;
    this.offlineTimer = null;
    this.statusTimer = null;

    // misc
    this.streamIsOnline = false;
    Vue.filter('plural', pluralize);

    // bindings
    this.handleNetworkingError = this.handleNetworkingError.bind(this);
    this.handlePlayerReady = this.handlePlayerReady.bind(this);
    this.handlePlayerPlaying = this.handlePlayerPlaying.bind(this);
    this.handleOfflineMode = this.handleOfflineMode.bind(this);
    this.getStreamStatus = this.getStreamStatus.bind(this);
  }

  init() {
    this.messagingInterface = new MessagingInterface();
    this.websocket = this.setupWebsocket();

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
        logo: TEMP_IMAGE,
        logoLarge: TEMP_IMAGE,
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
      mounted: this.vueAppMounted,
    });
  }
  vueAppMounted = () => {
    this.getConfig();
    this.messagingInterface.init();

    this.player = new OwncastPlayer();
    this.player.setupPlayerCallbacks({
      onReady: this.handlePlayerReady,
      onPlaying: this.handlePlayerPlaying,
      onEnded: this.handlePlayerEnded,
      onError: this.handlePlayerError,
    });
    this.player.init();
  };

  setConfigData = (data) => {
    this.vueApp.appVersion = data.version;
    this.vueApp.logo = data.logo.small;
    this.vueApp.logoLarge = data.logo.large;
    this.vueApp.socialHandles = data.socialHandles;
    this.vueApp.streamerName = data.name;
    this.vueApp.summary = data.summary && addNewlines(data.summary);
    this.vueApp.tags =  data.tags;
    this.vueApp.title = data.title;

    window.document.title = data.title;

    this.getExtraUserContent(`${URL_PREFIX}${data.extraUserInfoFileName}`);

    this.configData = data;
  }

  // for messaging
  setupWebsocket = () => {
    var ws = new WebSocket(URL_WEBSOCKET);  
    ws.onopen = (e) => {
      if (this.websocketReconnectTimer) {
        clearTimeout(this.websocketReconnectTimer);
      }
      this.messagingInterface.enableChat();
    };
    ws.onclose = (e) => {
      // connection closed, discard old websocket and create a new one in 5s
      this.websocket = null;
      this.messagingInterface.disableChat();
      this.handleNetworkingError('Websocket closed.');
      this.websocketReconnectTimer = setTimeout(this.setupWebsocket, TIMER_WEBSOCKET_RECONNECT);
    };
    // On ws error just close the socket and let it re-connect again for now.
    ws.onerror = e => {
      this.handleNetworkingError(`Stream status: ${e}`);
      ws.close();
    };
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
    };
    this.websocket = ws;
    this.messagingInterface.setWebsocket(this.websocket);
  };

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
      this.handleNetworkingError(`Fetch config: ${error}`);
    });
  }

  getStreamStatus = () => {
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
        this.handleOfflineMode();
        this.handleNetworkingError(`Stream status: ${error}`);
      });
  };


  getExtraUserContent = (path) => {
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
        this.handleNetworkingError(`Fetch extra content: ${error}`);
      });
  };

  updateStreamStatus = (status) => {
    // update UI
    this.vueApp.isOnline = status.online;
    this.vueApp.streamStatus = status.online ? MESSAGE_ONLINE : MESSAGE_OFFLINE;
    this.vueApp.viewerCount = status.viewerCount;
    this.vueApp.sessionMaxViewerCount = status.sessionMaxViewerCount;
    this.vueApp.overallMaxViewerCount = status.overallMaxViewerCount;

    if (status.online && !this.streamIsOnline) {
      // stream has just come online.
      this.handleOnlineMode();
    } else if (!status.online && this.streamIsOnline) {
      // stream has just gone offline.
      this.handleOfflineMode();
    }

    if (status.online) {
      this.player.setPoster();
    }

    this.streamIsOnline = status.online;
  };
  
  handleNetworkingError = (error) => {
    console.log(`>>> App Error: ${error}`)
  };

  handleOfflineMode = () => {
    this.streamIsOnline = false;

    this.vueApp.streamStatus = MESSAGE_OFFLINE;
    this.vueApp.isOnline = status.online;
    // TODO: start offline timer to disable chat.
  };

  handleOnlineMode = () => {
    this.player.startPlayer();
  }

  // when videojs player is ready, start polling for stream
  handlePlayerReady = () => {
    this.getStreamStatus();
    this.statusTimer = setInterval(this.getStreamStatus, TIMER_STATUS_UPDATE);
  };


  handlePlayerPlaying = () => {
    // do something?
  };


  handlePlayerEnded = () => {
    // do something?
    this.handleOfflineMode();
  };

  handlePlayerError = () => {
    // do something?
    this.handleOfflineMode();
    // stop timers?
  };
};
