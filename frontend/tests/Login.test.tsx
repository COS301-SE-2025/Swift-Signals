import React from "react";
import { render, screen} from "@testing-library/react";
import '@testing-library/jest-dom';
import { BrowserRouter } from "react-router-dom";
import Dashboard from "../src/pages/Dashboard";

console.log(React)

const mockNavigate = jest.fn();
jest.mock("react-router-dom", () => ({
  ...jest.requireActual("react-router-dom"),
  useNavigate: () => mockNavigate,
}));

jest.mock("../src/components/Footer", () => () => <div>Footer</div>);
jest.mock("../src/components/Navbar", () => () => <div>Navbar</div>);
jest.mock("../src/components/HelpMenu", () => () => <div>HelpMenu</div>);
// eslint-disable-next-line @typescript-eslint/no-explicit-any
jest.mock("../src/components/MapModal", () => (props: any) => (
  <div>MapModal {props.isOpen ? "(Open)" : "(Closed)"}</div>
));

jest.mock("chart.js", () => {
  const actual = jest.requireActual("chart.js");
  return {
    ...actual,
    Chart: jest.fn().mockImplementation(() => ({
      destroy: jest.fn(),
    })),
    registerables: [],
  };
});

global.fetch = jest.fn();

describe("Dashboard", () => {
  beforeEach(() => {
    (fetch as jest.Mock).mockReset();
    localStorage.clear();
  });

  it("renders main dashboard elements", async () => {
    (fetch as jest.Mock).mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        intersections: [
          { id: "1", name: "Int 1", traffic_density: "TRAFFIC_DENSITY_HIGH" },
        ],
      }),
    });

    render(
      <BrowserRouter>
        <Dashboard />
      </BrowserRouter>
    );

    expect(screen.getByText("Navbar")).toBeInTheDocument();
    expect(screen.getByText("Footer")).toBeInTheDocument();
    expect(screen.getByText("HelpMenu")).toBeInTheDocument();
    expect(await screen.findByText("Total Intersections")).toBeInTheDocument();
    expect(await screen.findByText("1")).toBeInTheDocument();
  });
});
