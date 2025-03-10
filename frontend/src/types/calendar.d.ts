type CalendarModel = {
  id: string;
  source: string;
  name: string;
  url: string;
  color: string;
}

type CalendarModelChanges = {
  name: boolean;
  color: boolean;
}