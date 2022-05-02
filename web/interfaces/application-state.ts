export enum AppState {
  Loading, // Initial loading state as config + status is loading.
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
  Available, // Normal state. Chat can be visible and used.
  NotAvailable, // Chat features are not available.
  Loading, // Chat is connecting and loading history.
  Offline, // Chat is offline/disconnected for some reason but is visible.
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
    default:
      return ChatVisibilityState.Hidden;
  }
}
