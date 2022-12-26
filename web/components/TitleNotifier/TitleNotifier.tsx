/**
 * This component is responsible for updating the title of the page when
 * different state changes occur.
 * If the stream live state changes, or chat messages come in while the
 * page is backgrounded, this component will update the title to reflect it. *
 * @component
 */
import { FC, useEffect, useState } from 'react';
import { useRecoilValue } from 'recoil';
import { serverStatusState, chatMessagesAtom } from '../stores/ClientConfigStore';

export type TitleNotifierProps = {
  name: string;
};

export const TitleNotifier: FC<TitleNotifierProps> = ({ name }) => {
  const chatMessages = useRecoilValue(chatMessagesAtom);
  const serverStatus = useRecoilValue(serverStatusState);

  const [backgrounded, setBackgrounded] = useState(false);

  const { online } = serverStatus;

  const setTitle = (title: string) => {
    document.title = title;
  };

  const onBlur = () => {
    setBackgrounded(true);
  };

  const onFocus = () => {
    setBackgrounded(false);
    setTitle(name);
  };

  const listenForEvents = () => {
    // Listen for events that should update the title
    window.addEventListener('blur', onBlur);
    window.addEventListener('focus', onFocus);
  };

  const removeEvents = () => {
    window.removeEventListener('blur', onBlur);
    window.removeEventListener('focus', onFocus);
  };

  useEffect(() => {
    listenForEvents();

    return () => {
      removeEvents();
    };
  }, [name]);

  useEffect(() => {
    if (!backgrounded || !online) {
      return;
    }

    // Only alert on real chat messages from people.
    const lastMessage = chatMessages[chatMessages.length - 1];
    if (lastMessage.type !== 'CHAT') {
      return;
    }

    setTitle(`💬 :: ${name}`);
  }, [chatMessages, name]);

  useEffect(() => {
    if (!backgrounded) {
      return;
    }

    if (online) {
      setTitle(` 🟢 :: ${name}`);
    } else if (!online) {
      setTitle(` 🔴 :: ${name}`);
    }
  }, [online, name]);

  return null;
};
