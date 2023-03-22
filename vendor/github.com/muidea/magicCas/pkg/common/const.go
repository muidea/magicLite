package common

const (
	v1Version      = "/api/v1"
	defaultVersion = v1Version
)

const ApiVersion = defaultVersion

const (
	CreateEventMask = "/#/create/"
	DeleteEventMask = "/#/delete/+"
	UpdateEventMask = "/#/update/+"
	NotifyEventMask = "/#/notify/+"
)

const (
	Create = iota
	Delete
	Update
	Load
)

const (
	AuthSession    = "authSession"
	AuthNamespace  = "X-Namespace"
	AuthEntity     = "authEntity"
	AuthRole       = "authRole"
	AuthEntityView = "authEntityView"
	AuthRoleView   = "authRoleView"
)
