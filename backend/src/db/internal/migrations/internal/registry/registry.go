package registry

import (
	"luna-backend/common"
	"luna-backend/db/internal/migrations/types"
)

type MigrationFunc func(*types.MigrationQueries) error

type Migration struct {
	Ver common.Version
	Fun MigrationFunc
}

type MigrationRegistry struct {
	migrations [][][]*Migration
}

var registry *MigrationRegistry

func GetRegistry() *MigrationRegistry {
	if registry == nil {
		registry = &MigrationRegistry{}
	}
	return registry
}

func RegisterMigration(version common.Version, fun MigrationFunc) {
	reg := GetRegistry()
	migrations := reg.migrations

	for len(migrations) <= version.Major {
		migrations = append(migrations, [][]*Migration{})
	}
	majorMigrations := migrations[version.Major]

	for len(majorMigrations) <= version.Minor {
		majorMigrations = append(majorMigrations, []*Migration{})
	}
	migrations[version.Major] = majorMigrations
	minorMigrations := majorMigrations[version.Minor]

	for len(minorMigrations) <= version.Patch {
		minorMigrations = append(minorMigrations, nil)
	}
	migrations[version.Major][version.Minor] = minorMigrations
	migrations[version.Major][version.Minor][version.Patch] = &Migration{Ver: version, Fun: fun}

	reg.migrations = migrations
}

func GetMigrations(lastVersion common.Version) []*Migration {
	reg := GetRegistry()
	migrations := reg.migrations
	selectedMigrations := []*Migration{}

	for major := lastVersion.Major; major < len(migrations); major++ {
		for minor := lastVersion.Minor; minor < len(migrations[major]); minor++ {
			for patch := lastVersion.Patch; patch < len(migrations[major][minor]); patch++ {
				migration := migrations[major][minor][patch]
				if migration == nil || major == lastVersion.Major && minor == lastVersion.Minor && patch == lastVersion.Patch {
					continue
				}
				selectedMigrations = append(selectedMigrations, migration)
			}
		}
	}

	return selectedMigrations
}
