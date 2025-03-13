type EventModel = {
  id: string;
  calendar: string;
  name: string;
  desc: string;
  color: string;
  date: {
    start: Date;
    end: Date;
    allDay: boolean;
    recurrence: any;
  }
};

type EventModelChanges = {
  name: boolean;
  desc: boolean;
  color: boolean;
  date: boolean;
}