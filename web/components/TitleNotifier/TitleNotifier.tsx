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

  const setTitle = (title: string) => {
    console.trace('Debug: Setting title to', title);
    window.document.title = title;
  };

  const onBlur = () => {
    console.log('onBlur: Saving defaultTitle as', window.document.title);
    backgrounded = true;
    defaultTitle = window.document.title;
  };

  const onFocus = () => {
    backgrounded = false;
    setTitle(defaultTitle);
  };

  const listenForEvents = () => {
    // Listen for events that should update the title
    window.addEventListener('blur', onBlur);
    window.addEventListener('focus', onFocus);
  };

  useEffect(() => {
    console.log('useEffect: Saving defaultTitle as', window.document.title);

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
    setTitle(`💬 :: ${defaultTitle}`);
  }, [chatMessages]);

  useEffect(() => {
    if (!backgrounded) {
      return;
    }

    const { online } = serverStatus;

    if (online) {
      setTitle(` 🟢 :: ${defaultTitle}`);
    } else if (!online) {
      setTitle(` 🔴 :: ${defaultTitle}`);
    }
  }, [serverStatusState]);

  return null;
};
