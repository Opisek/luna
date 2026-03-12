package constants

import "time"

const LifetimeRamCache = 1 * time.Minute

const LifetimeFileCacheSoft = 5 * time.Minute
const LifetimeFileCacheHard = 24 * time.Hour

const MaxInviteDuration = 7 * 24 * time.Hour // 7 days

const MaxFormBytes = 50 * 1024 * 1024 // 50MB
