## systemd user service containerized owncast

[Podman](https://github.com/containers/podman/) version 3.4.0 (released September 2021)
supports running OCI containers (i.e. Docker containers) with _socket activation_.
This opens up the possibility to run an owncast container with _socket activation_ in
a systemd user service.

### Install

1. Clone this repository
    ```
    git clone URL
    ```

2. Install the systemd files
    ```
    mkdir -p ~/.config/systemd/user
    cp -r owncast/examples/systemd-with-container/owncast-*@* ~/.config/systemd/user
    ```

3. Optional: If __container-linux__ >= v2.179.0
   it's possible to remove `--security-opt label=disable`

   ```
   sed '/--security-opt label=disable/d' ~/.config/systemd/user/owncast@.service
   ```

4. Reload systemd

   ```
   systemctl --user daemon-reload
   ```

### Usage

1. Start the target _owncast@tcpdemo.target_
    ```
    systemctl --user start owncast@tcpdemo.target
    ```

2. Open  _http://localhost:8080_ in a web browser
