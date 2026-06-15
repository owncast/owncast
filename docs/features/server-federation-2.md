# Following / Federating with Other Owncast Servers: Phase 2 Accepting ActivityPub Status Updates

This document outlines the specifications for implementing the backend and API changes required to support inbound server federation status updates in Owncast using ActivityPub. This feature will allow Owncast servers to follow and display the status of other Owncast servers.

## Goal

The goal of this feature is to allow an Owncast server to federate status with other Owncast servers. This allows a network of Owncast servers to share information about their status, such as whether they are online or offline, and to share information about their streams, such as the title and description of the stream.

The specific goal of this change is

### ActivityPub Support

- Accept an ActivityPub "Arrive" activity from a federated Owncast server when its live stream comes online. https://www.w3.org/TR/activitystreams-vocabulary/#dfn-arrive
- Accept an ActivityPub "Leave" activity from a federated Owncast server when its live stream goes offline. https://www.w3.org/TR/activitystreams-vocabulary/#dfn-leave
- In both cases expect the actor to be the Owncast server itself.
- In both cases expect the object to be the URL of the federated server.
- Use the existing ActivityPub infrastructure in Owncast to handle the delivery and receipt of these activities via the inbox and worker queue.

### Database Changes

- Add a new table to store the list of federated Owncast servers that are being followed.
- Add fields to store the current status of each federated server (online/offline) and stream information (title, description, tags, thumbnail URL).
- Add timestamps for when the server was last seen online and when the last status update was received.
- Update the table when an "Arrive" or "Leave" activity is received.

## Requirements
- For each federated server exchange the following information:
	- Server logo
  - Server name
  - Stream title
  - Stream description
  - Tags associated with the stream
  - A thumbnail image of the current stream (if available)
