<br />
<p align="center">
  <a href="https://github.com/owncast/owncast" alt="Owncast">
    <img src="https://owncast.online/images/logo.png" alt="Logo" width="200">
  </a>
</p>

<br/>

  <p align="center">
    <strong>Take control over your content and stream it yourself.</strong>
    <br />
    <a href="https://owncast.online"><strong>Explore the docs »</strong></a>
    <br />
    <a href="https://watch.owncast.online/">View Demo</a>
    ·
    <a href="https://broadcast.owncast.online/">Use Our Server for Testing</a>
    ·
    <a href="https://owncast.online/faq/">FAQ</a>
    ·
    <a href="https://github.com/owncast/owncast/issues">Report Bug</a>
  </p>
</p>

# Owncast-for-windows

> Note: Owncast currently **does not support Windows Operating System**. This is a workaround to run Owncast on Windows.

This document list out the steps in detail to install and run Owncast in Windows using Windows Subsystem for Linux, specifically **WSL2**.  
The steps provided are for both the use cases whether you want to use Owncast or build and install Owncast for contributing to the project.

---

### Installing WSL2 in Windows.

There are lots of tutorials available online (videos and docs both) on how to install WSL2.  
Here are the official documents from Microsoft -> [Install Linux on Windows with WSL](https://learn.microsoft.com/en-us/windows/wsl/setup/environment)
Some points to remember ->

- Preferable method to install WSL2 is by using the `wsl --install `. If you are facing issues with this method you can look at - [Manual installation steps for older versions of WSL](https://learn.microsoft.com/en-us/windows/wsl/install-manual)
- Make sure you have enabled the Virtual Machine feature. (ignore if used wsl --install method)
- Make sure you have WSL2
- Installed your Linux distribution of choice and make sure you installed the latest available version (Preferably Ubuntu)

### Setting up WSL2 and the distribution of your choice.

After basic setup, you can look into setting WSL2 for development. Here is the link for a detailed document by Microsoft - [https://learn.microsoft.com/en-us/windows/wsl/setup/environment](https://learn.microsoft.com/en-us/windows/wsl/setup/environment)

<hr>

## Now for the users who just want to use Owncast

### Installing Owncast for streaming.

One of the many ways to install Owncast is to use the Quick Installer script.
In the WSL2 terminal run, the following command

```
curl -s https://owncast.online/install.sh | bash
cd owncast
./owncast
```

If everything is installed correctly you will be able to access the admin interface at port `localhost:8080/admin`.  
Here are other ways listed you can use to install Owncast -[Owncast Quickstart](https://owncast.online/quickstart/)
If you face any errors you can look at the Error List at the end of the document.

<hr>

## Next steps are for developers and users looking to build and install Owncast.

### Make sure all the prerequisites are installed in WSL2.

Here is the list for all the prerequisites required ->

- C compiler, such as [GCC compiler](https://gcc.gnu.org/install/download.html) or a [Musl-compatible compiler](https://musl.libc.org/)
- npm (Node Package Manager) is installed as `sudo apt install npm`.
- Node.js is installed (LTS Version) `sudo apt install nodejs`.
- [ffmpeg](https://ffmpeg.org/download.html)
- Install the [Go toolchain](https://golang.org/dl/) (1.18 or above).

---

These steps are for building and install Owncast on Windows for development and contributing to this project.

### Installing Owncast Backend

The Owncast backend is a service written in Go.

1. Clone the repo. `git clone https://github.com/owncast/owncast`
1. `go run main.go` will run from the source.
1. Visit `http://yourserver:8080` to access the web interface or `http://yourserver:8080/admin` to access the admin.
1. Point your [broadcasting software](https://owncast.online/docs/broadcasting/) at your new server and start streaming.

### Installing Owncast Frontend

The frontend is the web interface that includes the player, chat, embed components, and other UI.

1. This project lives in the `web` directory.
1. Run `npm install` to install the Javascript dependencies.
1. Run `npm run dev`

Further information about the development is [here](https://owncast.online/development/)

---

#### Congratulations !!! If you have reached here. You have successfully installed and run Owncast locally on Windows. Happy Streaming.

You can contribute to Owncast here [https://github.com/owncast/owncast](https://github.com/owncast/owncast)

<hr>
### Some errors you can face while following these steps.

#### Error 1.

**You have an older version of Nodejs installed in the WSL2.**  
To solve this issue you can look at nvm. Here is one tutorial - [Node-Version-Manager](https://www.digitalocean.com/community/tutorials/how-to-install-node-js-on-ubuntu-20-04#option-3-installing-node-using-the-node-version-manager).

<hr>

#### Error 2.

**The broadcasting Software failed to connect to the server.**  
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

Now check again, you will be able to stream.

---

<!-- CONTACT -->

## Contact

Thanks for installing Owncast if you face any problems. You can contact here.

Project chat: [Join us on Rocket.Chat](https://owncast.rocket.chat/home) if you want to contribute, follow along, or if you have questions.

Gabe Kangas - [@gabek@social.gabekangas.com](https://social.gabekangas.com/gabek) - email [gabek@real-ity.com](mailto:gabek@real-ity.com)

Project Link: [https://github.com/owncast/owncast](https://github.com/owncast/owncast)
