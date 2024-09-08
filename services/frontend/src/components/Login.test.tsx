import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import Login from "./Login";

describe("Login", () => {
  it("renders the login heading", () => {
    render(<Login />);
    expect(screen.getByRole("heading", { name: /login/i })).toBeInTheDocument();
  });

  it("renders the login with Google link", () => {
    render(<Login />);
    const loginLink = screen.getByRole("link", { name: /login with google/i });
    expect(loginLink).toBeInTheDocument();
    expect(loginLink).toHaveAttribute(
      "href",
      "http://localhost:8080/auth/google",
    );
  });

  it("has the correct structure", () => {
    const { container } = render(<Login />);

    // Check for the heading
    const heading = screen.getByRole("heading", { name: /login/i, level: 1 });
    expect(heading).toBeInTheDocument();

    // Check for the link
    const loginLink = screen.getByRole("link", { name: /login with google/i });
    expect(loginLink).toBeInTheDocument();
    expect(loginLink).toHaveAttribute(
      "href",
      "http://localhost:8080/auth/google",
    );

    // Check the order of elements
    const mainDiv = container.firstChild as HTMLElement;
    expect(mainDiv.children).toHaveLength(2);
    expect(mainDiv.children[0]).toBe(heading);
    expect(mainDiv.children[1]).toBe(loginLink);
  });
});
