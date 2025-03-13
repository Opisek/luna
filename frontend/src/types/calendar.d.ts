type CalendarModel = {
  id: string;
  source: string;
  name: string;
  desc: string;
  color: string;
}

type CalendarModelChanges = {
  name: boolean;
  desc: boolean;
  color: boolean;
}