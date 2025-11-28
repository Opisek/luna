package perms

type Permission string

const (
	ReadEvents           Permission = "read_events"
	AddEvents            Permission = "add_events"
	EditEvents           Permission = "edit_events"
	DeleteEvents         Permission = "delete_events"
	ReadCalendars        Permission = "read_calendars"
	AddCalendars         Permission = "add_calendars"
	EditCalendars        Permission = "edit_calendars"
	DeleteCalendars      Permission = "delete_calendars"
	ReadSources          Permission = "read_sources"
	AddSources           Permission = "add_sources"
	EditSources          Permission = "edit_sources"
	DeleteSources        Permission = "delete_sources"
	ManageInvites        Permission = "manage_invites"
	ManageUsers          Permission = "manage_users"
	ManageGlobalSettings Permission = "manage_global_settings"
	ManageUserSettings   Permission = "manage_user_settings"
	ManageSessions       Permission = "manage_sessions"
	ManageResources      Permission = "manage_resources"
	Administrative       Permission = "administrative"
)

// Administrative:
// Whether the token is allowed to perform actions that require admin privileges.
// This does not automatically grant all other permissions.
// This does not replace checking whether the user is an administrator.
// To check whether some action xyz that requires admin privileges amongst other privileges is allowed:
// 1. Check whether the user is an administrator
// 2. Check whether the token has administrative permission
// 3. Check whether the token has additional permissions

var allPermList = []Permission{
	ReadEvents,
	AddEvents,
	EditEvents,
	DeleteEvents,
	ReadCalendars,
	AddCalendars,
	EditCalendars,
	DeleteCalendars,
	ReadSources,
	AddSources,
	EditSources,
	DeleteSources,
	ManageInvites,
	ManageUsers,
	ManageGlobalSettings,
	ManageUserSettings,
	ManageSessions,
	ManageResources,
	Administrative,
}

type TokenPermissions struct {
	hasAll      bool
	permissions map[Permission]bool
}

func AllPermissions() *TokenPermissions {
	return &TokenPermissions{
		hasAll:      true,
		permissions: map[Permission]bool{},
	}
}

func FromList(perms []Permission) *TokenPermissions {
	permMap := make(map[Permission]bool)
	for _, p := range perms {
		permMap[p] = true
	}

	return &TokenPermissions{
		hasAll:      false,
		permissions: permMap,
	}
}

func (tp *TokenPermissions) ToList() []Permission {
	if tp.hasAll {
		list := make([]Permission, len(allPermList))
		copy(list, allPermList)
		return list
	}
	perms := make([]Permission, 0, len(tp.permissions))
	for p := range tp.permissions {
		perms = append(perms, p)
	}
	return perms
}

func (tp *TokenPermissions) Has(perm Permission) bool {
	if tp.hasAll {
		return true
	}
	granted, ok := tp.permissions[perm]
	return ok && granted
}
