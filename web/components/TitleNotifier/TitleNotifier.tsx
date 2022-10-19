/**
 * This component is responsible for updating the title of the page when
 * different state changes occur.
 * If the stream live state changes, or chat messages come in while the
 * page is backgrounded, this component will update the title to reflect it. *
 * @component
 */
import { FC, useEffect } from 'react';
import { useRecoilValue } from 'recoil';
import { serverStatusState, chatMessagesAtom } from '../stores/ClientConfigStore';

export const TitleNotifier: FC = () => {
  const chatMessages = useRecoilValue(chatMessagesAtom);
  const serverStatus = useRecoilValue(serverStatusState);

  let backgrounded = false;
  let defaultTitle = '';

  const onBlur = () => {
    console.log('onBlur');
    backgrounded = true;
    defaultTitle = document.title;
  };

  const onFocus = () => {
    console.log('onFocus');
    backgrounded = false;
    window.document.title = defaultTitle;
  };

  const listenForEvents = () => {
    // Listen for events that should update the title
    console.log('listenForEvents');
    window.addEventListener('blur', onBlur);
    window.addEventListener('focus', onFocus);
  };

  useEffect(() => {
    defaultTitle = window.document.title;
    listenForEvents();

    return () => {
      window.removeEventListener('focus', onFocus);
      window.removeEventListener('blur', onBlur);
    };
  }, []);

  useEffect(() => {
    const { online } = serverStatus;

    if (!backgrounded || !online) {
      return;
    }

    window.document.title = `💬 :: ${defaultTitle}`;
  }, [chatMessages]);

  useEffect(() => {
    if (!backgrounded) {
      return;
    }

    const { online } = serverStatus;

    if (online) {
      window.document.title = ` 🟢 :: ${defaultTitle}`;
    } else if (!online) {
      window.document.title = ` 🔴 :: ${defaultTitle}`;
    }
  }, [serverStatusState]);

  return null;
};
