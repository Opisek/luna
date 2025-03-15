export const FRONTEND_VERSION = "0.1.0";

export enum VersionCompatibility {
  Unknown,
  Compatible,
  FrontendOutdatedMinor,
  FrontendOutdatedMajor,
  BackendOutdatedMinor,
  BackendOutdatedMajor
} 

export function isCompatibleWithBackend(backendVersion: string) {
  if (backendVersion === null || backendVersion === "" || backendVersion === undefined) {
    return VersionCompatibility.Unknown;
  }

  const frontendSplit = FRONTEND_VERSION.split(".");
  const backendSplit = backendVersion.split(".");

  if (frontendSplit.length !== 3 || backendSplit.length !== 3) {
    return VersionCompatibility.Unknown;
  }

  const frontendMajor = parseInt(frontendSplit[0]);
  const backendMajor = parseInt(backendSplit[0]);
  const frontendMinor = parseInt(frontendSplit[1]);
  const backendMinor = parseInt(backendSplit[1]);
  const frontendPatch = parseInt(frontendSplit[2]);
  const backendPatch = parseInt(backendSplit[2]);

  if (
    isNaN(frontendMajor) || frontendMajor < 0 ||
    isNaN(backendMajor) || backendMajor < 0 ||
    isNaN(frontendMinor) || frontendMinor < 0 ||
    isNaN(backendMinor) || backendMinor < 0 ||
    isNaN(frontendPatch) || frontendPatch < 0 ||
    isNaN(backendPatch) || backendPatch
  ) {
    return VersionCompatibility.Unknown;
  }

  if (frontendMajor > backendMajor) {
    return VersionCompatibility.FrontendOutdatedMajor;
  } else if (frontendMajor < backendMajor) {
    return VersionCompatibility.BackendOutdatedMajor;
  }

  if (frontendMinor > backendMinor) {
    return VersionCompatibility.FrontendOutdatedMinor;
  } else if (frontendMinor < backendMinor) {
    return VersionCompatibility.BackendOutdatedMinor;
  }

  return VersionCompatibility.Compatible;
}