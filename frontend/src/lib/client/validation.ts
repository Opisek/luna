const urlRegex = /https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,4}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)/;

const validResponse = {
  valid: true,
  message: ""
}

const invalidResponse = (message: string): Validity => {
  return {
    valid: false,
    message: message
  }
}

export const isValidUsername: InputValidation = (username) => {
  if (username.length < 3)
    return invalidResponse("Username must be at least 3 characters long.");
  if (username.length > 25)
    return invalidResponse("Username must be at most 20 characters long.");
  if (!/^[a-zA-Z0-9]*$/.test(username))
    return invalidResponse("Username must only contain letters and numbers.");
  return validResponse;
}

export const isValidPassword: InputValidation = (password) => {
  if (password.length < 8)
    return invalidResponse("Password must be at least 8 characters long.");
  if (password.length > 50)
    return invalidResponse("Password must be at most 50 characters long.");
  return validResponse;
}

export const isValidRepeatPassword = (repeatPassword: string): InputValidation => {
  return (password) => {
    if (password !== repeatPassword)
      return invalidResponse("Passwords do not match.");
    return validResponse;
  }
}

export const isValidUrl: InputValidation = (url) => {
  if (!url.startsWith("http://") && !url.startsWith("https://"))
    return invalidResponse("URL must start with \"http://\" or \"https://\".");
  if (!urlRegex.test(url))
    return invalidResponse("The URL contains illegal characters or is invalid.");
  return validResponse;
}