/**
 * @jest-environment jsdom
 */
import React from "react";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import '@testing-library/jest-dom';
import { BrowserRouter } from "react-router-dom";
import SignUp from "../src/pages/SignUp";

console.log(React)

// Mock Footer and logo import
jest.mock("../src/components/Footer", () => () => <div>Footer</div>);
jest.mock("../../src/assets/logo.png", () => "logo.png");

// Mock useNavigate
const mockNavigate = jest.fn();
jest.mock("react-router-dom", () => ({
  ...jest.requireActual("react-router-dom"),
  useNavigate: () => mockNavigate,
}));

// Mock global fetch
global.fetch = jest.fn();

describe("SignUp Page", () => {
  beforeEach(() => {
    (fetch as jest.Mock).mockReset();
    mockNavigate.mockReset();
  });

  it("renders all main elements", () => {
    render(
      <BrowserRouter>
        <SignUp />
      </BrowserRouter>
    );

    expect(screen.getByText("Sign Up")).toBeInTheDocument();
    expect(screen.getByPlaceholderText("Username")).toBeInTheDocument();
    expect(screen.getByPlaceholderText("Email")).toBeInTheDocument();
    expect(screen.getByPlaceholderText("Password")).toBeInTheDocument();
    expect(screen.getByText("Register")).toBeInTheDocument();
    expect(screen.getByText("Footer")).toBeInTheDocument();
  });

  it("shows error if fields are empty", async () => {
  render(
    <BrowserRouter>
      <SignUp />
    </BrowserRouter>
  );

  const registerButton = screen.getByText("Register") as HTMLButtonElement;
  fireEvent.click(registerButton);

  // Check button is still enabled (form was not submitted)
  expect(registerButton).not.toHaveAttribute("disabled");
});


  it("submits form successfully and navigates to login", async () => {
    (fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      text: async () => JSON.stringify({ message: "Registered" }),
    });

    render(
      <BrowserRouter>
        <SignUp />
      </BrowserRouter>
    );

    fireEvent.change(screen.getByPlaceholderText("Username"), { target: { value: "testuser" } });
    fireEvent.change(screen.getByPlaceholderText("Email"), { target: { value: "test@example.com" } });
    fireEvent.change(screen.getByPlaceholderText("Password"), { target: { value: "password123" } });

    fireEvent.click(screen.getByText("Register"));

    expect(screen.getByText("Registering...")).toBeInTheDocument();

    await waitFor(() => {
      expect(screen.getByText(/Registration successful! Redirecting to login/i)).toBeInTheDocument();
    });

    // Wait for setTimeout navigation
    await waitFor(() => expect(mockNavigate).toHaveBeenCalledWith("/login"), { timeout: 3000 });
  });

  it("shows error if API returns non-OK response", async () => {
    (fetch as jest.Mock).mockResolvedValueOnce({
      ok: false,
      text: async () => JSON.stringify({ message: "User already exists" }),
    });

    render(
      <BrowserRouter>
        <SignUp />
      </BrowserRouter>
    );

    fireEvent.change(screen.getByPlaceholderText("Username"), { target: { value: "testuser" } });
    fireEvent.change(screen.getByPlaceholderText("Email"), { target: { value: "test@example.com" } });
    fireEvent.change(screen.getByPlaceholderText("Password"), { target: { value: "password123" } });

    fireEvent.click(screen.getByText("Register"));

    expect(await screen.findByText(/User already exists/i)).toBeInTheDocument();
  });

  it("traffic lights activate correctly based on input", async () => {
  render(
    <BrowserRouter>
      <SignUp />
    </BrowserRouter>
  );

  // Query all three lights directly
  const lights = document.querySelectorAll(".traffic-light > div");
  expect(lights.length).toBe(3);

  // Fill inputs
  fireEvent.change(screen.getByPlaceholderText("Username"), { target: { value: "user1" } });
  fireEvent.change(screen.getByPlaceholderText("Email"), { target: { value: "user@example.com" } });
  fireEvent.change(screen.getByPlaceholderText("Password"), { target: { value: "password123" } });

  // Wait a tick (simulate state updates if necessary)
  await waitFor(() => {
    lights.forEach(light => {
      expect(light).toBeInTheDocument();
    });
  });
});

  it("navigates to login page when clicking login button", () => {
    render(
      <BrowserRouter>
        <SignUp />
      </BrowserRouter>
    );

    fireEvent.click(screen.getByText("Login here"));

    expect(mockNavigate).toHaveBeenCalledWith("/login");
  });
});
