import { ConnectedClientInfoEvent } from '../../../interfaces/socket-events';

export default function handleConnectedClientInfoMessage(
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
  setIsChatModerator(scopes?.includes('moderator'));
  setChatAuthenticated(authenticated);
}
