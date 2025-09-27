// tests/MapModal.test.tsx
import React from "react";
import { render, screen, fireEvent } from "@testing-library/react";
import MapModal from "../src/components/MapModal";

// MOCK react-leaflet components to avoid rendering real maps
jest.mock("react-leaflet", () => ({
  MapContainer: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  TileLayer: () => <div data-testid="tilelayer" />,
  Marker: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  Popup: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}));

// MOCK lucide-react X icon
jest.mock("lucide-react", () => ({
  X: () => <span>XIcon</span>,
}));

describe("MapModal Component", () => {
  const mockOnClose = jest.fn();
  const mockOnSimulate = jest.fn();

  const intersections = [
    {
      id: "1",
      name: "Main Street",
      details: {
        address: "123 Main St, Block A",
        city: "Pretoria",
        province: "Gauteng",
        latitude: -25.7479,
        longitude: 28.2293,
      },
    },
    {
      id: "2",
      name: "Second Avenue",
      details: {
        address: "456 Second Ave, Block B",
        city: "Pretoria",
        province: "Gauteng",
        latitude: -25.75,
        longitude: 28.23,
      },
    },
  ];

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("renders nothing when isOpen is false", () => {
    const { container } = render(
      <MapModal
        isOpen={false}
        onClose={mockOnClose}
        intersections={intersections}
        onSimulate={mockOnSimulate}
      />
    );
    expect(container.firstChild).toBeNull();
  });

  it("renders modal with title and close button when open", () => {
    render(
      <MapModal
        isOpen={true}
        onClose={mockOnClose}
        intersections={intersections}
        onSimulate={mockOnSimulate}
      />
    );

    expect(screen.getByText("Intersections Map")).toBeInTheDocument();
    const closeButton = screen.getByText("XIcon");
    fireEvent.click(closeButton);
    expect(mockOnClose).toHaveBeenCalled();
  });

  it("renders all intersections with simulate buttons", () => {
    render(
      <MapModal
        isOpen={true}
        onClose={mockOnClose}
        intersections={intersections}
        onSimulate={mockOnSimulate}
      />
    );

    // Each intersection name is rendered
    intersections.forEach((i) => {
      expect(screen.getByText(i.name)).toBeInTheDocument();
      // Simulate button for each intersection
      const button = screen.getByText("Simulate", { selector: "button" });
      fireEvent.click(button);
      expect(mockOnSimulate).toHaveBeenCalledWith(i.id, i.name);
    });
  });
});
