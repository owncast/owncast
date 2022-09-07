import { ConnectedClientInfoEvent } from '../../../interfaces/socket-events';

export function handleConnectedClientInfoMessage(
  message: ConnectedClientInfoEvent,
  setChatDisplayName: (string) => void,
  setChatDisplayColor: (number) => void,
  setChatUserId: (number) => void,
  setIsChatModerator: (boolean) => void,
  setChatAuthenticated: (boolean) => void,
) {
  const { user } = message;
  const { id, displayName, displayColor, scopes, authenticated } = user;
  setChatDisplayName(displayName);
  setChatDisplayColor(displayColor);
  setChatUserId(id);
  setIsChatModerator(scopes?.includes('MODERATOR'));
  setChatAuthenticated(authenticated);
}
export default handleConnectedClientInfoMessage;
