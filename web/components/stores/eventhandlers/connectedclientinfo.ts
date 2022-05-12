import { ConnectedClientInfoEvent } from '../../../interfaces/socket-events';

export default function handleConnectedClientInfoMessage(message: ConnectedClientInfoEvent) {
  console.log('connected client', message);
}
