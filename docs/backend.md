# Owncast Backend Architecture

This is a work in progress document detailing the future backend architecture of Owncast. It should be seen as a living document until a refactor of the backend is complete.

## Dependencies

Dependencies are services that are required across the application. This can be things like the chat service or a data repository for config values or user data.

Note: A better name that "dependencies" might be clearer. Perhaps "services" or "providers".

TODO: Have a complete list of dependencies.

### Data Repositories

The repository pattern provides a layer of abstraction between the application and the data store, allowing the application to interact with the data store through a well-defined interface, rather than directly accessing the data store. This helps to decouple the application from the data store.

TODO: List out the repositories and what they do.

Learn more about the [repository pattern](https://techinscribed.com/different-approaches-to-pass-database-connection-into-controllers-in-golang/).

### Application Controller

The `AppController` has references to all the dependencies and serves as an arbiter between consumers of these services and the services themselves.

A reference to the `AppController` is passed in to the all the core functionality in the application and each package would have its own interface that `AppController` implements. This can include getting access to dependency services like getting access to the chat service, getting access to the config repository values or knowing application state such as if a stream is live or how many viewers are watching via metrics.

TODO: Show examples of how the application is passed in to packages and how to reference dependencies through it.

## Diagram

```mermaid
%% Owncast backend architecture.
%% See https://jojozhuang.github.io/tutorial/mermaid-cheat-sheet/
%% for a cheat sheet on creating this diagram.
%% Paste this document in https://mermaid.live as a quick way to edit.

%% TODO: Add links between nodes and the actual code.

%% This is a graph style diagram, Top Down.
graph TD


%% Define Nodes and Subgraphs

subgraph VideoPipeline[Video Pipeline]
    VideoTranscoder(fa:fa-video Video Transcoder)
    RTMPService[fa:fa-video RTMP Service]
end

subgraph ChatService[fa:fa-comment Chat Service]
end

subgraph Dependencies
    OutboundWebhooks[Outbound]

    App{Application}
    ChatService--->App
    Webhooks--->App

    ConfigRepository(fa:fa-hard-drive Config Repository)--->App
    UserRepository(fa:fa-hard-drive User Repository)--->App
    APRepository(fa:fa-hard-drive ActivityPub Repository)--->App
    NotificationsRepository(fa:fa-hard-drive Notifications Repository)--->App
    ChatRepository(fa:fa-hard-drive Chat Repository)

    Database(fa:fa-hard-drive Database)--->ConfigRepository
    Database--->UserRepository
    Database--->APRepository
    Database--->NotificationsRepository

    ChatRepository-->ChatService

    ApplicationState(fa:fa-list Application State)--->App
    GeoIP(fa:fa-globe GeoIP Lookup)--->App
    Statistics(fa:fa-list Statistics)--->App
end

subgraph VideoStorageProviders[Video Storage Providers]
    LocalStorage((fa:fa-hard-drive Local Storage))
    S3Storage((fa:fa-hard-drive S3 Storage))
end

subgraph ActivityPub
    ActivityPubInboundHandlers[fa:fa-hashtag Inbound]
    ActivityPubOutboundHandlers[fa:fa-hashtag Outbound]
end

subgraph Authentication
    IndieAuth[fa:fa-key IndieAuth]
    FediAuth[fa:fa-key FediAuth]
end

subgraph Notifications
    DiscordNotifier[fa:fa-comment Discord]
    BrowserNotifier[fa:fa-browser Browser]
end

subgraph WebServer[Web Server]
    ActivityPubHandlers[fa:fa-file ActivityPub Handlers]

    subgraph WebAssets[Web Assets]
        EmbeddedStaticFiles((fa:fa-file Embedded\nStatic Assets))
        OnDiskStaticFiles((fa:fa-file On Disk\nStatic Assets))
        WebApplication[fa:fa-file  Web\nApplication]
        PublicFiles[fa:fa-file Public\nDirectory]
    end

    subgraph HTTPHandlers[fa:fa-browser HTTP Handlers]
        AdminAPIs[Admin APIs]
        ThirdPartyAPIs[3rd Party APIs]
        WebSocket[WebSocket]
        subgraph ChatAPIs[Chat APIs]
            ChatUserRegistration[Chat User Registration]
            Emoji[Emojis]
            subgraph ChatAuthAPIs[Chat Authentication]
                FediAuth
                IndieAuth
            end
        end
        subgraph VideoAPIs[fa:fa-video Video APIs]
            ViewerPing[Viewer Ping]
            PlaybackMetrics[Playback health metrics]
        end
        ActivityPubHandlers
        ApplicationConfig[Application Config]
        ApplicationStatus[Application Status]
        Directory[Directory API]
        Followers[Followers]
    end
end

subgraph Streamer
    BroadcastingSoftware>fa:fa-video BroadcastingSoftware]
end

%% All the services and packages require access
%% to dependencies through a Application reference.
App-.->HTTPHandlers
App-.->VideoPipeline
App-.->ActivityPub
App-.->Authentication
App-.Stream went\nonline.->Notifications
App-.->DirectoryNotifier[Directory Notifier]

LocalStorage--HLS-->OnDiskStaticFiles

RTMPService>RTMP Ingest]--RTMP-->VideoTranscoder
VideoTranscoder--HLS-->VideoStorageProviders

%% Viewers
VideoPlayer-->VideoStorageProviders
WebBrowser-->WebAssets

%% Streamers
BroadcastingSoftware--RTMP-->RTMPService

%% Style the nodes

%% Define reusable styles for node types
classDef bigtext font-weight:bold,font-size:20px
classDef repository fill:#4F625B,color:#fff
classDef webservice fill:#6082B6,color:#fff
classDef rtmp fill:#608200, color:#fff
classDef inbound fill:#6544e9,stoke:#fff,color:#fff
classDef outboundservice fill:#2386e2,stroke:green,color:#fff
classDef storage fill:#42bea6,color:#fff

%% Assign styles to nodes
class App bigtext
class UserRepository repository
class ConfigRepository repository
class APRepository repository
class NotificationsRepository repository
class ChatRepository repository

class HTTPHandlers inbound
class ActivityPubInboundHandlers inbound
class InboundWebhooks inbound

class ActivityPubOutboundHandlers outboundservice
class DiscordNotifier outboundservice
class BrowserNotifier outboundservice
class OutboundWebhooks outboundservice
class DirectoryNotifier outboundservice

class BroadcastingSoftware rtmp
class RTMPService rtmp

class Database storage
class LocalStorage storage
class S3Storage storage

%% Customize the theme styles
%%{init: {'theme':'base', 'themeVariables': {'darkMode': true, 'lineColor': '#c3dafe', 'tertiaryTextColor': 'white', 'clusterBkg': '#2d3748', 'primaryTextColor': '#39373d', "edgeLabelBackground": "white", "fontFamily": "monospace"}}}%%
```
