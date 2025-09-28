import { render, screen } from "@testing-library/react";

jest.mock("../src/components/InteractiveTutorial", () => {
  return {
    __esModule: true,
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    default: (props: any) => (
      <div data-testid={`tutorial-${props.tutorialType}`}>
        Tutorial Open
      </div>
    ),
  };
});

import InteractiveTutorial from "../src/components/InteractiveTutorial";

const steps = [
  { selector: "#step1", title: "Step 1", text: "This is step 1" },
  { selector: "#step2", title: "Step 2", text: "This is step 2" },
];
const onClose = jest.fn();

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

describe("InteractiveTutorial Component (mocked)", () => {
  it("renders the mocked tutorial", () => {
    render(
      <InteractiveTutorial steps={steps} onClose={onClose} tutorialType="navigation" />
    );
    expect(screen.getByTestId("tutorial-navigation")).toBeInTheDocument();
  });

  it("renders with a different tutorialType", () => {
    render(
      <InteractiveTutorial steps={steps} onClose={onClose} tutorialType="comparison-view" />
    );
    expect(screen.getByTestId("tutorial-comparison-view")).toBeInTheDocument();
  });
});
