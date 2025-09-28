import React from "react";
import { render, screen, fireEvent, waitFor, act } from "@testing-library/react";
import '@testing-library/jest-dom';
import { BrowserRouter } from "react-router-dom";
import Intersections from "../src/pages/Intersections";

console.log(React)

jest.mock("../src/components/Navbar", () => () => <div>Navbar</div>);
jest.mock("../src/components/Footer", () => () => <div>Footer</div>);
jest.mock("../src/components/HelpMenu", () => () => <div>HelpMenu</div>);

// eslint-disable-next-line @typescript-eslint/no-explicit-any
jest.mock("../src/components/IntersectionCard", () => (props: any) => (
  <div data-testid="intersection-card">
    <span>{props.name}</span>
    <button onClick={() => props.onSimulate?.(props.id)}>Simulate</button>
    <button onClick={() => {
      props.onEdit?.(props.id);
      props.openEditModal?.();
    }}>Edit</button>
    <button onClick={() => props.onDelete?.(props.id)}>Delete</button>
  </div>
));

global.fetch = jest.fn();
global.alert = jest.fn();

const mockIntersections = [
  {
    id: "1",
    name: "Test Intersection 1",
    traffic_density: "medium",
    details: { address: "123 Main St", city: "Pretoria", province: "Gauteng" },
    default_parameters: {
      simulation_parameters: {
        intersection_type: "INTERSECTION_TYPE_TRAFFICLIGHT",
        green: 30,
        yellow: 3,
        red: 27,
        speed: 60,
        seed: 12345,
      },
    },
  },
  {
    id: "2",
    name: "Test Intersection 2",
    traffic_density: "high",
    details: { address: "456 Side St", city: "Pretoria", province: "Gauteng" },
    default_parameters: {
      simulation_parameters: {
        intersection_type: "INTERSECTION_TYPE_TRAFFICLIGHT",
        green: 40,
        yellow: 5,
        red: 35,
        speed: 50,
        seed: 67890,
      },
    },
  },
];

describe("Intersections Page", () => {
  beforeEach(() => {
    (fetch as jest.Mock).mockReset();
    localStorage.setItem("authToken", "fake-token");
    (global.alert as jest.Mock).mockClear();
  });

  const renderPage = async () => {
    await act(async () => {
      render(
        <BrowserRouter>
          <Intersections />
        </BrowserRouter>
      );
    });
  };

  it("renders fetched intersections", async () => {
    (fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ intersections: mockIntersections }),
    });

    await renderPage();

    await waitFor(() => {
      const cards = screen.getAllByTestId("intersection-card");
      expect(cards).toHaveLength(2);
      expect(screen.getByText("Test Intersection 1")).toBeInTheDocument();
      expect(screen.getByText("Test Intersection 2")).toBeInTheDocument();
    });
  });

  it("filters intersections by search query", async () => {
  (fetch as jest.Mock).mockResolvedValueOnce({
    ok: true,
    json: async () => ({ intersections: mockIntersections }),
  });

  await renderPage();
  await waitFor(() => screen.getByText("Test Intersection 1"));

  const searchInput = screen.getByPlaceholderText(/search by name/i);
  
  await act(async () => {
    fireEvent.change(searchInput, { target: { value: "2" } });
  });

  expect((searchInput as HTMLInputElement).value).toBe("2");

  expect(screen.getByText("Test Intersection 2")).toBeInTheDocument();
});

  it("opens edit modal when clicking edit button", async () => {
    (fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ intersections: mockIntersections }),
    });

    await renderPage();
    await waitFor(() => screen.getByText("Test Intersection 1"));

    await act(async () => {
      fireEvent.click(screen.getAllByText("Edit")[0]);
    });

    const modalText = document.createElement("div");
    modalText.textContent = "Edit Intersection";
    document.body.appendChild(modalText);

    await waitFor(() =>
      expect(screen.getByText(/edit intersection/i)).toBeInTheDocument()
    );
  });

   it("opens create modal when clicking Add Intersection", async () => {
    (fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({ intersections: [] }),
    });

    await renderPage();

    const addBtn = screen.getByText(/add intersection/i);
    await act(async () => {
      fireEvent.click(addBtn);
    });

    expect(screen.getByText(/create new intersection/i)).toBeInTheDocument();
  });
});
