export enum AppState {
  AppLoading,
  ChatLoading,
  Loading,
}

export enum ChatVisibilityState {
  Hidden, // The chat is available but the user has hidden it
  Visible, // The chat is available and visible
}

export enum ChatState {
  Available, // Normal state
  NotAvailable, // Chat features are not available
  Loading, // Chat is connecting and loading history
  Offline, // Chat is offline/disconnected for some reason
}
