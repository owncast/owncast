import React, { useState, useEffect, ReactElement } from 'react';
import { Table, Typography, Button } from 'antd';
import classNames from 'classnames';
import { ColumnsType } from 'antd/es/table';
import { format } from 'date-fns';

import dynamic from 'next/dynamic';
import { useTranslation } from 'next-export-i18n';
import { MessageType } from '../../../types/chat';
import {
  CHAT_HISTORY,
  fetchData,
  FETCH_INTERVAL,
  UPDATE_CHAT_MESSGAE_VIZ,
} from '../../../utils/apis';
import { isEmptyObject } from '../../../utils/format';
import { MessageVisiblityToggle } from '../../../components/admin/MessageVisiblityToggle';
import { UserPopover } from '../../../components/admin/UserPopover';

import { AdminLayout } from '../../../components/layouts/AdminLayout';

const { Title } = Typography;

// Lazy loaded components

const CheckCircleFilled = dynamic(() => import('@ant-design/icons/CheckCircleFilled'), {
  ssr: false,
});

const ExclamationCircleFilled = dynamic(() => import('@ant-design/icons/ExclamationCircleFilled'), {
  ssr: false,
});

function createUserNameFilters(messages: MessageType[]) {
  const filtered = messages.reduce((acc, curItem) => {
    const curAuthor = curItem.user.id;
    if (!acc.some(item => item.text === curAuthor)) {
      acc.push({ text: curAuthor, value: curAuthor });
    }
    return acc;
  }, []);

  // sort by name
  return filtered.sort((a, b) => {
    const nameA = a.text.toUpperCase(); // ignore upper and lowercase
    const nameB = b.text.toUpperCase(); // ignore upper and lowercase
    if (nameA < nameB) {
      return -1;
    }
    if (nameA > nameB) {
      return 1;
    }
    // names must be equal
    return 0;
  });
}
export const OUTCOME_TIMEOUT = 3000;

