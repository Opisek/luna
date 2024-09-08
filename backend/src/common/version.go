package common

import "fmt"

type Version struct {
	Major     int
	Minor     int
	Patch     int
	Extension string
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d%s", v.Major, v.Minor, v.Patch, v.Extension)
}

func (v1 *Version) IsGreaterThan(v2 *Version) bool {
	// Major
	if v1.Major > v2.Major {
		return true
	}
	if v1.Major < v2.Major {
		return false
	}

	// Minor
	if v1.Minor > v2.Minor {
		return true
	}
	if v1.Minor < v2.Minor {
		return false
	}

	// Patch
	if v1.Patch > v2.Patch {
		return true
	}
	if v1.Patch < v2.Patch {
		return false
	}

	// Equal implies not greater than
	return false
}

func (v1 *Version) IsEqualTo(v2 *Version) bool {
	return v1.Major == v2.Major && v1.Minor == v2.Minor && v1.Patch == v2.Patch
}

func ParseVersion(verstr string) (Version, error) {
	ver := Version{}

	_, err := fmt.Sscanf(verstr, "%d.%d.%d%s", &ver.Major, &ver.Minor, &ver.Patch, &ver.Extension)
	if err != nil {
		_, err = fmt.Sscanf(verstr, "%d.%d.%d", &ver.Major, &ver.Minor, &ver.Patch)
		ver.Extension = ""
	}

	return ver, err
}

func EmptyVersion() Version {
	return Version{
		Major:     0,
		Minor:     0,
		Patch:     0,
		Extension: "",
	}
}

func Ver(major, minor, patch int) Version {
	return Version{
		Major:     major,
		Minor:     minor,
		Patch:     patch,
		Extension: "",
	}
}
