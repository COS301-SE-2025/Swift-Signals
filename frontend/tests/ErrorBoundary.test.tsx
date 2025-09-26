import { render, screen } from "@testing-library/react";
import ErrorBoundary from "../src/components/ErrorBoundary";

describe("ErrorBoundary", () => {
  const originalConsoleError = console.error;

  beforeEach(() => {
    console.error = jest.fn();
  });

  afterEach(() => {
    console.error = originalConsoleError;
  });

  it("renders children when no error occurs", () => {
    render(
      <ErrorBoundary>
        <div>Child content</div>
      </ErrorBoundary>
    );
    expect(screen.getByText("Child content")).toBeInTheDocument();
  });

  it("renders fallback UI when a child throws an error", () => {
    const ProblemChild = () => {
      throw new Error("Test error");
    };

    render(
      <ErrorBoundary>
        <ProblemChild />
      </ErrorBoundary>
    );

    // Fallback UI is displayed
    expect(screen.getByText("Something went wrong.")).toBeInTheDocument();

    // console.error was called
    expect(console.error).toHaveBeenCalled();

    // Grab all arguments passed to console.error
    const errorArgs = (console.error as jest.Mock).mock.calls[0];

    // Type 'unknown' for safety in TS
    const error = (errorArgs as unknown[]).find((arg: unknown) => arg instanceof Error) as Error | undefined;

    expect(error).toBeInstanceOf(Error);
    expect(error?.message).toBe("Test error");
  });
});
