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
        author: "",//localStorage.author || "Viewer" + (Math.floor(Math.random() * 42) + 1),
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
  let url = "https://util.real-ity.com:8042/status";

  try {
    let response = await fetch(url);
    let status = await response.json(); // read response body and parse as JSON
    app.streamStatus = status.online
      ? "Stream is online."
      : "Stream is offline."
    
    app.viewerCount = status.viewerCount

  } catch (e) {
    app.streamStatus = "Stream server is offline."
    app.viewerCount = 0
  }

}

function setupWebsocket() {
  const protocol = location.protocol == "https:" ? "wss" : "ws"
  var ws = new WebSocket("wss://util.real-ity.com:8042/entry")
  
  ws.onmessage = (e) => {
    var model = JSON.parse(e.data)
    var message = new Message(model)

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
    setTimeout(setupWebsocket, 5000)
  }

  ws.onerror = (e) => {
    setTimeout(setupWebsocket, 5000)
  }

  window.ws = ws
}

setupApp()
getStatus()
setupWebsocket()
// setInterval(getStatus, 5000)

// HLS Video setup
function setupVideo() {
  var video = document.getElementById("video")
  var videoSrc = "hls/stream.m3u8"
  if (Hls.isSupported()) {
    var hls = new Hls()
    hls.loadSource(videoSrc)
    hls.attachMedia(video)
    hls.on(Hls.Events.MANIFEST_PARSED, function () {
      video.play()
    });
  }
  
  // hls.js is not supported on platforms that do not have Media Source
  // Extensions (MSE) enabled.
  //
  // When the browser has built-in HLS support (check using `canPlayType`),
  // we can provide an HLS manifest (i.e. .m3u8 URL) directly to the video
  // element through the `src` property. This is using the built-in support
  // of the plain video element, without using hls.js.
  //
  // Note: it would be more normal to wait on the 'canplay' event below however
  // on Safari (where you are most likely to find built-in HLS support) the
  // video.src URL must be on the user-driven white-list before a 'canplay'
  // event will be emitted; the last video event that can be reliably
  // listened-for when the URL is not on the white-list is 'loadedmetadata'.
  else if (video.canPlayType("application/vnd.apple.mpegurl")) {
    video.src = videoSrc
    video.addEventListener("loadedmetadata", function () {
      video.play()
    });
  }
}

function scrollSmoothToBottom(id) {
  var div = document.getElementById(id)
  $('#' + id).animate({
    scrollTop: div.scrollHeight - div.clientHeight
  }, 500)
}

function uuidv4() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
    var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
}
