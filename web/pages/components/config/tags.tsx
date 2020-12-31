/* eslint-disable react/no-array-index-key */
import React, { useContext, useEffect } from 'react';
import { Typography, Button, Tooltip } from 'antd';
import { CloseCircleOutlined } from '@ant-design/icons';
import { ServerStatusContext } from '../../../utils/server-status-context';

const { Title } = Typography;

function Tag({ label }) {
  return (
      <Button className="tag" type="text" shape="round">
        {label}
        <Tooltip title="Delete this tag.">
          <Button type="link" size="small" className="tag-delete">
            <CloseCircleOutlined />
          </Button>
        </Tooltip>

      </Button>
  );
}

export default function EditInstanceTags() {
  const serverStatusData = useContext(ServerStatusContext);
  const { serverConfig } = serverStatusData || {};

  const { instanceDetails } = serverConfig;
  const { tags = [] } = instanceDetails;
  console.log(tags)
  
  return (
    <div className="tag-editor-container">

      <Title level={3}>Add Tags</Title>
      <p>This is a great way to categorize your Owncast server on the Directory!</p>

      <div className="tag-current-tags">
        {tags.map((tag, index) => <Tag label={tag} key={`tag-${tag}-${index}`} />)}
      </div>
    </div>
  );
}