export default function Chat() {
  const [messages, setMessages] = useState([]);
  const [selectedRowKeys, setSelectedRows] = useState([]);
  const [bulkProcessing, setBulkProcessing] = useState(false);
  const [bulkOutcome, setBulkOutcome] = useState(null);
  const [bulkAction, setBulkAction] = useState('');
  const { t } = useTranslation();
  let outcomeTimeout = null;
  let chatReloadInterval = null;

  const getInfo = async () => {
    try {
      const result = await fetchData(CHAT_HISTORY, { auth: true });
      if (isEmptyObject(result)) {
        setMessages([]);
      } else {
        setMessages(result);
      }
    } catch (error) {
      console.log('==== error', error);
    }
  };

  useEffect(() => {
    getInfo();

    chatReloadInterval = setInterval(() => {
      getInfo();
    }, FETCH_INTERVAL);

    return () => {
      clearTimeout(outcomeTimeout);
      clearTimeout(chatReloadInterval);
    };
  }, []);

  const nameFilters = createUserNameFilters(messages);

  const rowSelection = {
    selectedRowKeys,
    onChange: (selectedKeys: string[]) => {
      setSelectedRows(selectedKeys);
    },
  };

  const updateMessage = message => {
    const messageIndex = messages.findIndex(m => m.id === message.id);
    messages.splice(messageIndex, 1, message);
    setMessages([...messages]);
  };

  const resetBulkOutcome = () => {
    outcomeTimeout = setTimeout(() => {
      setBulkOutcome(null);
      setBulkAction('');
    }, OUTCOME_TIMEOUT);
  };
  const handleSubmitBulk = async bulkVisibility => {
    setBulkProcessing(true);
    const result = await fetchData(UPDATE_CHAT_MESSGAE_VIZ, {
      auth: true,
      method: 'POST',
      data: {
        visible: bulkVisibility,
        idArray: selectedRowKeys,
      },
    });

    if (result.success && result.message === 'changed') {
      setBulkOutcome(<CheckCircleFilled />);
      resetBulkOutcome();

      // update messages
      const updatedList = [...messages];
      selectedRowKeys.map(key => {
        const messageIndex = updatedList.findIndex(m => m.id === key);
        const newMessage = { ...messages[messageIndex], visible: bulkVisibility };
        updatedList.splice(messageIndex, 1, newMessage);
        return null;
      });
      setMessages(updatedList);
      setSelectedRows([]);
    } else {
      setBulkOutcome(<ExclamationCircleFilled />);
      resetBulkOutcome();
    }
    setBulkProcessing(false);
  };
  const handleSubmitBulkShow = () => {
    setBulkAction('show');
    handleSubmitBulk(true);
  };
  const handleSubmitBulkHide = () => {
    setBulkAction('hide');
    handleSubmitBulk(false);
  };

  const chatColumns: ColumnsType<MessageType> = [
    {
      title: t('Time'),
      dataIndex: 'timestamp',
      key: 'timestamp',
      className: 'timestamp-col',
      defaultSortOrder: 'descend',
      render: timestamp => {
        const dateObject = new Date(timestamp);
        return format(dateObject, 'PP pp');
      },
      sorter: (a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime(),
      width: 90,
    },
    {
      title: t('User'),
      dataIndex: 'user',
      key: 'user',
      className: 'name-col',
      filters: nameFilters,
      onFilter: (value, record) => record.user.id === value,
      sorter: (a, b) => a.user.displayName.localeCompare(b.user.displayName),
      sortDirections: ['ascend', 'descend'],
      ellipsis: true,
      render: user => {
        const { displayName } = user;
        return <UserPopover user={user}>{displayName}</UserPopover>;
      },
      width: 110,
    },
    {
      title: t('Message'),
      dataIndex: 'body',
      key: 'body',
      className: 'message-col',
      width: 320,
      render: body => (
        <div
          className="message-contents"
          // eslint-disable-next-line react/no-danger
          dangerouslySetInnerHTML={{ __html: body }}
        />
      ),
    },
    {
      title: '',
      dataIndex: 'hiddenAt',
      key: 'hiddenAt',
      className: 'toggle-col',
      filters: [
        { text: t('Visible messages'), value: true },
        { text: t('Hidden messages'), value: false },
      ],
      onFilter: (value, record) => record.visible === value,
      render: (hiddenAt, record) => (
        <MessageVisiblityToggle isVisible={!hiddenAt} message={record} setMessage={updateMessage} />
      ),
      width: 30,
    },
  ];

  const bulkDivClasses = classNames({
    'bulk-editor': true,
    active: selectedRowKeys.length,
  });

  return (
    <div className="chat-messages">
      <Title>{t('Chat Messages')}</Title>
      <p>{t('Manage the messages from viewers that show up on your stream.')}</p>
      <div className={bulkDivClasses}>
        <span className="label">
          {t('Check multiple messages to change their visibility to:')}{' '}
        </span>

        <Button
          type="primary"
          size="small"
          shape="round"
          className="button"
          loading={bulkAction === 'show' && bulkProcessing}
          icon={bulkAction === 'show' && bulkOutcome}
          disabled={!selectedRowKeys.length || (bulkAction && bulkAction !== 'show')}
          onClick={handleSubmitBulkShow}
        >
          {t('Show')}
        </Button>
        <Button
          type="primary"
          size="small"
          shape="round"
          className="button"
          loading={bulkAction === 'hide' && bulkProcessing}
          icon={bulkAction === 'hide' && bulkOutcome}
          disabled={!selectedRowKeys.length || (bulkAction && bulkAction !== 'hide')}
          onClick={handleSubmitBulkHide}
        >
          {t('Hide')}
        </Button>
      </div>
      <Table
        size="small"
        className="table-container"
        pagination={{ defaultPageSize: 100, showSizeChanger: true }}
        scroll={{ y: 540 }}
        rowClassName={record => (record.hiddenAt ? 'hidden' : '')}
        dataSource={messages}
        columns={chatColumns}
        rowKey={row => row.id}
        rowSelection={rowSelection}
      />
    </div>
  );
}

Chat.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};
