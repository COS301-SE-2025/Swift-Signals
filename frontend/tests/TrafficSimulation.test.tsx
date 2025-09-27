import React from "react";
import { render, screen, act, fireEvent } from "@testing-library/react";
import TrafficSimulation from "../src/pages/TrafficSimulation";

console.log(React)

jest.mock("@react-three/fiber", () => ({
  Canvas: ({ children }: any) => <div data-testid="mock-canvas">{children}</div>,
  useFrame: jest.fn(),
}));
jest.mock("@react-three/drei", () => ({
  MapControls: () => <div data-testid="mock-map-controls" />,
  OrthographicCamera: () => <div data-testid="mock-camera" />,
}));

jest.mock("../src/components/SimulationUI", () => ({
  SimulationUI: (props: any) => (
    <div data-testid="simulation-ui">
      <button onClick={props.onPlayPause}>Toggle Play</button>
      <button onClick={props.onRestart}>Restart</button>
      <button onClick={() => props.onSpeedChange(10)}>Change Speed</button>
    </div>
  ),
}));

const mockFetch = jest.fn();
(global as any).fetch = mockFetch;
const mockGetItem = jest.fn();
Object.defineProperty(window, "localStorage", {
  value: {
    getItem: mockGetItem,
    setItem: jest.fn(),
    removeItem: jest.fn(),
    clear: jest.fn(),
  },
});

describe("TrafficSimulation Page", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("renders loading state initially", async () => {
    mockGetItem.mockReturnValueOnce("fake-token");
    mockFetch.mockImplementationOnce(() => new Promise(() => {}));

    await act(async () => {
      render(<TrafficSimulation intersectionId="1" isExpanded={true} />);
    });

    expect(
      screen.getByText(/Loading simulation data/i)
    ).toBeInTheDocument();
  });

  it("shows error if no intersectionId is provided", async () => {
    await act(async () => {
      render(<TrafficSimulation intersectionId="" isExpanded={true} />);
    });
    expect(
      screen.getByText(/No intersection ID provided/i),
    ).toBeInTheDocument();
  });

  it("shows error if auth token is missing", async () => {
    mockGetItem.mockReturnValueOnce(null);
    await act(async () => {
      render(<TrafficSimulation intersectionId="123" isExpanded={true} />);
    });
    expect(
      screen.getByText(/Authentication token not found/i),
    ).toBeInTheDocument();
  });

  it("renders with provided simulationData (bypasses fetch)", async () => {
    const mockSimData = {
      intersection: {
        nodes: [{ id: "n1", x: 0, y: 0, type: "road" }],
        edges: [],
        connections: [],
      },
      vehicles: [],
    };

    await act(async () => {
      render(
        <TrafficSimulation
          intersectionId="123"
          isExpanded={true}
          simulationData={mockSimData}
        />,
      );
    });

    expect(screen.getByTestId("mock-canvas")).toBeInTheDocument();
    expect(screen.getByTestId("simulation-ui")).toBeInTheDocument();
  });

  it("allows play/pause, restart, and speed change", async () => {
    const mockSimData = {
      intersection: {
        nodes: [{ id: "n1", x: 0, y: 0, type: "road" }],
        edges: [],
        connections: [],
      },
      vehicles: [],
    };

    await act(async () => {
      render(
        <TrafficSimulation
          intersectionId="123"
          isExpanded={true}
          simulationData={mockSimData}
        />,
      );
    });

    const playPauseBtn = screen.getByText(/Toggle Play/i);
    const restartBtn = screen.getByText(/Restart/i);
    const speedBtn = screen.getByText(/Change Speed/i);

    act(() => {
      fireEvent.click(playPauseBtn);
      fireEvent.click(restartBtn);
      fireEvent.click(speedBtn);
    });

    expect(screen.getByTestId("simulation-ui")).toBeInTheDocument();
  });

  it("shows retry button when error occurs", async () => {
    mockGetItem.mockReturnValueOnce("fake-token");
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 404,
      statusText: "Not Found",
    });

    await act(async () => {
      render(<TrafficSimulation intersectionId="123" isExpanded={true} />);
    });

    expect(
      screen.getByText(/Simulation data not found/i),
    ).toBeInTheDocument();

    const retryBtn = screen.getByText(/Retry/i);
    expect(retryBtn).toBeInTheDocument();
  });
});
