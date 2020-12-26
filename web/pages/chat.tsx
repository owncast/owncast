import React, { useState, useEffect } from "react";
import { Table, Typography, Tooltip, Switch, Button } from "antd";
import { CheckCircleFilled, ExclamationCircleFilled } from "@ant-design/icons";
import { RowSelectionType } from "antd/es/table/interface";
import { ColumnsType } from 'antd/es/table';
import format from 'date-fns/format'

import MessageVisiblityToggle from './components/message-visiblity-toggle';

import { CHAT_HISTORY, fetchData, UPDATE_CHAT_MESSGAE_VIZ } from "../utils/apis";
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
export const OUTCOME_TIMEOUT = 3000;

export default function Chat() {
  const [messages, setMessages] = useState([]);
  const [selectedRowKeys, setSelectedRows] = useState([]);
  const [bulkVisibility, setBulkVisibility] = useState(false);
  const [bulkProcessing, setBulkProcessing] = useState(false);
  const [bulkOutcome, setBulkOutcome] = useState(null);
  let outcomeTimeout = null;

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
    messages.splice(messageIndex, 1, message)
    setMessages([...messages]);
  };

  useEffect(() => {
    getInfo();
    return () => {
      clearTimeout(outcomeTimeout);
    };
  }, []);

  const nameFilters = createUserNameFilters(messages);
  
  const rowSelection: RowSelectionType = {
    selectedRowKeys,
    onChange: selectedKeys => {
      setSelectedRows(selectedKeys);
    },
  };

  const handleBulkToggle = checked => {
    setBulkVisibility(checked);
  };
  const resetBulkOutcome = () => {
    outcomeTimeout = setTimeout(() => { setBulkOutcome(null)}, OUTCOME_TIMEOUT);
  };
  const handleSubmitBulk = async () => {
    setBulkProcessing(true);
    const result = await fetchData(UPDATE_CHAT_MESSGAE_VIZ, {
      auth: true,
      method: 'POST',
      data: {
        visible: bulkVisibility,
        idArray: selectedRowKeys,
      },
    });

    if (result.success && result.message === "changed") {
      setBulkOutcome(<CheckCircleFilled />);
      resetBulkOutcome();

      // update messages
      const updatedList = [...messages];
      selectedRowKeys.map(key => {
        const messageIndex = updatedList.findIndex(m => m.id === key);
        const newMessage = {...messages[messageIndex], visible: bulkVisibility };
        updatedList.splice(messageIndex, 1, newMessage);
        return null;
      });
      setMessages(updatedList);
    } else {
      setBulkOutcome(<ExclamationCircleFilled />);
      resetBulkOutcome();
    }
    setBulkProcessing(false);
    
  }


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
      title: 'Show / Hide',
      dataIndex: 'visible',
      key: 'visible',
      filters: [{ text: 'visible', value: true }, { text: 'hidden', value: false }],
      onFilter: (value, record) => record.visible === value,
      render: (visible, record) => (
        <MessageVisiblityToggle
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
      <div className="bulk-editor">
        <span className="label">Check multiple messages to change their visibility to: </span>
        
        <Switch
          className="toggler"
          disabled={!selectedRowKeys.length}
          checkedChildren="show"
          unCheckedChildren="hide"
          onChange={handleBulkToggle}
          checked={bulkVisibility}
        />

        <Button
          type="primary"
          size="small"
          loading={bulkProcessing}
          disabled={!selectedRowKeys.length}
          icon={bulkOutcome}
          onClick={handleSubmitBulk}
        >
          Go
        </Button>
      </div>
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


