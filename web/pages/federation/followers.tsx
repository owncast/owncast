import React, { useEffect, useState, useContext } from 'react';
import { Table, Avatar, Button, Tabs } from 'antd';
import { ColumnsType, SortOrder } from 'antd/lib/table/interface';
import format from 'date-fns/format';
import { UserAddOutlined, UserDeleteOutlined } from '@ant-design/icons';
import { ServerStatusContext } from '../../utils/server-status-context';
import {
  FOLLOWERS,
  FOLLOWERS_PENDING,
  SET_FOLLOWER_APPROVAL,
  FOLLOWERS_BLOCKED,
  fetchData,
} from '../../utils/apis';
import { isEmptyObject } from '../../utils/format';

const { TabPane } = Tabs;
export interface Follower {
  link: string;
  username: string;
  image: string;
  name: string;
  timestamp: Date;
  approved: Date;
}

export default function FediverseFollowers() {
  const [followersPending, setFollowersPending] = useState<Follower[]>([]);
  const [followersBlocked, setFollowersBlocked] = useState<Follower[]>([]);
  const [followers, setFollowers] = useState<Follower[]>([]);
  const [totalCount, setTotalCount] = useState<number>(0);
  const [currentPage, setCurrentPage] = useState<number>(0);

  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig } = serverStatusData || {};
  const { federation } = serverConfig;
  const { isPrivate } = federation;

  const getFollowers = async () => {
    try {
      const limit = 50;
      const offset = currentPage * limit;
      const u = `${FOLLOWERS}?offset=${offset}&limit=${limit}`;

      // Active followers
      const result = await fetchData(u, { auth: true });
      const { results, total } = result;

      if (isEmptyObject(results)) {
        setFollowers([]);
      } else {
        setTotalCount(total);
        setFollowers(results);
      }

      // Pending follow requests
      const pendingFollowersResult = await fetchData(FOLLOWERS_PENDING, { auth: true });
      if (isEmptyObject(pendingFollowersResult)) {
        setFollowersPending([]);
      } else {
        setFollowersPending(pendingFollowersResult);
      }

      // Blocked/rejected followers
      const blockedFollowersResult = await fetchData(FOLLOWERS_BLOCKED, { auth: true });
      if (isEmptyObject(followersBlocked)) {
        setFollowersBlocked([]);
      } else {
        setFollowersBlocked(blockedFollowersResult);
      }
    } catch (error) {
      console.log('==== error', error);
    }
  };

  useEffect(() => {
    getFollowers();
  }, []);

  const columns: ColumnsType<Follower> = [
    {
      title: '',
      dataIndex: 'image',
      key: 'image',
      width: 90,
      render: image => <Avatar size={40} src={image || '/img/logo.svg'} />,
    },
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
      render: (_, follower) => (
        <a href={follower.link} target="_blank" rel="noreferrer">
          {follower.name || follower.username}
        </a>
      ),
    },
    {
      title: 'URL',
      dataIndex: 'link',
      key: 'link',
      render: (_, follower) => (
        <a href={follower.link} target="_blank" rel="noreferrer">
          {follower.link}
        </a>
      ),
    },
  ];

  function makeTable(data: Follower[], tableColumns: ColumnsType<Follower>) {
    return (
      <Table
        dataSource={data}
        columns={tableColumns}
        size="small"
        rowKey={row => row.link}
        pagination={{
          pageSize: 50,
          hideOnSinglePage: true,
          showSizeChanger: false,
          total: totalCount,
        }}
        onChange={pagination => {
          const page = pagination.current;
          setCurrentPage(page);
        }}
      />
    );
  }

  async function approveFollowRequest(request) {
    try {
      await fetchData(SET_FOLLOWER_APPROVAL, {
        auth: true,
        method: 'POST',
        data: {
          actorIRI: request.link,
          approved: true,
        },
      });

      // Refetch and update the current data.
      getFollowers();
    } catch (err) {
      console.error(err);
    }
  }

  async function rejectFollowRequest(request) {
    try {
      await fetchData(SET_FOLLOWER_APPROVAL, {
        auth: true,
        method: 'POST',
        data: {
          actorIRI: request.link,
          approved: false,
        },
      });

      // Refetch and update the current data.
      getFollowers();
    } catch (err) {
      console.error(err);
    }
  }

  const pendingColumns: ColumnsType<Follower> = [...columns];
  pendingColumns.unshift(
    {
      title: 'Approve',
      dataIndex: null,
      key: null,
      render: request => (
        <Button
          type="primary"
          icon={<UserAddOutlined />}
          onClick={() => {
            approveFollowRequest(request);
          }}
        />
      ),
      width: 50,
    },
    {
      title: 'Reject',
      dataIndex: null,
      key: null,
      render: request => (
        <Button
          type="primary"
          danger
          icon={<UserDeleteOutlined />}
          onClick={() => {
            rejectFollowRequest(request);
          }}
        />
      ),
      width: 50,
    },
  );

  pendingColumns.push({
    title: 'Requested',
    dataIndex: 'timestamp',
    key: 'requested',
    width: 200,
    render: timestamp => {
      const dateObject = new Date(timestamp);
      return <>{format(dateObject, 'P')}</>;
    },
    sorter: (a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime(),
    sortDirections: ['descend', 'ascend'] as SortOrder[],
    defaultSortOrder: 'descend' as SortOrder,
  });

  const blockedColumns: ColumnsType<Follower> = [...columns];
  blockedColumns.unshift({
    title: 'Approve',
    dataIndex: null,
    key: null,
    render: request => (
      <Button
        type="primary"
        icon={<UserAddOutlined />}
        size="large"
        onClick={() => {
          approveFollowRequest(request);
        }}
      />
    ),
    width: 50,
  });

  blockedColumns.push(
    {
      title: 'Requested',
      dataIndex: 'timestamp',
      key: 'requested',
      width: 200,
      render: timestamp => {
        const dateObject = new Date(timestamp);
        return <>{format(dateObject, 'P')}</>;
      },
      sorter: (a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime(),
      sortDirections: ['descend', 'ascend'] as SortOrder[],
      defaultSortOrder: 'descend' as SortOrder,
    },
    {
      title: 'Rejected/Blocked',
      dataIndex: 'timestamp',
      key: 'disabled_at',
      width: 200,
      render: timestamp => {
        const dateObject = new Date(timestamp);
        return <>{format(dateObject, 'P')}</>;
      },
      sorter: (a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime(),
      sortDirections: ['descend', 'ascend'] as SortOrder[],
      defaultSortOrder: 'descend' as SortOrder,
    },
  );

  const followersColumns: ColumnsType<Follower> = [...columns];

  followersColumns.push(
    {
      title: 'Added',
      dataIndex: 'timestamp',
      key: 'timestamp',
      width: 200,
      render: timestamp => {
        const dateObject = new Date(timestamp);
        return <>{format(dateObject, 'P')}</>;
      },
      sorter: (a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime(),
      sortDirections: ['descend', 'ascend'] as SortOrder[],
      defaultSortOrder: 'descend' as SortOrder,
    },
    {
      title: 'Remove',
      dataIndex: null,
      key: null,
      render: request => (
        <Button
          type="primary"
          danger
          icon={<UserDeleteOutlined />}
          onClick={() => {
            rejectFollowRequest(request);
          }}
        />
      ),
      width: 50,
    },
  );

  const pendingRequestsTab = isPrivate && (
    <TabPane
      tab={<span>Requests {followersPending.length > 0 && `(${followersPending.length})`}</span>}
      key="2"
    >
      <p>
        The following people are requesting to follow your Owncast server on the{' '}
        <a href="https://en.wikipedia.org/wiki/Fediverse" target="_blank" rel="noopener noreferrer">
          Fediverse
        </a>{' '}
        and be alerted to when you go live. Each must be approved.
      </p>
      {makeTable(followersPending, pendingColumns)}
    </TabPane>
  );

  return (
    <div className="followers-section">
      <Tabs defaultActiveKey="1">
        <TabPane
          tab={<span>Followers {followers.length > 0 && `(${followers.length})`}</span>}
          key="1"
        >
          <p>The following accounts get notified when you go live or send a post.</p>
          {makeTable(followers, followersColumns)}{' '}
        </TabPane>
        {pendingRequestsTab}
        <TabPane
          tab={<span>Blocked {followersBlocked.length > 0 && `(${followersBlocked.length})`}</span>}
          key="3"
        >
          <p>
            The following people were either rejected or blocked by you. You can approve them as a
            follower.
          </p>
          <p>{makeTable(followersBlocked, blockedColumns)}</p>
        </TabPane>
      </Tabs>
    </div>
  );
}
