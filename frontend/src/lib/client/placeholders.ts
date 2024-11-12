export const EmptySource: SourceModel = {
  id: '',
  name: '',
  type: '',
  settings: {},
  auth_type: '',
  auth: {},
  collapsed: false
}

export const EmptyEvent: EventModel = {
  id: '',
  calendar: '',
  name: '',
  desc: '',
  color: '',
  date: {
    start: new Date(),
    end: new Date(),
    allDay: false,
  }
}

export const NoOp = () => {};

export const PlaceholderDate = new Date(0);

export const EmptyOption: Option = {
  value: '',
  name: ''
}