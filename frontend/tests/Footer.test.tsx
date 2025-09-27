import { render, screen } from "@testing-library/react";
import Footer from "../src/components/Footer";

jest.mock("../../src/assets/scs-logo.png", () => "mocked-logo.png", { virtual: true });

jest.mock("../src/components/ThemeToggle", () => () => (
  <div data-testid="theme-toggle" />
));

describe("Footer", () => {
  it("renders the footer container", () => {
    render(<Footer />);
    expect(screen.getByRole("contentinfo")).toBeInTheDocument();
  });

  it("renders the footer text", () => {
    render(<Footer />);
    expect(
      screen.getByText("A Southern Cross Solutions Product")
    ).toBeInTheDocument();
  });

  it("renders the current year dynamically", () => {
    render(<Footer />);
    const year = new Date().getFullYear();
    expect(
      screen.getByText(`© ${year} Swift Signals. All rights reserved.`)
    ).toBeInTheDocument();
  });

  it("renders the ThemeToggle component", () => {
    render(<Footer />);
    expect(screen.getByTestId("theme-toggle")).toBeInTheDocument();
  });
});
