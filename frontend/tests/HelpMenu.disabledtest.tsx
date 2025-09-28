// HelpMenu.test.tsx
import React from "react";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { BrowserRouter } from "react-router-dom";
import HelpMenu from "../src/components/HelpMenu";
import "@testing-library/jest-dom";

// Mock fetch for chatbot API
global.fetch = jest.fn(() =>
  Promise.resolve({
    ok: true,
    json: () =>
      Promise.resolve({
        fulfillmentText: "Hello from bot",
        fulfillmentMessages: [],
      }),
  }),
) as jest.Mock;

const renderWithRouter = (ui: React.ReactElement) => {
  return render(<BrowserRouter>{ui}</BrowserRouter>);
};

describe("HelpMenu Component", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  test("renders help button and toggles menu open/close", () => {
    renderWithRouter(<HelpMenu />);
    const helpButton = screen.getByText(/HELP/i);
    expect(helpButton).toBeInTheDocument();

    // Open menu
    fireEvent.click(helpButton);
    expect(screen.getByText(/Swift Chat/i)).toBeInTheDocument();

    // Close menu
    const closeButton = screen.getAllByRole("button", { name: "" })[0];
    fireEvent.click(closeButton);
    expect(screen.queryByText(/Swift Chat/i)).not.toBeInTheDocument();
  });

  test("switches between chat and general help tabs", () => {
    renderWithRouter(<HelpMenu />);
    fireEvent.click(screen.getByText(/HELP/i));

    // Click General Help tab
    const generalTab = screen.getByText(/General Help/i);
    fireEvent.click(generalTab);
    expect(screen.getByText(/Tutorials/i)).toBeInTheDocument();

    // Click Chat tab
    const chatTab = screen.getByText(/Swift Chat/i);
    fireEvent.click(chatTab);
    expect(
      screen.getByPlaceholderText(/Type your message/i),
    ).toBeInTheDocument();
  });

  test("expands and collapses FAQ items", () => {
    renderWithRouter(<HelpMenu />);
    fireEvent.click(screen.getByText(/HELP/i));
    fireEvent.click(screen.getByText(/General Help/i));

    // Open FAQ section
    fireEvent.click(screen.getByText(/Frequently Asked Questions/i));
    const firstQuestion = screen.getByText(/What is Swift Signals/i);
    fireEvent.click(firstQuestion);

    expect(
      screen.getByText(/Swift Signals is a simulation-powered/i),
    ).toBeInTheDocument();

    // Collapse FAQ
    fireEvent.click(firstQuestion);
    expect(
      screen.queryByText(/Swift Signals is a simulation-powered/i),
    ).not.toBeVisible();
  });

  test("starts a tutorial from General Help", () => {
    renderWithRouter(<HelpMenu />);
    fireEvent.click(screen.getByText(/HELP/i));
    fireEvent.click(screen.getByText(/General Help/i));

    const dashboardTutorialButton = screen.getByText(/Dashboard Tutorial/i);
    fireEvent.click(dashboardTutorialButton);

    // InteractiveTutorial for Dashboard should render
    expect(screen.getByText(/Summary Cards/i)).toBeInTheDocument();
  });

  test("sends a chat message and displays bot response", async () => {
    renderWithRouter(<HelpMenu />);
    fireEvent.click(screen.getByText(/HELP/i));

    const input = screen.getByPlaceholderText(/Type your message/i);
    fireEvent.change(input, { target: { value: "Hello bot" } });
    fireEvent.keyPress(input, { key: "Enter", code: "Enter", charCode: 13 });

    await waitFor(() => {
      expect(screen.getByText(/Hello from bot/i)).toBeInTheDocument();
    });
  });

  test("quick replies render when bot provides options", async () => {
    (global.fetch as jest.Mock).mockImplementationOnce(() =>
      Promise.resolve({
        ok: true,
        json: () =>
          Promise.resolve({
            fulfillmentText: "Choose an option",
            fulfillmentMessages: [
              {
                payload: {
                  fields: {
                    richContent: {
                      listValue: {
                        values: [
                          {
                            listValue: {
                              values: [
                                {
                                  structValue: {
                                    fields: {
                                      options: {
                                        listValue: {
                                          values: [
                                            {
                                              structValue: {
                                                fields: {
                                                  text: {
                                                    stringValue: "Option 1",
                                                  },
                                                },
                                              },
                                            },
                                          ],
                                        },
                                      },
                                    },
                                  },
                                },
                              ],
                            },
                          },
                        ],
                      },
                    },
                  },
                },
              },
            ],
          }),
      }),
    );

    renderWithRouter(<HelpMenu />);
    fireEvent.click(screen.getByText(/HELP/i));

    const input = screen.getByPlaceholderText(/Type your message/i);
    fireEvent.change(input, { target: { value: "Show options" } });
    fireEvent.keyPress(input, { key: "Enter", code: "Enter", charCode: 13 });

    await waitFor(() => {
      expect(screen.getByText(/Option 1/i)).toBeInTheDocument();
    });
  });
});
