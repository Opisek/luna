type EventModel = {
  id: string;
  calendar: string;
  name: string;
  desc?: string;
  color: string;
  date: {
    start: Date;
    end: Date;
    allDay: boolean;
    recurrence?: {
      RRULE?: string,
      RDATE?: string,
      EXDATE?: string
    };
  };
  overridden: boolean;
  can_edit: boolean;
  can_delete: boolean;
};

type EventModelChanges = {
  name?: boolean;
  desc?: boolean;
  color?: boolean;
  date?: boolean;
}