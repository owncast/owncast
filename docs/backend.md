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

subgraph ChatService[fa:fa-comment Chat\nService]
end

subgraph Dependencies

    App{AppController}
    ChatService--->App
    Webhooks--->App
    ActivityPubOutboundHandlers[fa:fa-hashtag ActivityPub\nOutbound]--->App

    ConfigRepository(fa:fa-hard-drive Config\nRepository)--->App
    UserRepository(fa:fa-hard-drive User\nRepository)--->App
    NotificationsRepository(fa:fa-hard-drive Notifications\nRepository)--->App
    ChatRepository(fa:fa-hard-drive Chat\nRepository)
    APRepository(fa:fa-hard-drive ActivityPub\nRepository)

    Database(fa:fa-hard-drive Database)--->ConfigRepository
    Database--->UserRepository
    Database--->APRepository
    Database--->NotificationsRepository

    ChatRepository-->ChatService

    ApplicationState(fa:fa-list Application\nState)--->App
    GeoIP(fa:fa-globe GeoIP\nLookup)--->App
    Statistics(fa:fa-list Statistics)--->App
    OutboundWebhooks[Outbound\nWebhooks]--->App

end

subgraph VideoStorageProviders[Video Storage Providers]
    LocalStorage((fa:fa-hard-drive Local\nStorage))
    S3Storage((fa:fa-hard-drive S3\nStorage))
end

subgraph Authentication[Chat Authentication]
    IndieAuth[fa:fa-key IndieAuth]
    FediAuth[fa:fa-key FediAuth]
end

subgraph Notifications[External Notifications]
    DiscordNotifier[fa:fa-comment Discord]
    BrowserNotifier[fa:fa-browser Browser]
end

subgraph WebServer[Web Server]
    ActivityPubHandlers[fa:fa-file ActivityPub\nHandlers]

    subgraph WebAssets[Web Assets]
        EmbeddedStaticFiles((fa:fa-file Embedded\nStatic Assets))
        OnDiskStaticFiles((fa:fa-file On Disk\nStatic Assets))
        WebApplication[fa:fa-file  Web\nApplication]
        PublicFiles[fa:fa-file Public\nDirectory]
    end

    subgraph HTTPHandlers[fa:fa-browser HTTP Handlers]
        AdminAPIs[Admin\nAPIs]
        ThirdPartyAPIs[3rd Party\nAPIs]
        WebSocket[WebSocket]
        subgraph ChatAPIs[Chat\nAPIs]
            ChatUserRegistration[Chat User\nRegistration]
            Emoji[Emojis]
            subgraph ChatAuthAPIs[Chat\nAuthentication]
                FediAuth
                IndieAuth
            end
        end
        subgraph VideoAPIs[fa:fa-video Video APIs]
            ViewerPing[Viewer\nPing]
            PlaybackMetrics[Playback\nHealth Metrics]
        end
        ActivityPubHandlers
        ApplicationConfig[Application\nConfig]
        ApplicationStatus[Application\nStatus]
        Directory[Directory\nAPI]
        Followers[Followers]
    end
end

subgraph Streamer
    BroadcastingSoftware>fa:fa-video Broadcasting\nSoftware]
end

%% All the services and packages require access
%% to dependencies through a Application reference.
App-.->HTTPHandlers
App-.->VideoPipeline
App-.->ActivityPub
App-.->Authentication
App-.->Notifications
App-.->DirectoryNotifier[Owncast\nDirectory]

LocalStorage--HLS-->OnDiskStaticFiles

RTMPService>RTMP Ingest]--RTMP-->VideoTranscoder
VideoTranscoder--HLS-->VideoStorageProviders

%% Streamers
BroadcastingSoftware--RTMP-->RTMPService

%% Style the nodes

%% Define reusable styles for node types
classDef bigtext font-weight:bold,font-size:40px
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
%%{init: {'theme':'base', 'themeVariables': {'darkMode': true, 'lineColor': '#c3dafe', 'tertiaryTextColor': 'white', 'clusterBkg': '#2d3748', 'primaryTextColor': '#39373d', "edgeLabelBackground": "white", "fontFamily": "monospace", "fontSize": "30px"}}}%%
```
