/* eslint-disable camelcase */
/* eslint-disable react/no-danger */
import React, { useState, useEffect } from 'react';
import { Typography, Skeleton } from 'antd';
import format from 'date-fns/format';

import { fetchExternalData } from '../utils/apis';

const { Title, Link } = Typography;

const OWNCAST_FEED_URL = 'https://owncast.online/news/index.json';
const OWNCAST_BASE_URL = 'https://owncast.online';

interface Article {
  title: string;
  url: string;
  content_html: string;
  date_published: string;
}

function ArticleItem({ title, url, content_html: content, date_published: date }: Article) {
  const dateObject = new Date(date);
  const dateString = format(dateObject, 'MMM dd, yyyy, HH:mm');
  return (
    <article>
      <p className="timestamp">{dateString}</p>
      <Title level={3}>
        <Link href={`${OWNCAST_BASE_URL}${url}`} target="_blank" rel="noopener noreferrer">
          {title}
        </Link>
      </Title>
      <div dangerouslySetInnerHTML={{ __html: content }} />
    </article>
  );
}

export default function NewsFeed() {
  const [feed, setFeed] = useState<Article[]>([]);
  const [loading, setLoading] = useState<Boolean>(true);

  const getFeed = async () => {
    setLoading(false);

    try {
      const result = await fetchExternalData(OWNCAST_FEED_URL);
      setFeed(result.items);
    } catch (error) {
      console.log('==== error', error);
    }
  };

  useEffect(() => {
    getFeed();
  }, []);

  const loadingSpinner = loading ? (<Skeleton loading={true} active />) : null;
  const noNews = !loading && feed.length === 0 ? (<div>No news.</div>) : null;

  return (
    <section className="news-feed form-module">
      <Title level={2}>News &amp; Updates from Owncast</Title>
      {loadingSpinner}
      {feed.map(item => (
        <ArticleItem {...item} key={item.url} />
      ))}

      {noNews}
    </section>
  );
}
