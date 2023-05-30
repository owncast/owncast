# Owncast on Windows

> Note: Owncast currently **does not natively support the Windows Operating System**, however it is possible to run Owncast on Windows using the Windows Subsystem for Linux (WSL2).

This document is a user-contributed document and the Owncast project does not actively maintain Windows support. Hopefully this can be helpful in pointing people in the right direction.

This document list out the steps in detail to install and run Owncast in Windows using Windows Subsystem for Linux, specifically **WSL2**.

Below are steps both for local development, contributing to the project and running it in production.

---

## Required: Installing WSL2 in Windows

There are lots of tutorials available online (videos and docs both) on how to install WSL2.  
Here are the official documents from Microsoft -> [Install Linux on Windows with WSL](https://learn.microsoft.com/en-us/windows/wsl/setup/environment)
Some points to remember ->

- Preferable method to install WSL2 is by using the `wsl --install `. If you are facing issues with this method you can look at - [Manual installation steps for older versions of WSL](https://learn.microsoft.com/en-us/windows/wsl/install-manual)
- Make sure you have enabled the Virtual Machine feature. (ignore if used wsl --install method)
- Make sure you have WSL2
- Installed your Linux distribution of choice and make sure you installed the latest available version (Preferably Ubuntu)

### Setting up WSL2 and the distribution of your choice

After basic setup, you can look into setting WSL2 for development. Here is the link for a detailed document by Microsoft - [https://learn.microsoft.com/en-us/windows/wsl/setup/environment](https://learn.microsoft.com/en-us/windows/wsl/setup/environment)

---

## Installing Owncast under WSL2

Once you're running WSL2 in Windows you can install Owncast the same way you would on any Linux distribution by following the [Quickstart](https://owncast.online/quickstart/) guide.

## Contributing to Owncast by performing local development

If you want to use your Windows machine to contribute to Owncast, you'll need to do so under WSL2 and make sure the following prerequisites are installed.

### Make sure all the prerequisites are installed in WSL2

Here is the list for all the prerequisites required ->

- C compiler, such as [GCC compiler](https://gcc.gnu.org/install/download.html) or a [Musl-compatible compiler](https://musl.libc.org/)
- npm (Node Package Manager) is installed as `sudo apt install npm`.
- Node.js is installed (LTS Version) `sudo apt install nodejs`.
- [ffmpeg](https://ffmpeg.org/download.html)
- Install the [Go toolchain](https://golang.org/dl/) (1.20 or above).

### Read more

Once your local development environment is setup, you can read more about how to contribute to Owncast [by reading the development document](https://owncast.online/development/).

## Some possible issues you can face while setting up WSL2

### You have an older version of Nodejs installed in the WSL2

To solve this issue you can look at nvm. Here is one tutorial - [Node-Version-Manager](https://www.digitalocean.com/community/tutorials/how-to-install-node-js-on-ubuntu-20-04#option-3-installing-node-using-the-node-version-manager).

### The broadcasting Software failed to connect to the server

This issue arises when you try to use `rtmp://localhost:1935/live` for example in OBS.
To solve this issue you need to find the correct IP address for the WSL2 you are running and use that instead of localhost.
You can use the below commands to find that ->  
Note: you can use either of these, whichever works for you.

- In WSL2 Terminal -
  `ip addr show eth0 | grep -oP '(?<=inet\s)\d+(\.\d+){3}'`
- In Windows Terminal -
  `wsl -- ip -o -4 -json addr list eth0`
  In this result look for "local": X.X.X.X

After finding the IP address in your broadcasting software make the server point to  
`rtmp://<your version of IP address>:1935/live`

Example in OBS-Studio ->
![image](https://user-images.githubusercontent.com/73140257/228762798-a0c56695-c692-4295-b11b-f2e85e867ce7.png)

## More resources

- [Windows Subsystem for Linux](https://learn.microsoft.com/en-us/windows/wsl/)
- [Owncast development documentation](https://owncast.online/development/)
- [Owncast quickstart guide](https://owncast.online/quickstart/)
- [Owncast README](https://github.com/owncast/owncast/blob/develop/README.md#building-from-source)
