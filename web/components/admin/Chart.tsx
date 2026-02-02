import { format } from 'date-fns';
import { FC, useRef } from 'react';

import { DownloadOutlined } from '@ant-design/icons';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  LogarithmicScale,
  ChartOptions,
} from 'chart.js';
import { Line } from 'react-chartjs-2';
import { Button } from 'antd';

ChartJS.register(
  CategoryScale,
  LogarithmicScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
);

export enum TimeWindow {
  Hour12 = 1,
  Hour24 = 2,
  Day7 = 3,
  Day30 = 4,
  Month3 = 5,
  Month6 = 6,
}

interface TimedValue {
  time: Date;
  value: number;
  pointStyle?: boolean | string;
  pointRadius?: number;
}

export type ChartProps = {
  data?: TimedValue[];
  title?: string;
  color: string;
  unit: string;
  yFlipped?: boolean;
  yLogarithmic?: boolean;
  dataCollections?: any[];
  minYValue?: number;
  yStepSize?: number;
  timeWindowKey?: TimeWindow;
};

function getLabelFormat(timeWindowKey?: TimeWindow) {
  switch (timeWindowKey) {
    case TimeWindow.Hour12:
    case TimeWindow.Hour24:
      return 'H:mm'; // 12/24h: hour

    case TimeWindow.Day7:
    case TimeWindow.Day30:
      return 'MMM d'; // 7/30d: day

    case TimeWindow.Month3:
    case TimeWindow.Month6:
      return 'MMM'; // 3/6mo: month

    default:
      return 'H:mm';
  }
}

function createGraphDataset(dataArray, labelFormat) {
  const dataValues = {};
  dataArray.forEach(item => {
    const dateObject = new Date(item.time);
    const dateString = format(dateObject, labelFormat);
    dataValues[dateString] = item.value;
  });
  return dataValues;
}

export const Chart: FC<ChartProps> = ({
  data,
  title,
  color,
  unit,
  dataCollections,
  yFlipped,
  yLogarithmic,
  minYValue,
  yStepSize = 0,
  timeWindowKey,
}) => {
  const renderData = [];
  const chartRef = useRef(null);

  const labelFormat = getLabelFormat(timeWindowKey);

  const downloadChart = () => {
    if (chartRef.current) {
      const link = document.createElement('a');
      link.download = 'chart.png';
      link.href = chartRef.current.canvas.toDataURL();
      link.click();
    }
  };

  if (data && data.length > 0) {
    renderData.push({
      id: title,
      label: title,
      backgroundColor: color,
      borderColor: color,
      borderWidth: 3,
      data: createGraphDataset(data, labelFormat),
    });
  }

  dataCollections.forEach(collection => {
    renderData.push({
      id: collection.name,
      label: collection.name,
      data: createGraphDataset(collection.data, labelFormat),
      backgroundColor: collection.color,
      borderColor: collection.color,
      borderWidth: 3,
      pointStyle: collection.pointStyle || 'circle',
      radius: collection.pointRadius || 1,
    });
  });

  const options = {
    responsive: true,
    clip: false,
    scales: {
      x: {
        title: { display: false },
        ticks: {
          autoSkip: true,
          maxTicksLimit: 10,
        },
      },
      y: {
        type: yLogarithmic ? ('logarithmic' as const) : ('linear' as const),
        reverse: yFlipped,
        min: minYValue,
        ticks: {
          stepSize: yStepSize,
        },
        title: {
          display: true,
          text: unit,
        },
      },
    },
  };

  return (
    <div className="line-chart-container">
      <Line
        ref={chartRef}
        data={{ datasets: renderData }}
        options={options as ChartOptions<'line'>}
        height="70vh"
      />
      <Button
        size="small"
        onClick={downloadChart}
        type="ghost"
        icon={<DownloadOutlined />}
        className="download-btn"
      />
    </div>
  );
};

Chart.defaultProps = {
  dataCollections: [],
  data: [],
  title: '',
  yFlipped: false,
  yLogarithmic: false,
};
