// tests/HelpMenu.test.tsx
import React from "react";
import { render, screen, fireEvent } from "@testing-library/react";
import "@testing-library/jest-dom"; // adds .toBeInTheDocument() and other matchers
import { MemoryRouter } from "react-router-dom";
import HelpMenu from "../src/components/HelpMenu";

console.log(React)

const renderHelpMenu = () =>
  render(
    <MemoryRouter>
      <HelpMenu />
    </MemoryRouter>
  );

describe("HelpMenu component", () => {
  test("opens and closes the help menu", () => {
    renderHelpMenu();

    const helpButton = screen.getByRole("button", { name: /help/i });
    fireEvent.click(helpButton);

    expect(screen.getByText(/Swift Chat/i)).toBeInTheDocument();

    const closeButton = screen.getAllByRole("button", { name: /times/i })[0];
    fireEvent.click(closeButton);

    expect(screen.queryByText(/Swift Chat/i)).not.toBeInTheDocument();
  });

  test("switches tabs between chat and general help", () => {
    renderHelpMenu();
    fireEvent.click(screen.getByRole("button", { name: /help/i }));

    const generalTab = screen.getByRole("button", { name: /general help/i });
    fireEvent.click(generalTab);
    expect(screen.getByText(/Tutorials/i)).toBeInTheDocument();

    const chatTab = screen.getByRole("button", { name: /swift chat/i });
    fireEvent.click(chatTab);
    expect(screen.getByPlaceholderText(/Type your message/i)).toBeInTheDocument();
  });

  test("opens the dashboard tutorial", () => {
    renderHelpMenu();
    fireEvent.click(screen.getByRole("button", { name: /help/i }));

    fireEvent.click(screen.getByText(/Tutorials/i));

    const dashboardButton = screen.getByRole("button", { name: /dashboard tutorial/i });
    fireEvent.click(dashboardButton);

    expect(screen.getByText(/Summary Cards/i)).toBeInTheDocument();
  });

  test("opens the intersections tutorial", () => {
    renderHelpMenu();
    fireEvent.click(screen.getByRole("button", { name: /help/i }));

    fireEvent.click(screen.getByText(/Tutorials/i));

    const intersectionsButton = screen.getByRole("button", { name: /intersections tutorial/i });
    fireEvent.click(intersectionsButton);

    expect(screen.getByText(/Search Bar/i)).toBeInTheDocument();
  });

  test("toggles FAQ sections", () => {
    renderHelpMenu();

    // Use querySelector to target the top-level HELP button specifically
    const helpButton = screen.getByText("HELP").closest("button");
    if (!helpButton) throw new Error("Top-level HELP button not found");
    fireEvent.click(helpButton);

    // Open FAQ accordion
    const faqAccordion = screen.getByText(/Frequently Asked Questions/i);
    fireEvent.click(faqAccordion);

    // Open a specific FAQ
    const faqQuestion = screen.getByRole("button", { name: /What is Swift Signals/i });
    fireEvent.click(faqQuestion);

    expect(screen.getByText(/Swift Signals is a simulation-powered/i)).toBeInTheDocument();
  });
});
