import { Avatar, Col, Row } from 'antd';
import React, { FC } from 'react';
import cn from 'classnames';
import { Follower } from '../../../../interfaces/follower';
import styles from './SingleFollower.module.scss';

export type SingleFollowerProps = {
  follower: Follower;
};

export const SingleFollower: FC<SingleFollowerProps> = ({ follower }) => (
  <div>
    <a href={follower.link} target="_blank" rel="noreferrer">
      <Row wrap={false} className={cn([styles.follower, 'followers-follower'])}>
        <Col span={6} className={styles.avatarColumn}>
          <Avatar
            src={follower.image}
            alt={(follower.name || follower.username).charAt(0).toUpperCase()}
            className={styles.avatar}
            size="large"
          >
            {(follower.name || follower.username).charAt(0).toUpperCase()}
          </Avatar>
        </Col>
        <Col>
          <Row className={styles.username}>
            {follower.name || follower.username.split('@', 2)[0]}
          </Row>
          <Row className={styles.account}>{follower.username}</Row>
        </Col>
      </Row>
    </a>
  </div>
);
