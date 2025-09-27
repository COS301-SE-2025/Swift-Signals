// jest.setup.ts
import '@testing-library/jest-dom';

// jest.setup.ts
class PolyTextEncoder {
  encode(input: string) {
    const result = new Uint8Array(input.length);
    for (let i = 0; i < input.length; i++) {
      result[i] = input.charCodeAt(i);
    }
    return result;
  }
}

class PolyTextDecoder {
  decode(input: Uint8Array) {
    let result = "";
    for (let i = 0; i < input.length; i++) {
      result += String.fromCharCode(input[i]);
    }
    return result;
  }
}

// @ts-ignore
global.TextEncoder = PolyTextEncoder;
// @ts-ignore
global.TextDecoder = PolyTextDecoder;

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
