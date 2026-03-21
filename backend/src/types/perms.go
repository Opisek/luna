package types

type Permission string

const (
	PermReadEvents           Permission = "read_events"
	PermAddEvents            Permission = "add_events"
	PermEditEvents           Permission = "edit_events"
	PermDeleteEvents         Permission = "delete_events"
	PermReadCalendars        Permission = "read_calendars"
	PermAddCalendars         Permission = "add_calendars"
	PermEditCalendars        Permission = "edit_calendars"
	PermDeleteCalendars      Permission = "delete_calendars"
	PermReadSources          Permission = "read_sources"
	PermAddSources           Permission = "add_sources"
	PermEditSources          Permission = "edit_sources"
	PermDeleteSources        Permission = "delete_sources"
	PermManageInvites        Permission = "manage_invites"
	PermManageUsers          Permission = "manage_users"
	PermManageOauthClients   Permission = "manage_oauth_clients"
	PermManageGlobalSettings Permission = "manage_global_settings"
	PermManageUserSettings   Permission = "manage_user_settings"
	PermManageSessions       Permission = "manage_sessions"
	PermManageResources      Permission = "manage_resources"
	PermAdministrative       Permission = "administrative"
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
	PermReadEvents,
	PermAddEvents,
	PermEditEvents,
	PermDeleteEvents,
	PermReadCalendars,
	PermAddCalendars,
	PermEditCalendars,
	PermDeleteCalendars,
	PermReadSources,
	PermAddSources,
	PermEditSources,
	PermDeleteSources,
	PermManageInvites,
	PermManageUsers,
	PermManageOauthClients,
	PermManageGlobalSettings,
	PermManageUserSettings,
	PermManageSessions,
	PermManageResources,
	PermAdministrative,
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

func TokenPermsFromPermsList(perms []Permission) *TokenPermissions {
	permMap := make(map[Permission]bool)
	for _, p := range perms {
		permMap[p] = true
	}

	return &TokenPermissions{
		hasAll:      false,
		permissions: permMap,
	}
}

func TokenPermsFromStringList(permStrs []string) *TokenPermissions {
	permSet := make(map[Permission]bool)
	for _, perm := range allPermList {
		permSet[perm] = true
	}

	validPerms := make([]Permission, 0, len(permStrs))
	for _, perm := range permStrs {
		if _, ok := permSet[Permission(perm)]; ok {
			validPerms = append(validPerms, Permission(perm))
		}
	}

	return TokenPermsFromPermsList(validPerms)
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

func (tp *TokenPermissions) Equals(other *TokenPermissions) bool {
	if tp.hasAll != other.hasAll {
		return false
	}
	if len(tp.permissions) != len(other.permissions) {
		return false
	}
	for perm, granted := range tp.permissions {
		otherGranted, ok := other.permissions[perm]
		if !ok || granted != otherGranted {
			return false
		}
	}
	return true
}
