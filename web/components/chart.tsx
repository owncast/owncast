import ChartJs from 'chart.js/auto';
import Chartkick from 'chartkick';
import format from 'date-fns/format';
import { LineChart } from 'react-chartkick';

// from https://github.com/ankane/chartkick.js/blob/master/chart.js/chart.esm.js
Chartkick.use(ChartJs);

interface TimedValue {
  time: Date;
  value: number;
}

interface ChartProps {
  data?: TimedValue[];
  title?: string;
  color: string;
  unit: string;
  yFlipped?: boolean;
  yLogarithmic?: boolean;
  dataCollections?: any[];
}

function createGraphDataset(dataArray) {
  const dataValues = {};
  dataArray.forEach(item => {
    const dateObject = new Date(item.time);
    const dateString = format(dateObject, 'H:mma');
    dataValues[dateString] = item.value;
  });
  return dataValues;
}

export default function Chart({
  data,
  title,
  color,
  unit,
  dataCollections,
  yFlipped,
  yLogarithmic,
}: ChartProps) {
  const renderData = [];

  if (data && data.length > 0) {
    renderData.push({
      name: title,
      color,
      data: createGraphDataset(data),
    });
  }

  dataCollections.forEach(collection => {
    renderData.push({
      name: collection.name,
      data: createGraphDataset(collection.data),
      color: collection.color,
      dataset: collection.options,
    });
  });

  // ChartJs.defaults.scales.linear.reverse = true;

  const options = {
    scales: {
      y: { reverse: false, type: 'linear' },
      x: {
        type: 'time',
      },
    },
  };

  options.scales.y.reverse = yFlipped;
  options.scales.y.type = yLogarithmic ? 'logarithmic' : 'linear';

  return (
    <div className="line-chart-container">
      <LineChart
        xtitle="Time"
        ytitle={title}
        suffix={unit}
        legend="bottom"
        color={color}
        data={renderData}
        download={title}
        library={options}
      />
    </div>
  );
}

Chart.defaultProps = {
  dataCollections: [],
  data: [],
  title: '',
  yFlipped: false,
  yLogarithmic: false,
};
