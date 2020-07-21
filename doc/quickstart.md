## Quickstart

There are two quick ways to get up and running, depending on your preference.  One is to simply download the service and run it, and the other is through Docker, if Docker is your thing.


## Download and run
1. Install [`ffmpeg`](https://ffmpeg.org/download.html) if you haven't.
1. Make a directory to run the service from, and download a release from https://github.com/gabek/owncast/releases into that directory.
1. Unzip the release's archive: `unzip owncast-linux-x.x.x.zip`.
1. [Edit `config.yaml` as detailed below](#configure).  Specifically your stream key and `ffmpeg` location.
1. Run `./owncast` to start the service.

---

or.....

## Through docker

1. Download the code: `git clone https://github.com/gabek/owncast`
1. Copy `config-example.yaml` to `config.yaml`
1. [Edit `config.yaml`](#configure) with a file editor of your choice and change the path of ffmpeg by appending `ffmpegPath: /usr/bin/ffmpeg` at the top level of the yaml.
1. Edit your stream key to whatever you'd like it to be in the config.
1. If you ever make any future config file changes you must rerun the `docker build` step otherwise you can just run the `docker run` step to run the service going forward.
1. Run `docker build -t owncast .` and wait.  It may take a few minutes to build depending on the speed of your server.
1. Run `docker run -p 8080:8080 -p 1935:1935 -it owncast` to start the service.

---

### Configure

1. Edit `config.yaml` and change the path of ffmpeg to where your copy is.
1. In this default configuration there will be a single video quality available, simply whatever is being sent to the server is being distributed to the viewers.  The video is also going to be distributed from the server running the service in this case.
1. Continue to edit the config file and customize with your own details, links and info.  See [More Configuration](configuration.md) to find additional ways to configure video quality.

### Test
1. Point your broadcasting software at your new server using `rtmp://yourserver/live` and the stream key you set above and start your stream.
1. Access your server in your web browser by visiting `http://yourserver:8080`.


### That's it!
