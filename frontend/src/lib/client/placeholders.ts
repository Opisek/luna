import type { Option } from "../../types/options"
import { PermissionKeys } from "../../types/permissions"

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

export const NoChangesSource: SourceModelChanges = {
  name: false,
  type: false,
  settings: false,
  auth: false
}

export const EmptyCalendar: CalendarModel = {
  id: "",
  source: "",
  name: "",
  desc: "",
  color: "",
  overridden: false
}

export const DefaultTokenPermissions = [
  PermissionKeys.ReadEvents,
  PermissionKeys.AddEvents,
  PermissionKeys.EditEvents,
  PermissionKeys.DeleteEvents,
  PermissionKeys.ReadCalendars,
  PermissionKeys.AddCalendars,
  PermissionKeys.EditCalendars,
  PermissionKeys.DeleteCalendars,
  PermissionKeys.ReadSources,
  PermissionKeys.AddSources,
  PermissionKeys.EditSources,
  PermissionKeys.DeleteSources
]

export const EmptySession: Session = {
  session_id: "",
  user_id: "",
  created_at: new Date(0),
  last_seen: new Date(0),
  user_agent: "",
  initial_ip_address: "",
  last_ip_address: "",
  is_short_lived: false,
  is_api: false,
  location: "",
  permissions: DefaultTokenPermissions
}

export const EmptyRegistrationInvite: RegistrationInvite = {
  invite_id: "",
  author: "",
  email: "",
  created_at: new Date(0),
  expires_at: new Date(0),
  code: ""
}

export const AllChangesCalendar: CalendarModelChanges = {
  name: true,
  desc: true,
  color: true
}

export const NoChangesCalendar: CalendarModelChanges = {
  name: false,
  desc: false,
  color: false
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
  },
  overridden: false
}

export const AllChangesEvent: EventModelChanges = {
  name: true,
  desc: true,
  color: true,
  date: true
}

export const NoChangesEvent: EventModelChanges = {
  name: false,
  desc: false,
  color: false,
  date: false
}

export const NoOp = () => {};
export const AsyncNoOp = async () => {};

export const EmptyOption: Option<string> = {
  value: '',
  name: ''
}