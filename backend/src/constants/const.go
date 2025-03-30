package constants

import "time"

const LifetimeCacheSoft = 5 * time.Minute
const LifetimeCacheHard = 24 * time.Hour
const MaxFormBytes = 50 * 1024 * 1024 // 50MB
