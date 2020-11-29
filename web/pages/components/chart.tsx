import { LineChart } from 'react-chartkick'
import styles from '../../styles/styles.module.css';
import 'chart.js';

const defaultProps = {
  active: false,
  payload: Object,
  unit: '',
};

interface TimedValue {
  time: Date;
  value: number;
}

interface ChartProps {
  data?: TimedValue[],
  title?: string,
  color: string,
  unit: string,
  dataCollections?: any[],
}

function createGraphDataset(dataArray) {
  const dataValues = {};
  dataArray.forEach(item => {
    const dateObject = new Date(item.time);
    const dateString = `${dateObject.getFullYear()  }-${  dateObject.getMonth()  }-${  dateObject.getDay()  } ${  dateObject.getHours()  }:${  dateObject.getMinutes()}`;
    dataValues[dateString] = item.value;
  })
  return dataValues;
}

export default function Chart({ data, title, color, unit, dataCollections }: ChartProps) {
  const renderData = [];

  if (data && data.length > 0) {
    renderData.push({
      name: title,
      color,
      data: createGraphDataset(data)
    });
  }

  dataCollections.forEach(collection => {
    renderData.push(
      {name: collection.name, data: createGraphDataset(collection.data), color: collection.color}
    )
  });

  return (
    <div className={styles.lineChartContainer}>
      <LineChart
        xtitle="Time"
        ytitle={title}
        suffix={unit}
        legend="bottom"
        color={color}
        data={renderData}
        download={title}
      />
    </div>
  );
}

Chart.defaultProps = {
  dataCollections: [],
  data: [],
  title: '',
};
