import styles from '../../styles/styles.module.css';
import { LineChart } from 'react-chartkick'
import 'chart.js'

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

export default function Chart({ data, title, color, unit, dataCollections }: ChartProps) {
  var renderData = [];

  if (data && data.length > 0) {
    renderData.push({
      name: title,
      color: color,
      data: createGraphDatasetFromObject(data)
    });
  }

  dataCollections.forEach(collection => {
    renderData.push(
      {name: collection.name, data: createGraphDatasetFromObject(collection.data), color: collection.color}
    )
  });

  return (
    <LineChart
    className={styles.lineChartContainer}
    xtitle="Time"
    ytitle={title}
    suffix={unit}
    legend={"bottom"}
    color={color}
    data={renderData}
    download={title}
    />
  )
}

Chart.defaultProps = {
  dataCollections: [],
  data: [],
  title: '',
};

function createGraphDatasetFromObject(dataArray) {
  var dataValues = {};
  dataArray.forEach(item => {
    const dateObject = new Date(item.time);
    const dateString = dateObject.getFullYear() + '-' + dateObject.getMonth() + '-' + dateObject.getDay() + ' ' + dateObject.getHours() + ':' + dateObject.getMinutes();
    dataValues[dateString] = item.value;
  })
  return dataValues;
}