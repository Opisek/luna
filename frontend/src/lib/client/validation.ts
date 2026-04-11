import { t } from "@sveltia/i18n";

const characterRegex = /^[a-zA-Z0-9]*$/;
// RFC 1738 3.3, RFC 1034 3.5, RFC 3986 3 
// IPvFuture is not included
// Ignoring this detail, the following is a strict supset of possible HTTP(S) URLs:
// - Some invalid IPv6 addresses are accepted
// - The port can exceed the 16 bit limit
// - Some invalid hostnames can be accepted
// - Some invalid queries can be accepted
// - Length limits specific to HTTP(S) are not considered
// - Valid TLDs are not considered
const urlRegex = /https?:\/\/((\d+\.){3}\d+|\[([a-fA-F0-9]{0,4}:){2,7}[a-zA-Z0-9]{0,4}\]|(([a-zA-Z0-9-\._~!\$&'\(\)\*\+,;=]|\%[a-fA-F]{2})\.?)+)?(:\d{1,5})?(\/([a-zA-Z0-9-\._~!\$&'\(\)\*\+,;=:@]|\%[a-fA-F]{2})+)*\/?(\?([a-zA-Z0-9-\._~!\$&'\(\)\*\+,;=:@\/\?]|\%[a-fA-F]{2})*)?(\#([a-zA-Z0-9-\._~!\$&'\(\)\*\+,;=:@\/\?]|\%[a-fA-F]{2})*)?/;
const emailRegex = /(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])/
const inviteCodeRegex = /^[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{4}$/

export const valid = {
  valid: true,
  message: ""
}

const invalidResponse = (message: string): Validity => {
  return {
    valid: false,
    message: message
  }
}

export const alwaysValid: InputValidation = () => Promise.resolve(valid);
export const alwaysValidFile: FileValidation = () => Promise.resolve(valid);

export const isValidUsername: InputValidation = async (username) => {
  if (username.length < 3)
    return invalidResponse(t("validation.username.short", { values: { len: 3 } }));
  if (username.length > 25)
    return invalidResponse(t("validation.username.long", { values: { len: 25 } }));
  if (!characterRegex.test(username))
    return invalidResponse(t("validation.username.characters"));
  return valid;
}

export const isValidPassword: InputValidation = async (password) => {
  if (password.length < 8)
    return invalidResponse(t("validation.password.short", { values: { len: 8 } }));
  if (password.length > 1000)
    return invalidResponse(t("validation.password.long", { values: { len: 1000 } }));
  return valid;
}

export const isValidRepeatPassword = (repeatPassword: string): InputValidation => {
  return async (password) => {
    if (password !== repeatPassword)
      return invalidResponse(t("validation.password.match"));
    return valid;
  }
}

export const isValidUrl: InputValidation = async (url) => {
  if (!url.startsWith("http://") && !url.startsWith("https://"))
    return invalidResponse(t("validation.url.proto"));
  if (!urlRegex.test(url))
    return invalidResponse(t("validation.url.invalid"));
  return valid;
}

export const isValidEmail: InputValidation = async (email) => {
  if (!emailRegex.test(email))
    return invalidResponse(t("validation.email.invalid"));
  return valid;
}

export const isValidPath: InputValidation = async (path) => {
  if (path.length < 1)
    return invalidResponse(t("validation.path.empty"));
  return valid;
}

export const isValidInviteCode: InputValidation = async (inviteCode) => {
  if (!inviteCodeRegex.test(inviteCode))
    return invalidResponse(t("validation.invite.invalid"));
  return valid;
}

export const isValidFile: FileValidation = async (files) => {
  if (files.length === 0)
    return invalidResponse(t("validation.file.empty"));
  if (files.length > 1)
    return invalidResponse(t("validation.file.multiple"));
  if (files.item(0) === null)
    return invalidResponse(t("validation.file.null"));
  if ((files.item(0) as File).size > 50000000)
    return invalidResponse(t("validation.file.size", { values: { size: "50MB" } }));
  return valid;
}

export const isValidIcalFile: FileValidation = async (files) => {
  const fileValidation = await isValidFile(files);
  if (!fileValidation.valid)
    return fileValidation;

  const file = files.item(0) as File;

  if (file.type !== "text/calendar")
    return invalidResponse(t("validation.file.type", { values: { type: "text/calendar." } }));
  
  const content = await file.text();

  if (!content.includes("BEGIN:VCALENDAR") || !content.includes("END:VCALENDAR"))
    return invalidResponse(t("validation.file.invalid"));

  return valid;
}
