import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter, useNavigate, useLocation } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Callback, Logout } from "./Callback";
import * as apiModule from "../api/apiClient";
import * as useApiErrorModule from "../hooks/useApiError";

// Mock react-router-dom
vi.mock("react-router-dom", async () => {
  const actual = await vi.importActual("react-router-dom");
  return {
    ...actual,
    useNavigate: vi.fn(),
    useLocation: vi.fn(),
  };
});

// Mock the entire api module
vi.mock("../api/apiClient", () => ({
  api: {
    url: vi.fn().mockReturnThis(),
    get: vi.fn().mockReturnThis(),
    post: vi.fn().mockReturnThis(),
    json: vi.fn(),
  },
}));

// Mock useApiError hook
vi.mock("../hooks/useApiError", () => ({
  useApiError: vi.fn(),
}));

let queryClient: QueryClient;

beforeEach(() => {
  queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });
});

afterEach(() => {
  queryClient.clear();
});

const renderWithProviders = (component: React.ReactElement) => {
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>{component}</MemoryRouter>
    </QueryClientProvider>,
  );
};

describe("Callback", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.resetAllMocks();
    (useLocation as any).mockReturnValue({ search: "?code=test" });
    (useNavigate as any).mockReturnValue(vi.fn());
    (useApiErrorModule.useApiError as any).mockReturnValue("");
  });

  it("handles successful login", async () => {
    (apiModule.api.url as any).mockReturnValue({
      get: vi.fn().mockReturnThis(),
      json: vi.fn().mockResolvedValue({}),
    });
    const navigate = vi.fn();
    (useNavigate as any).mockReturnValue(navigate);

    renderWithProviders(<Callback />);

    await waitFor(() => {
      expect(screen.getByText(/login successful/i)).toBeInTheDocument();
    });
    expect(navigate).toHaveBeenCalledWith("/user-info");
  });

  it("handles login error", async () => {
    (apiModule.api.url as any).mockReturnValue({
      get: vi.fn().mockReturnThis(),
      json: vi.fn().mockRejectedValue(new Error("Login failed")),
    });
    (useApiErrorModule.useApiError as any).mockReturnValue("Login failed");

    renderWithProviders(<Callback />);

    await waitFor(() => {
      expect(screen.getByText(/Login failed/i)).toBeInTheDocument();
    });
  });

  it("handles offline scenario", async () => {
    (apiModule.api.url as any).mockReturnValue({
      get: vi.fn().mockReturnThis(),
      json: vi.fn().mockRejectedValue(new Error("Offline: Network error")),
    });
    (useApiErrorModule.useApiError as any).mockReturnValue(
      "You are currently offline. Some features may be unavailable.",
    );

    renderWithProviders(<Callback />);

    await waitFor(() => {
      expect(
        screen.getByText(/You are currently offline/i),
      ).toBeInTheDocument();
    });
  });
});

describe("Logout", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.resetAllMocks();
    (useNavigate as any).mockReturnValue(vi.fn());
    (useApiErrorModule.useApiError as any).mockReturnValue("");
  });

  it("calls logout API and navigates to home", async () => {
    (apiModule.api.url as any).mockReturnValue({
      get: vi.fn().mockReturnValue({ json: vi.fn().mockResolvedValue({}) }),
    });
    const navigate = vi.fn();
    (useNavigate as any).mockReturnValue(navigate);

    renderWithProviders(<Logout />);

    await waitFor(() => {
      expect(apiModule.api.url).toHaveBeenCalledWith("/auth/logout");
      expect(navigate).toHaveBeenCalledWith("/");
    });
  });

  it("handles logout error", async () => {
    (apiModule.api.url as any).mockReturnValue({
      get: vi.fn().mockReturnValue({
        json: vi.fn().mockRejectedValue(new Error("Logout failed")),
      }),
    });
    (useApiErrorModule.useApiError as any).mockReturnValue("Logout failed");

    renderWithProviders(<Logout />);

    await waitFor(() => {
      expect(screen.getByText(/Logout failed/i)).toBeInTheDocument();
    });
  });

  it("handles auth error during logout", async () => {
    (apiModule.api.url as any).mockReturnValue({
      get: vi.fn().mockReturnValue({
        json: vi.fn().mockRejectedValue(new Error("AuthError")),
      }),
    });
    (useApiErrorModule.useApiError as any).mockImplementation(() => {
      const navigate = useNavigate();
      navigate("/login");
      return "Auth Error";
    });

    const navigate = vi.fn();
    (useNavigate as any).mockReturnValue(navigate);

    renderWithProviders(<Logout />);

    await waitFor(() => {
      expect(navigate).toHaveBeenCalledWith("/login");
    });
  });
});
