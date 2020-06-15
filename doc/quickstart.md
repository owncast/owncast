## Quickstart

These steps will utilize Docker, as going from a brand new sever to running the service is easiest done within a container.

1. Create or login to an existing server somewhere.
1. Install Docker.  If it's a Debian based linux machine `sudo apt-get install docker.io`.
1. Install git `sudo apt-get install git`.
1. Download the code.  `git clone https://github.com/gabek/owncast`
1. Copy `config/config-example.yaml` to `config/config.yaml`
1. Edit `config/config.yaml` and change the path of ffmpeg to `/usr/bin/ffmpeg`
1. Edit your stream key to whatever you'd like it to be in the config.
1. Make any other config changes.
1. Run `docker build -t owncast .` and wait.  It may take a few minutes to build depending on the speed of your server.
1. Run `docker run -p 8080:8080 -p 1935:1935 -it owncast` to start the service.
1. Point your broadcasting software at your new server using `rtmp://yourserver/live` and the stream key you set above and start your stream.
1. Access your server in your web browser by visiting `http://yourserver:8080`.
1. If you ever make any future config file changes you must rerun the `docker build` step otherwise you can just run the `docker run` step to run the service going forward.
