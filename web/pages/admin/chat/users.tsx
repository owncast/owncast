import React, { useState, useEffect, useContext, ReactElement } from 'react';
import { Tabs } from 'antd';
import { ServerStatusContext } from '../../../utils/server-status-context';
import {
  CONNECTED_CLIENTS,
  fetchData,
  DISABLED_USERS,
  MODERATORS,
  BANNED_IPS,
} from '../../../utils/apis';
import { UserTable } from '../../../components/admin/UserTable';
import { ClientTable } from '../../../components/admin/ClientTable';
import { BannedIPsTable } from '../../../components/admin/BannedIPsTable';

import { AdminLayout } from '../../../components/layouts/AdminLayout';

export const FETCH_INTERVAL = 10 * 1000; // 10 sec

export default function ChatUsers() {
  const context = useContext(ServerStatusContext);
  const { online } = context || {};

  const [disabledUsers, setDisabledUsers] = useState([]);
  const [ipBans, setIPBans] = useState([]);
  const [clients, setClients] = useState([]);
  const [moderators, setModerators] = useState([]);

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

    try {
      const result = await fetchData(MODERATORS);
      setModerators(result);
    } catch (error) {
      console.error('error fetching moderators', error);
    }

    try {
      const result = await fetchData(BANNED_IPS);
      setIPBans(result);
    } catch (error) {
      console.error('error fetching banned ips', error);
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

  const connectedUserTabTitle = (
    <span>Connected {online ? `(${clients.length})` : '(offline)'}</span>
  );

  const bannedUsersTabTitle = <span>Banned Users ({disabledUsers.length})</span>;
  const bannedUsersTable = <UserTable data={disabledUsers} />;

  const bannedIPTabTitle = <span>IP Bans ({ipBans.length})</span>;
  const bannedIpTable = <BannedIPsTable data={ipBans} />;

  const moderatorUsersTabTitle = <span>Moderators ({moderators.length})</span>;
  const moderatorTable = <UserTable data={moderators} />;

  const items = [
    { label: connectedUserTabTitle, key: '1', children: connectedUsers },
    { label: bannedUsersTabTitle, key: '2', children: bannedUsersTable },
    { label: bannedIPTabTitle, key: '3', children: bannedIpTable },
    { label: moderatorUsersTabTitle, key: '4', children: moderatorTable },
  ];

  return <Tabs defaultActiveKey="1" items={items} />;
}

ChatUsers.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};
