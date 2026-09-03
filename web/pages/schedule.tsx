import type { ReactElement } from 'react';
import { ScheduleEventPage } from '../components/ui/ScheduleEventPage/ScheduleEventPage';

export default function Schedule() {
  return <ScheduleEventPage />;
}

Schedule.getLayout = function getLayout(page: ReactElement) {
  return page;
};
