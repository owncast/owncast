import React, { useState, useEffect } from "react";
import { Table, Typography, Tooltip, Switch, Button, Result } from "antd";
import { RowSelectionType } from "antd/es/table/interface";
import { ColumnsType } from 'antd/es/table';
import format from 'date-fns/format'

import ToggleSwitch from './components/toggle';

import { CHAT_HISTORY, fetchData } from "../utils/apis";
import { MessageType } from '../types/chat';


const { Title } = Typography;

function createUserNameFilters(messages: MessageType[]) {
  const filtered = messages.reduce((acc, curItem) => {
    const curAuthor = curItem.author;
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

export default function Chat() {
  const [messages, setMessages] = useState([]);
  const [selectedRowKeys, setSelectedRows] = useState([]);

  const getInfo = async () => {
    try {
      const result = await fetchData(CHAT_HISTORY, { auth: false });
      setMessages(result);
    } catch (error) {
      console.log("==== error", error);
    }
  };

  const updateMessage = message => {
    const messageIndex = messages.findIndex(m => m.id === message.id);
    console.log("====update?", message, messages[messageIndex])
    messages.splice(messageIndex, 1, message)
    setMessages([...messages]);
  };

  useEffect(() => {
    getInfo();
  }, []);

  const nameFilters = createUserNameFilters(messages);
  
  const rowSelection: RowSelectionType = {
    selectedRowKeys,
    onChange: selectedKeys => {
      setSelectedRows(selectedKeys);
    },
  };


  const chatColumns: ColumnsType<MessageType> = [
    {
      title: 'Time',
      dataIndex: 'timestamp',
      key: 'timestamp',
      className: 'timestamp-col',
      defaultSortOrder: 'descend',
      render: (timestamp) => {
        const dateObject = new Date(timestamp);
        return format(dateObject, 'PP pp');
      },
      sorter: (a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime(),
      width: 80,
    },
    {
      title: 'User',
      dataIndex: 'author',
      key: 'author',
      className: 'name-col',
      filters: nameFilters,
      onFilter: (value, record) => record.author === value,
      sorter: (a, b) => a.author.localeCompare(b.author),
      sortDirections: ['ascend', 'descend'],  
      ellipsis: true,
      render: author => (
        <Tooltip placement="topLeft" title={author}>
          {author}
        </Tooltip>
      ),
      width: 80,
    },
    {
      title: 'Message',
      dataIndex: 'body',
      key: 'body',
      className: 'message-col',
      width: 230,
    },
    {
      title: 'Show/ Hide',
      dataIndex: 'visible',
      key: 'visible',
      filters: [{ text: 'visible', value: true }, { text: 'hidden', value: false }],
      onFilter: (value, record) => record.visible === value,
      render: (visible, record) => (
        <ToggleSwitch
          isVisible={visible}
          message={record}
          setMessage={updateMessage}
        />
      ),
      width: 60,
    },
  ];

  
  return (
    <div className="chat-messages">
      <Title level={2}>Chat Messages</Title>
      <p>click things and stuff</p>
      <Button
        type="primary"
        // onClick={}
        disabled={!selectedRowKeys.length}
        loading={false}
      >
        Bulk toggle
      </Button>
      <Switch />
      <Table
        size="small"
        className="messages-table"
        pagination={{ pageSize: 100 }} 
        scroll={{ y: 540 }}
        rowClassName={record => !record.visible ? 'hidden' : ''}
        dataSource={messages}
        columns={chatColumns}
        rowKey={(row) => row.id}
        rowSelection={rowSelection}
      />
  </div>)
}


