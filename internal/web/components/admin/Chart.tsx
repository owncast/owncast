import format from 'date-fns/format';
import { FC } from 'react';

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
} from 'chart.js';
import { Line } from 'react-chartjs-2';

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
};

function createGraphDataset(dataArray) {
  const dataValues = {};
  dataArray.forEach(item => {
    const dateObject = new Date(item.time);
    const dateString = format(dateObject, 'H:mma');
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
}) => {
  const renderData = [];

  if (data && data.length > 0) {
    renderData.push({
      id: title,
      label: title,
      backgroundColor: color,
      borderColor: color,
      borderWidth: 3,
      data: createGraphDataset(data),
    });
  }

  dataCollections.forEach(collection => {
    renderData.push({
      id: collection.name,
      label: collection.name,
      data: createGraphDataset(collection.data),
      backgroundColor: collection.color,
      borderColor: collection.color,
      borderWidth: 3,
      pointStyle: collection.pointStyle || 'circle',
      radius: collection.pointRadius || 1,
    });
  });

  const options = {
    responsive: true,

    scales: {
      y: {
        type: yLogarithmic ? ('logarithmic' as const) : ('linear' as const),
        reverse: yFlipped,
        title: {
          display: true,
          text: unit,
        },
      },
    },
  };

  return (
    <div className="line-chart-container">
      <Line data={{ datasets: renderData }} options={options} height="70vh" />
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
