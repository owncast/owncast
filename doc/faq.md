<!-- vscode-markdown-toc -->
* [Who is this for?](#Whoisthisfor)
* [Why would I use this instead of Twitch, Facebook Live, or YouTube live?](#WhywouldIusethisinsteadofTwitchFacebookLiveorYouTubelive)
* [What can I customize?](#WhatcanIcustomize)
* [What kind of server do I need to run Owncast?](#WhatkindofserverdoIneedtorunOwncast)
* [When would I need a more powerful server?](#WhenwouldIneedamorepowerfulserver)
* [What are bitrates?  Why would I want more?](#WhatarebitratesWhywouldIwantmore)

<!-- vscode-markdown-toc-config
	numbering=true
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->


##  1. <a name='Whoisthisfor'></a>Who is this for?

Owncast is for people who are live streamers, or who wants to host live streams for others.


##  2. <a name='WhywouldIusethisinsteadofTwitchFacebookLiveorYouTubelive'></a>Why would I use this instead of Twitch, Facebook Live, or YouTube live?

Owncast might be a good alternative if you're somebody who doesn't want to rely on the large companies or wants the ability to build something completely custom that is more in tune with the experience they want to offer.  As a bonus it allows you to offer a live streaming experience that is without tracking, invasive analytics or advertising.


##  3. <a name='WhatcanIcustomize'></a>What can I customize?
You can edit the included config file to specify your site name, logo and social networking links.

Additionally, out of the box there is a fully functional web site with built-in chat and a video player.  It's HTML + CSS + Javascript that you can edit directly.  It's yours.  You could also disable included web interface completely and instead embed your stream into your existing web site.  Build something cool!


##  4. <a name='WhatkindofserverdoIneedtorunOwncast'></a>What kind of server do I need to run Owncast?

You need a publicly accessible Linux or macOS server on the internet.  Something like [Linode](https://www.linode.com/products/shared/) or [Digital Ocean](https://www.digitalocean.com/products/droplets/) are good options and start at $5/mo.


##  5. <a name='WhenwouldIneedamorepowerfulserver'></a>When would I need a more powerful server?

The more bitrates you support the more processing power is required.  You can easily run three bitrates on something like a $10/mo dedicated server.


##  6. <a name='WhatarebitratesWhywouldIwantmore'></a>What are bitrates?  Why would I want more?

Bitrates specify the quality of the video.  The more bitrates you support the wider range of network conditions you can support.  For example, a user on their broadband connection at home would want the full quality you have available.  But if they're on a slow wireless connection on their phone a lower bitrate would result in less buffering and a smoother experience.  This would be a two bitrate configuration and allow for offering two distinct video qualities to your users.

[Read more on Wikipedia](https://en.wikipedia.org/wiki/Adaptive_bitrate_streaming) about adaptive bitrate streaming.

