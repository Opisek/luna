export function callApi(endpoint: string, init?: RequestInit): Promise<Response> {
  // @ts-ignore
  return fetch(`${process.env.API_URL}/api/${endpoint}`, init)
}