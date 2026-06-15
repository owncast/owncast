# Following / Federating with Other Owncast Servers: Phase 1 UI

This document outlines the specifications for implementing the user interface changes required to support server federation in Owncast. This feature will allow Owncast servers to follow and display the status of other Owncast servers.

This allows visitors to an Owncast server to see a list of other Owncast servers that are being followed, along with their current status (online/offline) and stream information if they are online. This creates a network effect, allowing users to discover and navigate to other Owncast servers easily. It also gives somewhere for somebody to go if the server they are on is offline.

The user interface should give the feel of browsing through Twitch or YouTube, with each followed server represented by a card that includes a thumbnail, server name, stream title, description, and tags.

This document focuses solely on the user interface changes required to support server federation. The backend and API changes required to support this feature will be added separately.

## Goal

The goal of this feature is to allow an Owncast server to federate status with other Owncast servers. This allows a network of Owncast servers to share information about their status, such as whether they are online or offline, and to share information about their streams, such as the title and description of the stream.

## Requirements

### Web User Interface

- The web user interface changes should be done using Ant Design v4 components, not Ant Design v5.
- Styling should be done using CSS modules and follow the existing styling conventions in the Owncast codebase. This means using sass scss files.
- Any new components should be added to Storybook for testing and prototyping.

### Primary user-facing interface changes

- Add a new tab to the main web interface called "Streams".
- If the tab bar is currently hidden because there were previously no tabs to show, it should become visible when the "Streams" tab is added.
- The "Streams" tab should display a list of all federated Owncast servers and their current status (online/offline).
- Each Owncast server should be represented by a card in the style of Twitch or YouTube with a thumbnail on top and text information below.
- For each online server, display the following information:
	- Server logo
  - Server name
  - Stream title
  - Stream description
  - Tags associated with the stream
  - A thumbnail image of the current stream (if available)
- For each offline server, display the following information:
	- Server logo
  - Server name
  - Status: Offline
	- Tags
- Allow users to click on a server card to navigate to the server's URL.

### Admin Interface

- Add a new admin section under "Social" for following servers.
- The admin interface should allow administrators to add a federated servers.
- Administrators should be able to remove federated servers from the list.
- Validate the URL of the federated server when adding it to ensure it is a valid URL.
- Display the current status of each federated server (online/offline) in the admin interface.
