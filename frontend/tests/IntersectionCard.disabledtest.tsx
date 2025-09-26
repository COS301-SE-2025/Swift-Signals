// tests/IntersectionCard.test.tsx
import React from "react";
import { render, screen, fireEvent } from "@testing-library/react";
import IntersectionCard from "../src/components/IntersectionCard";

console.log(React)

// MOCK ALL IMAGE IMPORTS TO A STRING SO THEY DON'T TRIGGER FILE RESOLUTION
jest.mock("../src/assets/placeholder.png", () => "", { virtual: true });

describe("IntersectionCard Component", () => {
  const mockSimulate = jest.fn();
  const mockEdit = jest.fn();
  const mockDelete = jest.fn();

  const props = {
    id: "123",
    name: "Main Street [Downtown]",
    location: "City Center, Block A",
    lanes: "4-way",
    image: "", // triggers placeholder
    onSimulate: mockSimulate,
    onEdit: mockEdit,
    onDelete: mockDelete,
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("renders name, location, lanes, and placeholder image", () => {
    render(<IntersectionCard {...props} />);
    expect(screen.getByText("Main Street")).toBeInTheDocument();
    expect(screen.getByText("Location: City Center")).toBeInTheDocument();
    expect(screen.getByText("Type: 4-way")).toBeInTheDocument();

    const img = screen.getAllByRole("img")[0];
    expect(img).toHaveAttribute("alt", props.name);
  });

  it("calls onSimulate when Simulate button is clicked", () => {
    render(<IntersectionCard {...props} />);
    fireEvent.click(screen.getAllByText("Simulate")[0]);
    expect(mockSimulate).toHaveBeenCalledWith("123");
  });

  it("calls onEdit when Edit button is clicked", () => {
    render(<IntersectionCard {...props} />);
    fireEvent.click(screen.getAllByText("Edit")[0]);
    expect(mockEdit).toHaveBeenCalledWith("123");
  });

  it("calls onDelete when Delete button is clicked", () => {
    render(<IntersectionCard {...props} />);
    fireEvent.click(screen.getAllByText("Delete")[0]);
    expect(mockDelete).toHaveBeenCalledWith("123");
  });

  it("renders at least 6 buttons including mobile buttons", () => {
    render(<IntersectionCard {...props} />);
    const buttons = screen.getAllByRole("button");
    expect(buttons.length).toBeGreaterThanOrEqual(6);
  });
});
