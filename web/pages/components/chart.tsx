import { LineChart, XAxis, YAxis, Line, Tooltip, Legend } from "recharts";
import { timeFormat } from "d3-time-format";
import useWindowSize from '../../utils/hook-windowresize';
import styles from '../../styles/styles.module.css';

interface ToolTipProps {
  active?: boolean,
  payload?: {name: string, payload: {value: string, time: Date}}[],
  unit?: string
}

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

function CustomizedTooltip(props: ToolTipProps) {
  const { active, payload, unit } = props;
  if (active && payload && payload[0]) {
    const time = payload[0].payload ? timeFormat("%I:%M")(new Date(payload[0].payload.time)) : "";

    const tooltipDetails = payload.map(data => {
      return <div className="label" key={data.name}>
        {data.payload.value}{unit} {data.name}
        </div>
    });
    return (
      <span className="custom-tooltip">
        <strong>{time}</strong>
        {tooltipDetails}
      </span>
    );
  }
  return null;
}
CustomizedTooltip.defaultProps = defaultProps;

export default function Chart({ data, title, color, unit, dataCollections }: ChartProps) {
  if (!data && !dataCollections) {
    return null;
  }
  
  const windowSize = useWindowSize();
  const chartWidth = windowSize.width * .68;
  const chartHeight = chartWidth * .333;

  const timeFormatter = (tick: string) => {
    return timeFormat("%I:%M")(new Date(tick));
  };

  let ticks = [];
  if (dataCollections.length > 0) {
    ticks = dataCollections[0].data?.map((collection) => {
      return collection?.time;
    })
  } else if (data?.length > 0){
    ticks = data?.map(item => {
      return item?.time;
    });
  }

  const line = data && data?.length > 0 ? (
    <Line
    key="line"
    type="basis"
    dataKey="value"
    stroke={color}
    dot={null}
    strokeWidth={3}
    legendType="square"
    name={title}
    data={data}
  />
  ) : null;

  return (
    <div className={styles.lineChartContainer}>
      <LineChart width={chartWidth} height={chartHeight}>
        <XAxis
          dataKey="time"
          tickFormatter={timeFormatter}
          interval="preserveStartEnd"
          tickCount={5}
          minTickGap={15}
          domain={["dataMin", "dataMax"]}
          ticks={ticks}
        />
        <YAxis
          dataKey="value"
          interval="preserveStartEnd"
          unit={unit}
          domain={["dataMin", "dataMax"]}
        />
        <Tooltip content={<CustomizedTooltip unit={unit} />} />
        <Legend />
        {line}
        {dataCollections?.map((s) => (
          <Line
            dataKey="value"
            data={s.data}
            name={s.name}
            key={s.name}
            type="basis"
            stroke={s.color}
            dot={null}
            strokeWidth={3}
            legendType="square"
          />
        ))}
      </LineChart>
    </div>
  );
}

Chart.defaultProps = {
  dataCollections: [],
  data: [],
  title: '',
};
