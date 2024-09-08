import wretch, { Wretch } from "wretch";

const BASE_URL = import.meta.env.VITE_API_BASE_URL;

type Middleware = (
  next: (url: string, opts: RequestInit) => Promise<Response>,
) => (url: string, opts: RequestInit) => Promise<Response>;

class Mutex {
  private locked = false;
  private waitQueue: Array<() => void> = [];

  async acquire(): Promise<void> {
    if (!this.locked) {
      this.locked = true;
      return;
    }
    return new Promise<void>((resolve) => {
      this.waitQueue.push(resolve);
    });
  }

  release(): void {
    if (this.waitQueue.length > 0) {
      const next = this.waitQueue.shift();
      if (next) next();
    } else {
      this.locked = false;
    }
  }
}

let refreshPromise: Promise<void> | null = null;
const refreshMutex = new Mutex();

const authMiddleware: Middleware = (next) => async (url, opts) => {
  const performRequest = async () => {
    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), 10000); // 10 seconds timeout

      const response = await next(url, { ...opts, signal: controller.signal });
      clearTimeout(timeoutId);

      if (response.status === 401) {
        throw new Error("Unauthorized");
      }
      return response;
    } catch (error) {
      if (error instanceof Error) {
        if (
          error.name === "TypeError" ||
          error.message.includes("NetworkError") ||
          error.message.includes("Failed to fetch") ||
          error.name === "AbortError" // Handle timeout
        ) {
          return new Response(
            JSON.stringify({
              offline: true,
              debug: error.toString(),
              timestamp: new Date().toISOString(),
            }),
            {
              status: 200,
              headers: {
                "Content-Type": "application/json",
                "X-Offline": "true",
              },
            },
          );
        }
      }

      if (error instanceof Response && error.status >= 500) {
        return new Response(
          JSON.stringify({ serverError: true, status: error.status }),
          {
            status: 200,
            headers: {
              "Content-Type": "application/json",
              "X-Server-Error": "true",
            },
          },
        );
      }

      throw error;
    }
  };

  const refreshToken = async () => {
    await refreshMutex.acquire();
    try {
      if (!refreshPromise) {
        refreshPromise = wretch(BASE_URL)
          .url("/auth/refresh")
          .options({ credentials: "include" })
          .post()
          .res()
          .then(() => {
            refreshPromise = null;
          });
      }
      await refreshPromise;
    } finally {
      refreshMutex.release();
    }
  };

  try {
    return await performRequest();
  } catch (error) {
    if (
      error instanceof Error &&
      (error.message === "Unauthorized" || error.name === "AbortError")
    ) {
      try {
        await refreshToken();
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

export const isOfflineResponse = (response: Response): boolean =>
  response.headers.get("X-Offline") === "true";

export const isServerErrorResponse = (response: Response): boolean =>
  response.headers.get("X-Server-Error") === "true";
