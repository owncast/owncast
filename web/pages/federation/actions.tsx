import React, { useEffect, useState } from 'react';
import { Table, Typography } from 'antd';
import { ColumnsType } from 'antd/lib/table/interface';
import format from 'date-fns/format';
import { FEDERATION_ACTIONS, fetchData } from '../../utils/apis';

import { isEmptyObject } from '../../utils/format';

const { Title, Paragraph } = Typography;

export interface Action {
  iri: string;
  actorIRI: string;
  type: string;
  timestamp: Date;
}

export default function FediverseActions() {
  const [actions, setActions] = useState<Action[]>([]);

  const getActions = async () => {
    try {
      const result = await fetchData(FEDERATION_ACTIONS, { auth: true });
      if (isEmptyObject(result)) {
        setActions([]);
      } else {
        setActions(result);
      }
    } catch (error) {
      console.log('==== error', error);
    }
  };

  useEffect(() => {
    getActions();
  }, []);

  const columns: ColumnsType<Action> = [
    {
      title: 'Action',
      dataIndex: 'type',
      key: 'type',
      width: 50,
      render: (_, record) => {
        let image;
        let title;
        switch (record.type) {
          case 'FEDIVERSE_ENGAGEMENT_REPOST':
            image = '/img/repost.svg';
            title = 'Share';
            break;
          case 'FEDIVERSE_ENGAGEMENT_LIKE':
            image = '/img/like.svg';
            title = 'Like';
            break;
          case 'FEDIVERSE_ENGAGEMENT_FOLLOW':
            image = '/img/follow.svg';
            title = 'Follow';
            break;
          default:
            image = '';
        }

        return (
          <div
            style={{
              width: '100%',
              height: '100%',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              flexDirection: 'column',
            }}
          >
            <img src={image} width="70%" alt={title} title={title} />
            <div style={{ fontSize: '0.7rem' }}>{title}</div>
          </div>
        );
      },
    },
    {
      title: 'From',
      dataIndex: 'actorIRI',
      key: 'from',
      render: (_, record) => <a href={record.actorIRI}>{record.actorIRI}</a>,
    },
    {
      title: 'When',
      dataIndex: 'timestamp',
      key: 'timestamp',
      render: (_, record) => {
        const dateObject = new Date(record.timestamp);
        return format(dateObject, 'P pp');
      },
    },
  ];

  function makeTable(data: Action[], tableColumns: ColumnsType<Action>) {
    return (
      <Table
        dataSource={data}
        columns={tableColumns}
        size="small"
        rowKey={row => row.iri}
        pagination={{ pageSize: 50 }}
      />
    );
  }

  return (
    <div>
      <Title level={3}>Fediverse Actions</Title>
      <Paragraph>
        Below is a list of actions that were taken by others in response to your posts as well as
        people who requested to follow you.
      </Paragraph>
      {makeTable(actions, columns)}
    </div>
  );
}
