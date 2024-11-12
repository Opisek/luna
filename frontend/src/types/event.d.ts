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
  }
};