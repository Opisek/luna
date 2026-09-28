type CalendarModel = {
  id: string;
  source: string;
  name: string;
  desc: string;
  color: string;
  overridden: boolean;
  can_edit: boolean;
  can_delete: boolean;
  can_add_events: boolean;
}

type CalendarModelChanges = {
  name: boolean;
  desc: boolean;
  color: boolean;
}