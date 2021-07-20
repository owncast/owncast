import React, { useState, useEffect, useContext } from 'react';
import { Typography } from 'antd';
import { ServerStatusContext } from '../../utils/server-status-context';
import { CONNECTED_CLIENTS, fetchData, DISABLED_USERS } from '../../utils/apis';
import UserTable from '../../components/user-table';
import ClientTable from '../../components/client-table';

const { Title } = Typography;

export const FETCH_INTERVAL = 10 * 1000; // 10 sec

export default function ChatUsers() {
  const context = useContext(ServerStatusContext);
  const { online } = context || {};

  const [disabledUsers, setDisabledUsers] = useState([]);
  const [clients, setClients] = useState([]);

  const getInfo = async () => {
    try {
      const result = await fetchData(DISABLED_USERS);
      setDisabledUsers(result);
    } catch (error) {
      console.log('==== error', error);
    }

    try {
      const result = await fetchData(CONNECTED_CLIENTS);
      setClients(result);
    } catch (error) {
      console.log('==== error', error);
    }
  };

  useEffect(() => {
    let getStatusIntervalId = null;

    getInfo();

    getStatusIntervalId = setInterval(getInfo, FETCH_INTERVAL);
    // returned function will be called on component unmount
    return () => {
      clearInterval(getStatusIntervalId);
    };
  }, [online]);

  const connectedUsers = online ? (
    <>
      <ClientTable data={clients} />
      <p className="description">
        Visit the{' '}
        <a
          href="https://owncast.online/docs/viewers/?source=admin"
          target="_blank"
          rel="noopener noreferrer"
        >
          documentation
        </a>{' '}
        to configure additional details about your viewers.
      </p>
    </>
  ) : (
    <p className="description">
      When a stream is active and chat is enabled, connected chat clients will be displayed here.
    </p>
  );
  return (
    <>
      <Title>Connected Chat Participants</Title>
      {connectedUsers}
      <br />
      <br />
      <Title>Banned Users</Title>
      <UserTable data={disabledUsers} />
    </>
  );
}
