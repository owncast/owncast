import { useEffect } from 'react';
import { ReactElement } from 'react-markdown/lib/react-markdown';
import { atom, useRecoilState } from 'recoil';
import { ServerStatus, makeEmptyServerStatus } from '../../models/ServerStatus';
import ServerStatusService from '../../services/StatusService';

export const serverStatusState = atom({
  key: 'serverStatusState',
  default: makeEmptyServerStatus(),
});

export function ServerStatusStore(): ReactElement {
  const [, setServerStatus] = useRecoilState<ServerStatus>(serverStatusState);

  const updateServerStatus = async () => {
    try {
      const status = await ServerStatusService.getStatus();
      setServerStatus(status);
      return status;
    } catch (error) {
      console.error(`serverStatusState -> getStatus() ERROR: \n${error}`);
      return null;
    }
  };

  useEffect(() => {
    setInterval(() => {
      updateServerStatus();
    }, 5000);
    updateServerStatus();
  }, []);

  return null;
}
