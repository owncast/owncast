function setupApp() {
  Vue.filter('plural', function (string, count) {
    if (count === 1) {
      return string
    } else {
      return string + "s"
    }
  })

  window.app = new Vue({
    el: "#app",
    data: {
      streamStatus: "",
      viewerCount: 0,
    },
  });

  window.messagesContainer = new Vue({
    el: "#messages-container",
    data: {
      messages: []
    }
  })

  window.chatForm = new Vue({
    el: "#chatForm",
    data: {
      message: {
        author: localStorage.author || "Viewer" + (Math.floor(Math.random() * 42) + 1),
        body: ""
      }
    },
    methods: {
      submitChatForm: function (e) {
        const message = new Message(this.message)
        message.id = uuidv4()
        localStorage.author = message.author
        const messageJSON = JSON.stringify(message)
        window.ws.send(messageJSON)
        e.preventDefault()

        this.message.body = ""
      }
    }
  });
}

async function getStatus() {
  const url = "/status";

  try {
    const response = await fetch(url);
    const status = await response.json(); // read response body and parse as JSON
    app.streamStatus = status.online
      ? "Stream is online."
      : "Stream is offline."
    
    app.viewerCount = status.viewerCount

  } catch (e) {
    app.streamStatus = "Stream server is offline."
    app.viewerCount = 0
  }

}

var websocketReconnectTimer
function setupWebsocket() {
  clearTimeout(websocketReconnectTimer)

  const protocol = location.protocol == "https:" ? "wss" : "ws"
  var ws = new WebSocket(protocol + "://" + location.host + "/entry")
  
  ws.onmessage = (e) => {
    const model = JSON.parse(e.data)
    const message = new Message(model)

    const existing = this.messagesContainer.messages.filter(function (item) {
      return item.id === message.id
    })
    
    if (existing.length === 0 || !existing) {
      this.messagesContainer.messages.push(message)
      scrollSmoothToBottom("messages-container")
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

function scrollSmoothToBottom(id) {
  const div = document.getElementById(id)
  $('#' + id).animate({
    scrollTop: div.scrollHeight - div.clientHeight
  }, 500)
}

function uuidv4() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
    const r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
}
