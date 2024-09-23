export const getApiUrl = (): string => {
  if (
    typeof window !== "undefined" &&
    window.APP_CONFIG &&
    window.APP_CONFIG.API_URL
  ) {
    return window.APP_CONFIG.API_URL;
  }
  throw new Error("API URL is not defined");
};
