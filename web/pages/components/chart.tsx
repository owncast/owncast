import { LineChart, XAxis, YAxis, Line, Tooltip, Legend } from "recharts";
import { timeFormat } from "d3-time-format";

export default function Chart({ data, color, unit }) {
  const CustomizedTooltip = (props) => {
    const { active, payload } = props;
    if (active && payload && payload[0]) {
      const time = payload[0].payload
        ? timeFormat("%I:%M")(new Date(payload[0].payload.time), {
            nearestTo: 1,
          })
        : "";
      return (
        <div className="custom-tooltip">
          <p className="label">
            <strong>{time}</strong> {payload[0].payload.value} %
          </p>
        </div>
      );
    }
    return null;
  };

  const timeFormatter = (tick) => {
    return timeFormat("%I:%M")(new Date(tick), {
      nearestTo: 1,
    });
  };

  return (
    <LineChart width={1200} height={400} data={data}>
      <XAxis
        dataKey="time"
        // type="number"
        tickFormatter={timeFormatter}
        interval="preserveStartEnd"
        tickCount={5}
        minTickGap={15}
        domain={["dataMin", "dataMax"]}
      />
      <YAxis
        dataKey="value"
        interval="preserveStartEnd"
        unit={unit}
        domain={["dataMin", "dataMax"]}
      />
      <Tooltip content={<CustomizedTooltip />} />
      <Legend />
      <Line
        type="monotone"
        dataKey="value"
        stroke={color}
        dot={null}
        strokeWidth={3}
      />
    </LineChart>
  );
}
