import React, { useState, useEffect } from "react";
import { Table, Typography } from "antd";
import { ColumnsType } from 'antd/es/table';
import format from 'date-fns/format'

import {
  CHAT_HISTORY,
  fetchDataFromMain,
} from "../utils/apis";

const { Title } = Typography;

interface Message {
  key: string;
  author: string;
  body: string;
  id: string;
  name: string;
  timestamp: string;
  type: string;
  visible: boolean;
}


function createUserNameFilters(messages) {
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

  const getInfo = async () => {
    try {
      const result = await fetchDataFromMain(CHAT_HISTORY);
      setMessages(result);
    } catch (error) {
      console.log("==== error", error);
    }
  };

  useEffect(() => {
    getInfo();
  }, []);

  
  const nameFilters = createUserNameFilters(messages);
  const rowSelection = {
    onChange: (selectedRowKeys, selectedRows) => {
      console.log(`selectedRowKeys: ${selectedRowKeys}`, 'selectedRows: ', selectedRows);
    },
    getCheckboxProps: record => ({
      disabled: record.name === 'Disabled User', // Column configuration not to be checked
      name: record.name,
    }),
  };

  const chatColumns:  ColumnsType<Message> = [
    {
      title: 'Time',
      dataIndex: 'timestamp',
      key: 'timestamp',
      defaultSortOrder: 'descend',
      render: (timestamp) => {
        const dateObject = new Date(timestamp);
        return format(dateObject, 'P pp');
      },
      sorter: (a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime(),
    },
    {
      title: 'User',
      dataIndex: 'author',
      key: 'author',
      className: 'name-col',
      filters: nameFilters,
      onFilter: (value, record) => record.name.indexOf(value) === 0,
      sorter: (a, b) => a.name.toUppercase() - b.name.toUppercase(),
      sortDirections: ['ascend', 'descend'],  
    },
    {
      title: 'Message',
      dataIndex: 'body',
      key: 'body',
      className: 'message-col',
    },
    {
      title: 'Show/Hide',
      dataIndex: 'visible',
      key: 'visible',
    },
  ];

  
  return (
    <div className="chat-message">
      <Title level={2}>Chat Messages</Title>
      <Table
        size="small"
        pagination={{ pageSize: 100 }} 
        scroll={{ y: 540 }}
        rowClassName={record => !record.visible ? 'hidden' : ''}
        dataSource={messages}
        columns={chatColumns}
        rowKey={(row) => row.id}

        rowSelection={{
          type: "checkbox",
          ...rowSelection,
        }}

      />
  </div>)
}


