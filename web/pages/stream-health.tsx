/* eslint-disable react/no-unescaped-entities */
// import { BulbOutlined, LaptopOutlined, SaveOutlined } from '@ant-design/icons';
import { Row, Col, Typography, Space, Statistic, Card, Alert } from 'antd';
import React, { ReactNode, useEffect, useState } from 'react';
import { ClockCircleOutlined, WarningOutlined, WifiOutlined } from '@ant-design/icons';
import { fetchData, FETCH_INTERVAL, API_STREAM_HEALTH_METRICS } from '../utils/apis';
import Chart from '../components/chart';

interface TimedValue {
  time: Date;
  value: Number;
}

interface DescriptionBoxProps {
  title: String;
  description: ReactNode;
}

function DescriptionBox({ title, description }: DescriptionBoxProps) {
  return (
    <div className="description-box">
      <Typography.Title>{title}</Typography.Title>
      <Typography.Paragraph>{description}</Typography.Paragraph>
    </div>
  );
}

export default function StreamHealth() {
  const [errors, setErrors] = useState<TimedValue[]>([]);
  const [qualityVariantChanges, setQualityVariantChanges] = useState<TimedValue[]>([]);

  const [lowestLatency, setLowestLatency] = useState<TimedValue[]>();
  const [highestLatency, setHighestLatency] = useState<TimedValue[]>();
  const [medianLatency, setMedianLatency] = useState<TimedValue[]>([]);

  const [medianSegmentDownloadDurations, setMedianSegmentDownloadDurations] = useState<
    TimedValue[]
  >([]);
  const [maximumSegmentDownloadDurations, setMaximumSegmentDownloadDurations] = useState<
    TimedValue[]
  >([]);
  const [minimumSegmentDownloadDurations, setMinimumSegmentDownloadDurations] = useState<
    TimedValue[]
  >([]);
  const [minimumPlayerBitrate, setMinimumPlayerBitrate] = useState<TimedValue[]>([]);
  const [medianPlayerBitrate, setMedianPlayerBitrate] = useState<TimedValue[]>([]);
  const [maximumPlayerBitrate, setMaximumPlayerBitrate] = useState<TimedValue[]>([]);
  const [availableBitrates, setAvailableBitrates] = useState<Number[]>([]);
  const [segmentLength, setSegmentLength] = useState(0);

  const getMetrics = async () => {
    try {
      const result = await fetchData(API_STREAM_HEALTH_METRICS);
      setErrors(result.errors);
      setQualityVariantChanges(result.qualityVariantChanges);

      setHighestLatency(result.highestLatency);
      setLowestLatency(result.lowestLatency);
      setMedianLatency(result.medianLatency);

      setMedianSegmentDownloadDurations(result.medianSegmentDownloadDuration);
      setMaximumSegmentDownloadDurations(result.maximumSegmentDownloadDuration);
      setMinimumSegmentDownloadDurations(result.minimumSegmentDownloadDuration);

      setMinimumPlayerBitrate(result.minPlayerBitrate);
      setMedianPlayerBitrate(result.medianPlayerBitrate);
      setMaximumPlayerBitrate(result.maxPlayerBitrate);

      setAvailableBitrates(result.availableBitrates);
      setSegmentLength(result.segmentLength - 0.3);
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    let getStatusIntervalId = null;

    getMetrics();
    getStatusIntervalId = setInterval(getMetrics, FETCH_INTERVAL); // runs every 1 min.

    // returned function will be called on component unmount
    return () => {
      clearInterval(getStatusIntervalId);
    };
  }, []);

  const noData = (
    <div>
      <Typography.Title>Stream Performance</Typography.Title>
      <Typography.Paragraph>
        Data has not yet been collected. Once a stream has begun and viewers are watching this page
        will be available.
      </Typography.Paragraph>
    </div>
  );
  if (!errors?.length) {
    return noData;
  }

  if (!medianLatency?.length) {
    return noData;
  }

  if (!medianSegmentDownloadDurations?.length) {
    return noData;
  }

  const errorChart = [
    {
      name: 'Errors',
      color: '#B63FFF',
      options: { radius: 3 },
      data: errors,
    },
    {
      name: 'Quality changes',
      color: '#2087E2',
      options: { radius: 2 },
      data: qualityVariantChanges,
    },
  ];

  const latencyChart = [
    {
      name: 'Median stream latency',
      color: '#00FFFF',
      options: { radius: 2 },
      data: medianLatency,
    },
    {
      name: 'Lowest stream latency',
      color: '#02FD0D',
      options: { radius: 2 },
      data: lowestLatency,
    },
    {
      name: 'Highest stream latency',
      color: '#B63FFF',
      options: { radius: 2 },
      data: highestLatency,
    },
  ];

  const segmentDownloadDurationChart = [
    {
      name: 'Median download duration',
      color: '#00FFFF',
      options: { radius: 2 },
      data: medianSegmentDownloadDurations,
    },
    {
      name: 'Max download duration',
      color: '#B63FFF',
      options: { radius: 2 },
      data: maximumSegmentDownloadDurations,
    },
    {
      name: 'Min download duration',
      color: '#02FD0D',
      options: { radius: 2 },
      data: minimumSegmentDownloadDurations,
    },
    {
      name: `Approximate limit`,
      color: '#003FFF',
      data: medianSegmentDownloadDurations.map(item => ({
        time: item.time,
        value: segmentLength,
      })),
      options: { radius: 0 },
    },
  ];

  const bitrateChart = [
    {
      name: 'Lowest viewer bitrate',
      color: '#B63FFF',
      data: minimumPlayerBitrate,
      options: { radius: 2 },
    },
    {
      name: 'Median viewer bitrate',
      color: '#00FFFF',
      data: medianPlayerBitrate,
      options: { radius: 2 },
    },
    {
      name: 'Maximum viewer bitrate',
      color: '#02FD0D',
      data: maximumPlayerBitrate,
      options: { radius: 2 },
    },
  ];

  availableBitrates.forEach(bitrate => {
    bitrateChart.push({
      name: `${bitrate} kbps stream`,
      color: '#003FFF',
      data: minimumPlayerBitrate.map(item => ({
        time: item.time,
        value: 1200,
      })),
      options: { radius: 0 },
    });
  });

  const currentSpeed = bitrateChart[0]?.data[bitrateChart[0].data.length - 1]?.value;
  const currentDownloadSeconds =
    medianSegmentDownloadDurations[medianSegmentDownloadDurations.length - 1]?.value;
  const lowestVariant = availableBitrates[0]; // TODO: get lowest bitrate from available bitrates

  const latencyMedian = medianLatency[medianLatency.length - 1]?.value || 0;
  const latencyMax = highestLatency[highestLatency.length - 1]?.value || 0;
  const latencyMin = lowestLatency[lowestLatency.length - 1]?.value || 0;
  const latencyStat = (Number(latencyMax) + Number(latencyMin) + Number(latencyMedian)) / 3;

  let recentErrorCount = 0;
  const errorValueCount = errorChart[0]?.data.length || 0;
  if (errorValueCount > 5) {
    const values = errorChart[0].data.slice(-5);
    recentErrorCount = values.reduce((acc, curr) => acc + Number(curr.value), 0);
  } else {
    recentErrorCount = errorChart[0].data.reduce((acc, curr) => acc + Number(curr.value), 0);
  }
  const showStats = currentSpeed > 0 || currentDownloadSeconds > 0 || recentErrorCount > 0;
  let bitrateError = null;
  let speedError = null;

  if (currentSpeed !== 0 && currentSpeed < lowestVariant) {
    bitrateError = `At least one of your viewers is playing your stream at ${currentSpeed}kbps, slower than ${lowestVariant}kbps, the lowest quality you made available, experiencing buffering. Consider adding a lower quality with a lower bitrate if the errors over time warrant this.`;
  }

  if (currentDownloadSeconds > segmentLength) {
    speedError =
      'Your viewers may be consuming your video slower than required. This may be due to slow networks or your latency configuration. Consider adding a lower quality with a lower bitrate or experiment with increasing the latency buffer setting.';
  }

  const errorStatColor = recentErrorCount > 0 ? '#B63FFF' : '#FFFFFF';
  const statStyle = {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    height: '80px',
  };
  return (
    <>
      <Typography.Title>Stream Performance</Typography.Title>
      <Typography.Paragraph>
        This tool hopes to help you identify and troubleshoot problems you may be experiencing with
        your stream. It aims to aggregate experiences across your viewers, meaning one viewer with
        an exceptionally bad experience may throw off numbers for the whole, especially with a low
        number of viewers.
      </Typography.Paragraph>
      <Typography.Paragraph>
        The data is only collected by those using the Owncast web interface and is unable to gain
        insight into external players people may be using such as VLC, MPV, QuickTime, etc.
      </Typography.Paragraph>
      <Space direction="vertical" size="middle">
        <Row
          gutter={[16, 16]}
          justify="space-around"
          style={{ display: showStats ? 'flex' : 'none' }}
        >
          <Col>
            <Card type="inner">
              <div style={statStyle}>
                <Statistic
                  title="Viewer Playback Speed"
                  value={`${currentSpeed}`}
                  prefix={<WifiOutlined style={{ marginRight: '5px' }} />}
                  precision={0}
                  suffix="kbps"
                />
              </div>
            </Card>
          </Col>
          <Col>
            <Card type="inner">
              <div style={statStyle}>
                <Statistic
                  title="Viewer Latency"
                  value={`${latencyStat}`}
                  prefix={<ClockCircleOutlined style={{ marginRight: '5px' }} />}
                  precision={0}
                  suffix="seconds"
                />
              </div>
            </Card>
          </Col>
          <Col>
            <Card type="inner">
              <div style={statStyle}>
                <Statistic
                  title="Recent Playback Errors"
                  value={`${recentErrorCount || 0}`}
                  valueStyle={{ color: errorStatColor }}
                  prefix={<WarningOutlined style={{ marginRight: '5px' }} />}
                  suffix=""
                />{' '}
              </div>
            </Card>
          </Col>
        </Row>

        <Card>
          <DescriptionBox
            title="Video Segment Download"
            description={
              <>
                <Typography.Paragraph>
                  Once a video segment takes too long to download a viewer will experience
                  buffering. If you see slow downloads you should offer a lower quality for your
                  viewers, or find other ways, possibly an external storage provider, a CDN or a
                  faster network, to improve your stream quality. Increasing your latency buffer can
                  also help for some viewers.
                </Typography.Paragraph>
                <Typography.Paragraph>
                  In short, once the pink line consistently gets near the blue line, your stream is
                  likely experiencing problems for viewers.
                </Typography.Paragraph>
              </>
            }
          />
          {speedError && (
            <Alert message="Slow downloads" description={speedError} type="error" showIcon />
          )}
          <Chart
            title="Seconds"
            dataCollections={segmentDownloadDurationChart}
            color="#FF7700"
            unit="s"
            yLogarithmic
          />
        </Card>
        <Card>
          <DescriptionBox
            title="Player Network Speed"
            description={
              <>
                <Typography.Paragraph>
                  The playback bitrate of your viewers. Once somebody's bitrate drops below the
                  lowest video variant bitrate they will experience buffering. If you see viewers
                  with slow connections trying to play your video you should consider offering an
                  additional, lower quality.
                </Typography.Paragraph>
                <Typography.Paragraph>
                  In short, once the pink line gets near the lowest blue line, your stream is likely
                  experiencing problems for at least one of your viewers.
                </Typography.Paragraph>
              </>
            }
          />
          {bitrateError && (
            <Alert
              message="Low bandwidth viewers"
              description={bitrateError}
              type="error"
              showIcon
            />
          )}
          <Chart
            title="Lowest Player Bitrate"
            dataCollections={bitrateChart}
            color="#FF7700"
            unit="kbps"
            yLogarithmic
          />
        </Card>
        <Card>
          <DescriptionBox
            title="Errors and Quality Changes"
            description={
              <>
                <Typography.Paragraph>
                  Recent number of errors, including buffering, and quality changes from across all
                  your viewers. Errors can occur for many reasons, including browser issues,
                  plugins, wifi problems, and they don't all represent fatal issues or something you
                  have control over.
                </Typography.Paragraph>
                A quality change is not necessarily a negative thing, but if it's excessive and
                coinciding with errors you should consider adding another quality variant.
                <Typography.Paragraph />
              </>
            }
          />
          <Chart title="#" dataCollections={errorChart} color="#FF7700" unit="" />
        </Card>
        <Card>
          <DescriptionBox
            title="Viewer Latency"
            description="An approximate number of seconds that your viewers are behind your live video. The more people buffer the further behind they will be. High latency itself is not a problem, but some people care about this more than others."
          />
          <Chart title="Seconds" dataCollections={latencyChart} color="#FF7700" unit="s" />
        </Card>
      </Space>
    </>
  );
}
