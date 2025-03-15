export const EmptySource: SourceModel = {
  id: '',
  name: '',
  type: 'caldav',
  settings: {
    url: '',
  },
  auth_type: 'none',
  auth: {
    username: '',
    password: '',
    token: '',
  }
}

export const AllChangesSource: SourceModelChanges = {
  name: true,
  type: true,
  settings: true,
  auth: true
}

export const EmptyCalendar: CalendarModel = {
  id: "",
  source: "",
  name: "",
  desc: "",
  color: ""
}

export const AllChangesCalendar: CalendarModelChanges = {
  name: true,
  desc: true,
  color: true
}

export const PlaceholderDate = new Date(0);

export const EmptyEvent: EventModel = {
  id: '',
  calendar: '',
  name: '',
  desc: '',
  color: '',
  date: {
    start: PlaceholderDate,
    end: PlaceholderDate,
    allDay: false,
    recurrence: false
  }
}

export const AllChangesEvent: EventModelChanges = {
  name: true,
  desc: true,
  color: true,
  date: true
}

export const NoOp = () => {};

export const EmptyOption: Option = {
  value: '',
  name: ''
}