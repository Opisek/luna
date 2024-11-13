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
  },
  collapsed: false
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
  }
}

export const NoOp = () => {};

export const EmptyOption: Option = {
  value: '',
  name: ''
}