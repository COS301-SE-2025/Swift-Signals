import React from "react";
import { render, screen, fireEvent } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import WelcomePage from "../src/pages/WelcomePage";

console.log(React)

const mockNavigate = jest.fn();
jest.mock("react-router-dom", () => ({
  ...(jest.requireActual("react-router-dom") as any),
  useNavigate: () => mockNavigate,
}));

jest.mock("../src/components/Carousel", () => (_props: any) => (
  <div data-testid="carousel">Carousel Component</div>
));

describe("WelcomePage", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("renders the welcome page correctly", () => {
    render(
      <MemoryRouter>
        <WelcomePage />
      </MemoryRouter>
    );

    const logo = screen.getByAltText("Logo");
    expect(logo).toBeInTheDocument();

    const heading = screen.getByText(/Welcome to Swift Signals!/i);
    expect(heading).toBeInTheDocument();

    const loginButton = screen.getByText("Login");
    const registerButton = screen.getByText("Register");
    expect(loginButton).toBeInTheDocument();
    expect(registerButton).toBeInTheDocument();

    const carousel = screen.getByTestId("carousel");
    expect(carousel).toBeInTheDocument();
  });

  it("navigates to /login when Login button is clicked", () => {
    render(
      <MemoryRouter>
        <WelcomePage />
      </MemoryRouter>
    );

    const loginButton = screen.getByText("Login");
    fireEvent.click(loginButton);

    expect(mockNavigate).toHaveBeenCalledWith("/login");
  });

  it("navigates to /signup when Register button is clicked", () => {
    render(
      <MemoryRouter>
        <WelcomePage />
      </MemoryRouter>
    );

    const registerButton = screen.getByText("Register");
    fireEvent.click(registerButton);

    expect(mockNavigate).toHaveBeenCalledWith("/signup");
  });
});
