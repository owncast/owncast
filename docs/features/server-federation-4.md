# Following / Federating with Other Owncast Servers: Phase 4 Sending ActivityPub stream ended status update

This document outlines the specifications for implementing the backend and API changes required to support outbound server federation status updates in Owncast using ActivityPub. This feature will allow Owncast servers to follow and display the status of other Owncast servers.

## Goal

The goal of this feature is to allow an Owncast server to federate status with other Owncast servers. This allows a network of Owncast servers to share information about their status, such as whether they are online or offline, and to share information about their streams, such as the title and description of the stream.

The specific goal of this change is to send an ActivityPub "Leave" activity to all federated Owncast servers when the local server's live stream ends.

The recieving server should update the status of the local server to offline when it receives this activity, and update the server metadata accordingly.

### ActivityPub Support

- Accept an ActivityPub "Leave" activity from a federated Owncast server when its live stream ends. https://www.w3.org/TR/activitystreams-vocabulary/#dfn-leave
- Send an ActivityPub "Leave" activity from the local Owncast server when its live stream ends. https://www.w3.org/TR/activitystreams-vocabulary/#dfn-leave
- In both cases the custom server metadata should be included.
- In both cases expect the actor to be the Owncast server itself.
- In both cases expect the object to be the URL of the federated server.
- Use the existing ActivityPub infrastructure in Owncast to handle the sending, delivery and receipt of these activities via the inbox, outbox and worker queues.

### Database Changes

- No database changes are required for this functionality.

## Requirements
- For each federated server exchange the following information:
	- Server logo
  - Server name
  - Stream title
  - Stream description
  - Tags associated with the stream
