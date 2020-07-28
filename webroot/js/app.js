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
    this.disableChatTimer = null;
    this.streamDurationTimer = null;

    // misc
    this.streamStatus = null;

    Vue.filter('plural', pluralize);

    // bindings
    this.vueAppMounted = this.vueAppMounted.bind(this);
    this.setConfigData = this.setConfigData.bind(this);
    this.setupWebsocket = this.setupWebsocket.bind(this);
    this.getStreamStatus = this.getStreamStatus.bind(this);
    this.getExtraUserContent = this.getExtraUserContent.bind(this);
    this.updateStreamStatus = this.updateStreamStatus.bind(this);
    this.handleNetworkingError = this.handleNetworkingError.bind(this);
    this.handleOfflineMode = this.handleOfflineMode.bind(this);
    this.handleOnlineMode = this.handleOnlineMode.bind(this);
    this.handleNetworkingError = this.handleNetworkingError.bind(this);
    this.handlePlayerReady = this.handlePlayerReady.bind(this);
    this.handlePlayerPlaying = this.handlePlayerPlaying.bind(this);
    this.handlePlayerEnded = this.handlePlayerEnded.bind(this);
    this.handlePlayerError = this.handlePlayerError.bind(this);
    this.setCurrentStreamDuration = this.setCurrentStreamDuration.bind(this);
  }

  init() {
    this.messagingInterface = new MessagingInterface();
    this.websocket = this.setupWebsocket();

    this.vueApp = new Vue({
      el: '#app-container',
      data: {
        playerOn: false,
        messages: [],
        overallMaxViewerCount: 0,
        sessionMaxViewerCount: 0,
        streamStatus: MESSAGE_OFFLINE, // Default state.
        viewerCount: 0,
        isOnline: false,
  
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
  // do all these things after Vue.js has mounted, else we'll get weird DOM issues.
  vueAppMounted() {
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

    this.getChatHistory();
  };

  setConfigData(data) {
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

  // websocket for messaging
  setupWebsocket() {
    var ws = new WebSocket(URL_WEBSOCKET);  
    ws.onopen = (e) => {
      if (this.websocketReconnectTimer) {
        clearTimeout(this.websocketReconnectTimer);
      }

      // If we're "online" then enable the chat.
      if (this.streamStatus && this.streamStatus.online) {
        this.messagingInterface.enableChat();
      }
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
      this.handleNetworkingError(`Socket error: ${json.Parse(e)}`);
      ws.close();
    };
    ws.onmessage = (e) => {
      const model = JSON.parse(e.data);

      // Send PONGs
      if (model.type === SOCKET_MESSAGE_TYPES.PING) {
        this.sendPong(ws);
        return;
      // Ignore non-chat messages (such as keepalive PINGs)
      } else if (model.type !== SOCKET_MESSAGE_TYPES.CHAT) {
        return; 
      }
      const message = new Message(model);
      this.addMessage(message);
    };
    this.websocket = ws;
    this.messagingInterface.setWebsocket(this.websocket);
  };

  sendPong(ws) {
    try {
      const pong = { type: 'PONG' };
      ws.send(JSON.stringify(pong));
    } catch (e) {
      console.log('PONG error:', e);
    }
  }

  addMessage(message) {
    const existing = this.vueApp.messages.filter(function (item) {
      return item.id === message.id;
    })
    if (existing.length === 0 || !existing) {
      this.vueApp.messages = [...this.vueApp.messages, message];
    }
  }

  // fetch /config data
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

  // fetch stream status
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
        this.handleOfflineMode();
        this.handleNetworkingError(`Stream status: ${error}`);
      });
  };

  // fetch content.md
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
        this.handleNetworkingError(`Fetch extra content: ${error}`);
      });
  };

  // fetch chat history
  getChatHistory() {
    fetch(URL_CHAT_HISTORY)
    .then(response => {
      if (!response.ok) {
        throw new Error(`Network response was not ok ${response.ok}`);
      }
      return response.json();
    })
    .then(data => {
      const formattedMessages = data.map(function (message) {
        return new Message(message);
      })
      this.vueApp.messages = formattedMessages.concat(this.vueApp.messages);
    })
    .catch(error => {
      this.handleNetworkingError(`Fetch getChatHistory: ${error}`);
    });
  }


  // handle UI things from stream status result
  updateStreamStatus(status = {}) {
    if (!status) {
      return;
    }
    // update UI
    this.vueApp.viewerCount = status.viewerCount;
    this.vueApp.sessionMaxViewerCount = status.sessionMaxViewerCount;
    this.vueApp.overallMaxViewerCount = status.overallMaxViewerCount;

    this.lastDisconnectTime = status.lastDisconnectTime;

    if (!this.streamStatus) {
      // display offline mode the first time we get status, and it's offline.
      if (!status.online) {
        this.handleOfflineMode();
      } else {
        this.handleOnlineMode();
      }
    } else {
      if (status.online && !this.streamStatus.online) {
        // stream has just come online.
        this.handleOnlineMode();
      } else if (!status.online && this.streamStatus.online) {
        // stream has just flipped offline.
        this.handleOfflineMode();
      }
    }

    // keep a local copy
    this.streamStatus = status;

    if (status.online) {
      // only do this if video is paused, so no unnecessary img fetches
      if (this.player.vjsPlayer && this.player.vjsPlayer.paused()) {
        this.player.setPoster();
      }
    }
  };
  
  // update vueApp.streamStatus text when online
  setCurrentStreamDuration() {
    // Default to something
    let streamDurationString = '';

    if (this.streamStatus.lastConnectTime) {
      const diff = (Date.now() - Date.parse(this.streamStatus.lastConnectTime)) / 1000;
      streamDurationString = secondsToHMMSS(diff);
    }
    this.vueApp.streamStatus = `${MESSAGE_ONLINE} ${streamDurationString}.`
  }
  
  handleNetworkingError(error) {
    console.log(`>>> App Error: ${error}`)
  };

  // stop status timer and disable chat after some time.
  handleOfflineMode() {
    this.vueApp.isOnline = false;
    clearInterval(this.streamDurationTimer);
    this.vueApp.streamStatus = MESSAGE_OFFLINE;
    if (this.streamStatus) {
      const remainingChatTime = TIMER_DISABLE_CHAT_AFTER_OFFLINE - (Date.now() - new Date(this.lastDisconnectTime));
      const countdown = (remainingChatTime < 0) ? 0 : remainingChatTime;
      this.disableChatTimer = setTimeout(this.messagingInterface.disableChat, countdown);
    }
  };

  // play video!
  handleOnlineMode() {
    this.vueApp.playerOn = true;
    this.vueApp.isOnline = true;
    this.vueApp.streamStatus = MESSAGE_ONLINE;

    this.player.startPlayer();
    clearTimeout(this.disableChatTimer);
    this.disableChatTimer = null;
    this.messagingInterface.enableChat();

    this.streamDurationTimer =
      setInterval(this.setCurrentStreamDuration, TIMER_STREAM_DURATION_COUNTER);
  }

  // when videojs player is ready, start polling for stream
  handlePlayerReady() {
    this.getStreamStatus();
    this.statusTimer = setInterval(this.getStreamStatus, TIMER_STATUS_UPDATE);
  };


  handlePlayerPlaying() {
    // do something?
  };


  // likely called some time after stream status has gone offline.
  // basically hide video and show underlying "poster"
  handlePlayerEnded() {
    this.vueApp.playerOn = false;
  };

  handlePlayerError() {
    // do something?
    this.handleOfflineMode();
    this.handlePlayerEnded();
  };
};
