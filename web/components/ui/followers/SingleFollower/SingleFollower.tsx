import { Avatar, Col, Row, Typography } from 'antd';
import React, { FC } from 'react';
import cn from 'classnames';
import { Follower } from '../../../../interfaces/follower';
import styles from './SingleFollower.module.scss';

export type SingleFollowerProps = {
  follower: Follower;
};

export const SingleFollower: FC<SingleFollowerProps> = ({ follower }) => (
  <div className={cn([styles.follower, 'followers-follower'])}>
    <a href={follower.link} target="_blank" rel="noreferrer">
      <Row wrap={false}>
        <Col span={6}>
          <Avatar src={follower.image} size="large" alt="Avatar" className={styles.avatar}>
            {follower.username.charAt(0).toUpperCase()}
          </Avatar>
        </Col>
        <Col span={16}>
          <Row className={styles.name}>
            <Typography.Text ellipsis>
              {follower.name || follower.username.split('@', 2)[0]}
            </Typography.Text>
          </Row>
          <Row className={styles.account}>
            <Typography.Text ellipsis>{follower.username}</Typography.Text>
          </Row>
        </Col>
      </Row>
    </a>
  </div>
);
