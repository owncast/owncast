import React, { useState, useEffect, useContext } from 'react';
import { Tabs } from 'antd';
import { ServerStatusContext } from '../../utils/server-status-context';
import {
  CONNECTED_CLIENTS,
  fetchData,
  DISABLED_USERS,
  MODERATORS,
  BANNED_IPS,
} from '../../utils/apis';
import UserTable from '../../components/user-table';
import ClientTable from '../../components/client-table';
import BannedIPsTable from '../../components/banned-ips-table';

const { TabPane } = Tabs;

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

  return (
    <Tabs defaultActiveKey="1">
      <TabPane tab={<span>Connected {online ? `(${clients.length})` : '(offline)'}</span>} key="1">
        {connectedUsers}
      </TabPane>
      <TabPane tab={<span>Banned Users ({disabledUsers.length})</span>} key="2">
        <UserTable data={disabledUsers} />
      </TabPane>
      <TabPane tab={<span>IP Bans ({ipBans.length})</span>} key="3">
        <BannedIPsTable data={ipBans} />
      </TabPane>
      <TabPane tab={<span>Moderators ({moderators.length})</span>} key="4">
        <UserTable data={moderators} />
      </TabPane>
    </Tabs>
  );
}
