import React, { useState, useEffect } from "react";
import { Table, Tag, Space, Button, Modal, Checkbox, Input, Typography } from 'antd';
import { DeleteOutlined, EyeTwoTone, EyeInvisibleOutlined } from '@ant-design/icons';
const { Title, Paragraph, Text } = Typography;

import format from 'date-fns/format'

import {
    fetchData,
    ACCESS_TOKENS,
    DELETE_ACCESS_TOKEN,
    CREATE_ACCESS_TOKEN,
} from "../utils/apis";

const scopeMapping = {
    'CAN_SEND_SYSTEM_MESSAGES': 'system chat',
    'CAN_SEND_MESSAGES': 'user chat',
};

function convertScopeStringToRenderString(scope) {
    if (!scope || !scopeMapping[scope]) {
        return "unknown";
    }

    return scopeMapping[scope].toUpperCase();
}

function NewTokenModal(props) {
    var selectedScopes = [];

    const scopes = [
        {
            value: 'CAN_SEND_SYSTEM_MESSAGES',
            label: 'Can send system chat messages',
            description: 'Can send chat messages as the offical system user.'
        },
        {
            value: 'CAN_SEND_MESSAGES',
            label: 'Can send user chat messages',
            description: 'Can send chat messages as any user name.'
        },
    ]

    function onChange(checkedValues) {
        selectedScopes = checkedValues
    }

    function saveToken() {
        props.onOk(name, selectedScopes)
    }

    const [name, setName] = useState('');

    return (
        <Modal title="Create New Access token" visible={props.visible} onOk={saveToken} onCancel={props.onCancel}>
            <p><Input value={name} placeholder="Access token name/description" onChange={(input) => setName(input.currentTarget.value)} /></p>

            <p>
                Select the permissions this access token will have.  It cannot be edited after it's created.
            </p>
            <Checkbox.Group options={scopes} onChange={onChange} />
        </Modal>
    )
}

export default function AccessTokens() {
    const [tokens, setTokens] = useState([]);
    const [isTokenModalVisible, setIsTokenModalVisible] = useState(false);

    const columns = [
        {
            title: '',
            key: 'delete',
            render: (text, record) => (
                <Space size="middle">
                    <Button onClick={() => handleDeleteToken(record.token)} icon={<DeleteOutlined />} />
                </Space>
            )
        },
        {
            title: 'Name',
            dataIndex: 'name',
            key: 'name',
        },
        {
            title: 'Token',
            dataIndex: 'token',
            key: 'token',
            render: (text, record) => (
               <Input.Password size="small" bordered={false} value={text}/>
            )
        },
        {
            title: 'Scopes',
            dataIndex: 'scopes',
            key: 'scopes',
            render: scopes => (
                <>
                    {scopes.map(scope => {
                        const color = 'purple';

                        return (
                            <Tag color={color} key={scope}>
                                {convertScopeStringToRenderString(scope)}
                            </Tag>
                        );
                    })}
                </>
            ),
        },
        {
            title: 'Last Used',
            dataIndex: 'lastUsed',
            key: 'lastUsed',
            render: (timestamp) => {
                const dateObject = new Date(timestamp);
                return format(dateObject, 'P p');
            },
        },
    ];

    const getAccessTokens = async () => {
        try {
            const result = await fetchData(ACCESS_TOKENS);
            setTokens(result);
        } catch (error) {
            handleError(error);
        }
    };

    useEffect(() => {
        getAccessTokens();
    }, []);

    async function handleDeleteToken(token) {
        try {
            const result = await fetchData(DELETE_ACCESS_TOKEN, { method: 'POST', data: { token: token } });
            getAccessTokens();
        } catch (error) {
            handleError(error);
        }
    }

    async function handleSaveToken(name: string, scopes: string[]) {
        try {
            const result = await fetchData(CREATE_ACCESS_TOKEN, { method: 'POST', data: { name: name, scopes: scopes } });
            getAccessTokens();
        } catch (error) {
            handleError(error);
        }
    }

    function handleError(error) {
        console.error("error", error);
        alert(error);
    }

    const showCreateTokenModal = () => {
        setIsTokenModalVisible(true);
    };

    const handleTokenModalSaveButton = (name, scopes) => {
        setIsTokenModalVisible(false);
        handleSaveToken(name, scopes);
    };

    const handleTokenModalCancel = () => {
        setIsTokenModalVisible(false);
    };

    return (
        <div>
            <Title>Access Tokens</Title>
            <Paragraph>
                Access tokens are used to allow external, 3rd party tools to perform specific actions on your Owncast server.
                They should be kept secure and never included in client code, instead they should be kept on a server that you control.
            </Paragraph>
            <Paragraph>
                Read more about how to use these tokens at _some documentation here_.
            </Paragraph>

            <Table rowKey="token" columns={columns} dataSource={tokens} pagination={false} />
            <Button onClick={showCreateTokenModal}>Create Access Token</Button>
            <NewTokenModal visible={isTokenModalVisible} onOk={handleTokenModalSaveButton} onCancel={handleTokenModalCancel} />
        </div>
    );
}
