import { ConnectedClientInfoEvent } from '../../../interfaces/socket-events';

export default function handleConnectedClientInfoMessage(
  message: ConnectedClientInfoEvent,
  setChatDisplayName: (string) => void,
) {
  const { user } = message;
  const { displayName } = user;
  setChatDisplayName(displayName);
}
