// tests/HelpMenu.test.tsx
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import HelpMenu from "../src/components/HelpMenu";
import { MemoryRouter } from "react-router-dom";

// Mock uuid
jest.mock("uuid", () => ({
  v4: () => "test-uuid",
}));

// Mock fetch
(globalThis as any).fetch = jest.fn(() =>
  Promise.resolve({
    ok: true,
    json: () =>
      Promise.resolve({
        fulfillmentText: "Hello from bot!",
        fulfillmentMessages: [],
      }),
  })
);

// Mock InteractiveTutorial
jest.mock("../src/components/InteractiveTutorial", () => ({
  __esModule: true,
  default: (props: any) => {
    return <div data-testid={`tutorial-${props.tutorialType}`}>Tutorial Open</div>;
  },
}));

// Suppress specific React warnings
beforeAll(() => {
  const originalConsoleError = console.error;
  console.error = (...args: unknown[]) => {
    if (
      typeof args[0] === "string" &&
      args[0].includes("React does not recognize the `dragConstraints` prop")
    ) {
      return;
    }
    originalConsoleError(...args);
  };
});

describe("HelpMenu Component", () => {
  it("renders the help button", () => {
    render(
      <MemoryRouter>
        <HelpMenu />
      </MemoryRouter>
    );
    expect(screen.getByText(/HELP/i)).toBeInTheDocument();
  });

  it("opens and closes the help menu on button click", () => {
    render(
      <MemoryRouter>
        <HelpMenu />
      </MemoryRouter>
    );
    const button = screen.getByRole("button", { name: /HELP/i });
    fireEvent.click(button);
    expect(screen.getByText(/Swift Chat/i)).toBeInTheDocument();
    fireEvent.click(button);
    expect(screen.queryByText(/Swift Chat/i)).not.toBeVisible();
  });

  it("switches between chat and general help tabs", () => {
    render(
      <MemoryRouter>
        <HelpMenu />
      </MemoryRouter>
    );
    fireEvent.click(screen.getByRole("button", { name: /HELP/i }));
    const generalTab = screen.getByText(/General Help/i);
    fireEvent.click(generalTab);
    expect(screen.getByText(/Tutorials/i)).toBeInTheDocument();
    const chatTab = screen.getByText(/Swift Chat/i);
    fireEvent.click(chatTab);
    expect(screen.getByPlaceholderText(/Type your message/i)).toBeInTheDocument();
  });

  it("can send a chat message", async () => {
    render(
      <MemoryRouter>
        <HelpMenu />
      </MemoryRouter>
    );
    fireEvent.click(screen.getByRole("button", { name: /HELP/i }));
    const input = screen.getByPlaceholderText(/Type your message/i);
    fireEvent.change(input, { target: { value: "Hello" } });
    fireEvent.keyPress(input, { key: "Enter", code: "Enter", charCode: 13 });

    await waitFor(() =>
      expect(screen.getByText("Hello from bot!")).toBeInTheDocument()
    );
  });

  it("toggles FAQ sections", () => {
    render(
      <MemoryRouter>
        <HelpMenu />
      </MemoryRouter>
    );
    fireEvent.click(screen.getByRole("button", { name: /HELP/i }));
    fireEvent.click(screen.getByText(/General Help/i));
    const firstFaq = screen.getByText("What is Swift Signals?");
    fireEvent.click(firstFaq);
    expect(
      screen.getByText(/Swift Signals is a simulation-powered/i)
    ).toBeVisible();
    fireEvent.click(firstFaq);
    expect(
      screen.getByText(/Swift Signals is a simulation-powered/i)
    ).not.toHaveClass("open");
  });

  it("launches tutorials from accordion", () => {
    render(
      <MemoryRouter>
        <HelpMenu />
      </MemoryRouter>
    );
    fireEvent.click(screen.getByRole("button", { name: /HELP/i }));
    fireEvent.click(screen.getByText(/General Help/i));
    fireEvent.click(screen.getByText(/Navigation Tutorial/i));
    expect(screen.getByTestId("tutorial-navigation")).toBeInTheDocument();
  });

  it("opens confirmation overlay when tutorial page differs", () => {
    render(
      <MemoryRouter initialEntries={["/somepath"]}>
        <HelpMenu />
      </MemoryRouter>
    );
    fireEvent.click(screen.getByRole("button", { name: /HELP/i }));
    fireEvent.click(screen.getByText(/Dashboard Tutorial/i));
    expect(screen.getByText(/Switch to Dashboard/i)).toBeInTheDocument();
  });
});
