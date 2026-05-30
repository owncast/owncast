# Proxmox Installation

This guide assumes you have proxmox CE 9.1.1+ already installed on a host system of your choosing. For installing proxmox [follow this guide here](https://proxmox.com/en/products/proxmox-virtual-environment/get-started). The host os is Debian 13 (Trixie). There is a yunohost script for Debian 12 Bookworm that also works nicely. 

### Getting Started
1. Log into proxmox pve host (root user)
2. Run this command as root user in pve shell
`var_cpu="4" var_ram="4096" var_disk="20" bash -c "$(curl -fsSL https://raw.githubusercontent.com/community-scripts/ProxmoxVE/main/ct/owncast.sh)"`
  Explanation of var's
  - var_cpu="4": owncast runs fine on 2 CPU cores. Consider more resources for multiple transcodes and resolutions beyond 1080p30
  - var_ram="4096": 4GB RAM this is much more than necessary. Default is 2048. 
  - var_disk="20": 20GB storage **This is the most important var here**. The script default assumes non-self contained storage setup (Advanced install). 
3. Select 'advanced install' if you wish to rename your instance. Default LXC name is 'owncast' It is recommended to change the container name if planning on running more than 1 owncast stream in proxmox.
4. follow the installation script, Defaults are fine. Depending on your specific setup, you may need to reconfigure the network settings, proxy or additional resources to get results. 
5. ***MAKE SURE YOU PAY ATTENTION HERE*** the default nvidia graphics driver can be mismatched if you use `apt` or likely any other package manager. DO NOT AUTOMATICALLY INSTALL graphics drivers via the owncast script. it will make manual install difficul later on. 
6. complete install script.
7. Verify owncast is running. your GPU will not be functioning as of this point. 
  Some recommendations before going onto the next steps
  - Ensure you can connect your owncast instance to the OBS stream. A 2 minute test stream can save you headache down the line
  - Check logs. Sometimes there can be package misconfigs. owncast logs will catch this. 

#### Information to note

In this configuration, most of the owncast data lives in /opt/owncast/*. The following is a tree diagram of the default install.

```
.
├── data
│   ├── backup
│   │   └── owncastdb.bak
│   ├── emoji
│   │   ├── blob
│   │   │   ├── ablobattention.gif
│   │   │   ├── ...
│   │   ├── conigliolo96
│   │   │   ├── conigliolo15.gif
│   │   │   ├── ...
│   │   ├── dog
│   │   │   ├── img001.svg
│   │   │   ├── ...
│   │   ├── mutant
│   │   │   ├── 8_ball.svg
│   │   │   ├── ...
│   │   └── thanks.png
│   ├── hls
│   │   ├── 0
│   │   │   ├── stream.m3u8
│   │   │   ├── stream-offline-0.ts
│   │   │   └── stream-offline-1.ts
│   │   └── stream.m3u8
│   ├── logo.png
│   ├── logs
│   │   ├── owncast.log -> owncast.log.yyyymmddhhmmss
│   │   ├── owncast.log.yyyymmddhhmmss
│   │   └── transcoder.log
│   ├── metrics
│   │   ├── p-xxxxxxxxxx-yyyyyyyy
│   │   │   ├── data
│   │   │   └── meta.json
│   │   ├── p-.........
│   │   └── wal
│   │       ├── 8
│   │       └── 9
│   ├── owncast.db
│   ├── owncast.db-shm
│   ├── owncast.db-wal
│   └── tmp
│       ├── offline-v2.tsxxxxxxxxxx
│       ├── offline-v2.tsxxxxxxxxxx
│       ├── offline-v2.tsxxxxxxxxxx
│       ├── preview.gif
│       └── thumbnail.jpg
└── owncast
```
Make sure you check the /logs/ folder for your most recent test stream. It will show you any configuration errors that may crash stream in the future.

### Enable GPU Passthrough (NVIDIA)

#### Attaching in Proxmox
The owncast installation script handles this automatically. If you did not see this happen, it is easier to re-run the installation script instead of manually attach the GPU. It may be necessary for certain hardware/software setups. 

#### Driver Setup (NVIDIA)
1. ensure ffmpeg, make, gcc, and cuda dependencies are installed [see this link](https://developer.nvidia.com/blog/nvidia-ffmpeg-transcoding-guide/) you may need to check pkg-config as well. also double check that iommu is enabled on pve host
  - If you run into hangups here, you may need to enable the 'non-free' debian repositories. 
  - Installation order tested:
    1. gcc
    2. ffmpeg
    3. make
    4. cuda
2. apt update, apt upgrade, reboot. (make sure your host and container system is up to date)
3. navigate to nvidia's driver repo [link to driver search](https://www.nvidia.com/en-us/drivers/)
4. select the driver for your OS + hardware combo. ensure compatibility, cant help you there. 
5. copy download link to repo -> navigate to owncast shell -> `wget https://www.nvidia.com/en-us/drivers/details/XXXXX/`
6. `chmod +x NVIDIA-*.run` the resulting nvidia driver file after downloading
7. run the driver installer script. 
8. you may get some errors, I used [this guide here to install properly](https://github.com/gma1n/LXC-JellyFin-GPU) if you do run into errors. make sure you have your dependencies configured pre-install. screwing this up can mean having to reinstall the LXC from script (step 1 of 1)
9. ensure `nvidia-smi` is working and the driver versions are identical to PVE host. I did this on a single PC node with 1 GPU. attaching the GPU may need to be done manually on multi-PC or multi-GPU setups. that can be done following a LXC GPU passthrough guide.
  When its working, should look something like this:
  ```
  +-----------------------------------------------------------------------------------------+
  | NVIDIA-SMI xxx.yyy                Driver Version: xxx.yyy        CUDA Version: 13.0     |
  +-----------------------------------------+------------------------+----------------------+
  | GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
  | Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
  |                                         |                        |               MIG M. |
  |=========================================+========================+======================|
  |   0  NVIDIA GeForce GTX 1660 Ti     Off |   00000000:01:00.0 Off |                  N/A |
  | 41%   45C    P0             15W /  120W |       0MiB /   6144MiB |      0%      Default |
  |                                         |                        |                  N/A |
  +-----------------------------------------+------------------------+----------------------+
  
  +-----------------------------------------------------------------------------------------+
  | Processes:                                                                              |
  |  GPU   GI   CI              PID   Type   Process name                        GPU Memory |
  |        ID   ID                                                               Usage      |
  |=========================================================================================|
  |  No running processes found                                                             |
  +-----------------------------------------------------------------------------------------+
  ```

**Congratulations! You have installed owncast on proxmox**

### Enable Hardware Acceleration

This part is the easy part. 
1. (If you havent already) Log into your owncast/admin interface 
  Owncast has great guides for each of these found [here in the documentation](https://owncast.online/docs/)
  - Set a secure login and password
  - Set a secure streamkey (use sha256 keygen, every single computer has it) `echo -n foobar | sha256sum` more secure methods are still recommended. If I can use default password to login to your instance I will do so and rickroll your viewers. You've been warned.
  - Setup HTTPS and SSL keys. Caddy is recommended for owncast. Viewers won't be able to watch unsigned .m3u8 files directly in browser. 
2. Navigate to >configuration>Video. There will be your encoder menu as well as latency buffer visible. Underneath latency buffer there's _'Advanced Settings'_ drop down. Clicking that reveals your video codecs. Select the codec that matches your hardware (NVIDIA GPU Acceleration). Only installed codecs can be seen.
  If Owncast is configured correctly, x264 will be shown by default. If none are visible, see #Common-Mistakes in this guide. 
3. Configure Additional Transcoding presets leveraging your GPU encoders and decoders! The rest should be handled by owncast
4. Test the stream wtih new encoders. Both locally, on wifi, and through cellular for best results. Have fun tuning! 


## Common Issues
*To add to section this create a pull request or reach out in the community discussion board*

- Q: My Owncast install has a CUDA error in the logs
  A: You must ensure CUDA drivers are installed properly for NVIDIA cards to work via passthrough see [CUDA installation guide](https://docs.nvidia.com/cuda/cuda-installation-guide-linux)
- Q: My Owncast stream spams 'disconnected' status overwhelming my Streaming PC's notifications. 
  A: There exists a encoder mismatch between your exported encoder and Owncast's decoder. The hardware on both ends must be compatible. Only NVIDIA<->NVIDIA Interfaces have been verified at this time. 
- Q: How can I increase beyond 1080p30? 
  A: Edit the 'advanced' settings under the video encoder presets. It is recommended to have at least 1 resolution <1080p30 for mobile/low data speed viewers. Beyond 1080p30 resolution significant bitrate is required for visual fidelity. This is only recommended on newer hardware on high speed (1GB+) WAN connections. It is possible on <300 MB/s, but not without compromise in other services. 
- Q: When running `nvidia-smi` I get an error saying that the hardware cannot communicate with the nvidia drivers.
  There are multiple issues that can cause this
  - A: Driver version mismatch between host-os and owncast-lxc. Ensure you installed the same GPU driver version that runs on your host OS. upgrading to a new driver within owncast may currently require a reinstall of the container. Can do anything if you are good enough with the terminal though.
  - A: Too Many driver versions installed. Try either reinstalling the owncast container with a manual GPU driver install, or remove all NVIDIA drivers from container and reinstall the correct version. Both worked in testing. 





Guide written by hefftv in 2026