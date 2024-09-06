import wretch, { Wretch } from "wretch";

const BASE_URL = "http://localhost:8080";

type Middleware = (
  next: (url: string, opts: RequestInit) => Promise<Response>,
) => (url: string, opts: RequestInit) => Promise<Response>;

const authMiddleware: Middleware = (next) => async (url, opts) => {
  const performRequest = async () => {
    try {
      const response = await next(url, opts);
      if (response.status === 401) {
        throw new Error("Unauthorized");
      }
      return response;
    } catch (error) {
      // Ignore connection errors and 5xx errors
      if (
        error instanceof Error &&
        (error.message.includes("Failed to fetch") ||
          (error instanceof Response && error.status >= 500))
      ) {
        console.warn("Network or server error occurred, ignoring:", error);
        return new Response(null, { status: 200 }); // Return a fake 200 response
      }
      throw error;
    }
  };

  try {
    return await performRequest();
  } catch (error) {
    if (error instanceof Error && error.message === "Unauthorized") {
      try {
        await wretch(BASE_URL)
          .url("/auth/refresh")
          .options({ credentials: "include" })
          .post()
          .res();
        return await performRequest();
      } catch (refreshError) {
        throw new Error("AuthError");
      }
    }
    throw error;
  }
};

export const api: Wretch = wretch(BASE_URL)
  .options({ credentials: "include" })
  .middlewares([authMiddleware]);

export const isAuthError = (error: unknown): boolean =>
  error instanceof Error && error.message === "AuthError";
