import { ConnectedClientInfoEvent } from '../../../interfaces/socket-events';

export default function handleConnectedClientInfoMessage(
  message: ConnectedClientInfoEvent,
  setChatDisplayName: (string) => void,
  setChatUserId: (number) => void,
  setIsChatModerator: (boolean) => void,
) {
  const { user } = message;
  const { id, displayName, scopes } = user;
  setChatDisplayName(displayName);
  setChatUserId(id);
  setIsChatModerator(scopes.includes('moderator'));
}
