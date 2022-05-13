export enum AppState {
  Loading, // Initial loading state as config + status is loading.
  Registering, // Creating a default anonymous chat account.
  Online, // Stream is active.
  Offline, // Stream is not active.
  OfflineWaiting, // Period of time after going offline chat is still available.
  Banned, // Certain features are disabled for this single user.
}

export enum ChatVisibilityState {
  Hidden, // The chat components are not available to the user.
  Visible, // The chat components are not available to the user visually.
}

export enum ChatState {
  Available = 'Available', // Normal state. Chat can be visible and used.
  NotAvailable = 'NotAvailable', // Chat features are not available.
  Loading = 'Loading', // Chat is connecting and loading history.
  Offline = 'Offline', // Chat is offline/disconnected for some reason but is visible.
}

export enum VideoState {
  Available, // Play button should be visible and the user can begin playback.
  Unavailable, // Play button not be visible and video is not available.
  Playing, // Playback is taking place and the play button should not be shown.
}

export function getChatState(state: AppState): ChatState {
  switch (state) {
    case AppState.Loading:
      return ChatState.NotAvailable;
    case AppState.Banned:
      return ChatState.NotAvailable;
    case AppState.Online:
      return ChatState.Available;
    case AppState.Offline:
      return ChatState.NotAvailable;
    case AppState.OfflineWaiting:
      return ChatState.Available;
    case AppState.Registering:
      return ChatState.Loading;
    default:
      return ChatState.Offline;
  }
}

export function getChatVisibilityState(state: AppState): ChatVisibilityState {
  switch (state) {
    case AppState.Loading:
      return ChatVisibilityState.Hidden;
    case AppState.Banned:
      return ChatVisibilityState.Hidden;
    case AppState.Online:
      return ChatVisibilityState.Visible;
    case AppState.Offline:
      return ChatVisibilityState.Hidden;
    case AppState.OfflineWaiting:
      return ChatVisibilityState.Visible;
    case AppState.Registering:
      return ChatVisibilityState.Visible;
    default:
      return ChatVisibilityState.Hidden;
  }
}
