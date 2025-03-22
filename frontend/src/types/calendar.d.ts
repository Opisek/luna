type CalendarModel = {
  id: string;
  source: string;
  name: string;
  desc: string;
  color: string;
  overridden: boolean;
}

type CalendarModelChanges = {
  name: boolean;
  desc: boolean;
  color: boolean;
}