export async function callApi(endpoint: string, init?: RequestInit): Promise<Response> {
  // @ts-ignore
  return await fetch(`${process.env.API_URL}/api/${endpoint}`, init).catch((error) => {
    throw error;
  });
}