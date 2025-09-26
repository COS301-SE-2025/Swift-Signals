// jest.setup.ts
import '@testing-library/jest-dom';

// jest.setup.ts

// Save original console.error
const originalConsoleError = console.error;

// Suppress specific React warning
beforeAll(() => {
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

afterAll(() => {
  console.error = originalConsoleError;
});
