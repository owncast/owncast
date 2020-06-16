videojs.hookOnce('setup', function (player) {
  if (window.WebKitPlaybackTargetAvailabilityEvent) {
    var videoJsButtonClass = videojs.getComponent('Button');
    var concreteButtonClass = videojs.extend(videoJsButtonClass, {

      // The `init()` method will also work for constructor logic here, but it is 
      // deprecated. If you provide an `init()` method, it will override the
      // `constructor()` method!
      constructor: function () {
        videoJsButtonClass.call(this, player);
      }, // notice the comma

      handleClick: function () {
        const videoElement = document.getElementsByTagName('video')[0];
        videoElement.webkitShowPlaybackTargetPicker();
      }
    });

    var concreteButtonInstance = player.controlBar.addChild(new concreteButtonClass());
    concreteButtonInstance.addClass("vjs-airplay");
  }

});