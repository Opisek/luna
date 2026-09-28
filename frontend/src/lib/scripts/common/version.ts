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

  let frontendMajor = parseInt(frontendSplit[0]);
  let backendMajor = parseInt(backendSplit[0]);
  let frontendMinor = parseInt(frontendSplit[1]);
  let backendMinor = parseInt(backendSplit[1]);
  let frontendPatch = parseInt(frontendSplit[2]);
  let backendPatch = parseInt(backendSplit[2]);

  if (frontendMajor === 0 && backendMajor === 0) {
    frontendMajor = frontendMinor;
    backendMajor = backendMinor;
    frontendMinor = frontendPatch;
    backendMinor = backendPatch;
    frontendPatch = 0;
    backendPatch = 0;
  }

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
    return VersionCompatibility.BackendOutdatedMajor;
  } else if (frontendMajor < backendMajor) {
    return VersionCompatibility.FrontendOutdatedMajor;
  }

  if (frontendMinor > backendMinor) {
    return VersionCompatibility.BackendOutdatedMinor;
  } else if (frontendMinor < backendMinor) {
    return VersionCompatibility.BackendOutdatedMajor;
  }

  return VersionCompatibility.Compatible;
}