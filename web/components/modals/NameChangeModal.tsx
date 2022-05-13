import { useRecoilValue } from 'recoil';
import WebsocketService from '../../services/websocket-service';
// import { setLocalStorage } from '../../utils/helpers';
import { websocketServiceAtom } from '../stores/ClientConfigStore';

/* eslint-disable @typescript-eslint/no-unused-vars */
interface Props {}

export default function NameChangeModal(props: Props) {
  const websocketService = useRecoilValue<WebsocketService>(websocketServiceAtom);

  // const handleNameChange = () => {
  //   // Send name change
  //   const nameChange = {
  //     type: SOCKET_MESSAGE_TYPES.NAME_CHANGE,
  //     newName,
  //   };
  //   websocketService.send(nameChange);

  //   // Store it locally
  //   setLocalStorage(KEY_USERNAME, newName);
  // };

  return <div>Name change modal component goes here</div>;
}
