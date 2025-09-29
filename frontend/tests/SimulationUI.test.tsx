import React from "react";
import { render, screen, fireEvent } from "@testing-library/react";
import { SimulationUI } from "../src/components/SimulationUI";

console.log(React);

describe("SimulationUI Component", () => {
  const mockPlayPause = jest.fn();
  const mockRestart = jest.fn();
  const mockSpeedChange = jest.fn();

  const defaultProps = {
    time: 12.5,
    vehicleCount: 10,
    isPlaying: false,
    speed: 5,
    onPlayPause: mockPlayPause,
    onRestart: mockRestart,
    onSpeedChange: mockSpeedChange,
    trafficLightStates: { N: "g", S: "r", E: "y", W: "g" },
    activeVehicles: 7,
    completedVehicles: 3,
    avgSpeed: 15,
    progress: 0.42,
    totalSimTime: 30,
    scale: 1,
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("renders main sections and data correctly", () => {
    render(<SimulationUI {...defaultProps} />);

    expect(screen.getByText("Simulation")).toBeInTheDocument();

    expect(screen.getByText("42%")).toBeInTheDocument();
    expect(screen.getByText("Progress")).toBeInTheDocument();

    expect(
      screen.getByText(
        `${defaultProps.time.toFixed(1)} / ${defaultProps.totalSimTime.toFixed(1)} s`,
      ),
    ).toBeInTheDocument();

    expect(screen.getByText("Total")).toBeInTheDocument();
    expect(
      screen.getByText(defaultProps.vehicleCount.toString()),
    ).toBeInTheDocument();
    expect(screen.getByText("Active")).toBeInTheDocument();
    expect(
      screen.getByText(defaultProps.activeVehicles.toString()),
    ).toBeInTheDocument();
    expect(screen.getByText("Completed")).toBeInTheDocument();
    expect(
      screen.getByText(defaultProps.completedVehicles.toString()),
    ).toBeInTheDocument();

    const avgSpeedKmh = (defaultProps.avgSpeed * 3.6).toFixed(1);
    expect(screen.getByText(`${avgSpeedKmh} km/h`)).toBeInTheDocument();

    Object.entries(defaultProps.trafficLightStates).forEach(([dir]) => {
      expect(screen.getByText(dir)).toBeInTheDocument();
    });

    expect(screen.getByText("Play")).toBeInTheDocument();
    expect(screen.getByText("Restart")).toBeInTheDocument();

    expect(
      screen.getByDisplayValue(defaultProps.speed.toString()),
    ).toBeInTheDocument();
  });

  it("calls onPlayPause and onRestart on button click", () => {
    render(<SimulationUI {...defaultProps} />);
    fireEvent.click(screen.getByText("Play"));
    expect(mockPlayPause).toHaveBeenCalled();

    fireEvent.click(screen.getByText("Restart"));
    expect(mockRestart).toHaveBeenCalled();
  });

  it("calls onSpeedChange when slider is changed", () => {
    render(<SimulationUI {...defaultProps} />);
    const slider = screen.getByRole("slider") as HTMLInputElement;
    fireEvent.change(slider, { target: { value: "10" } });
    expect(mockSpeedChange).toHaveBeenCalledWith(10);
  });

  it("updates Play/Pause button text when isPlaying prop changes", () => {
    const { rerender } = render(
      <SimulationUI {...defaultProps} isPlaying={false} />,
    );
    expect(screen.getByText("Play")).toBeInTheDocument();

    rerender(<SimulationUI {...defaultProps} isPlaying={true} />);
    expect(screen.getByText("Pause")).toBeInTheDocument();
  });
});
