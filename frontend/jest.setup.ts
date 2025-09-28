// jest.setup.ts
import "@testing-library/jest-dom";

jest.mock("/src/config", () => ({
  API_BASE_URL: "http://localhost:3000",
  CHATBOT_BASE_URL: "http://localhost:3000/chatbot",
}));

// Polyfills for TextEncoder/TextDecoder
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

// Mock import.meta.env for Jest
Object.defineProperty(globalThis, "import", {
  value: {
    meta: {
      env: {
        VITE_API_BASE_URL: "http://localhost:3000",
        VITE_CHATBOT_BASE_URL: "http://localhost:3000/chatbot",
      },
    },
  },
});

// Suppress specific React warnings
const originalConsoleError = console.error;
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
