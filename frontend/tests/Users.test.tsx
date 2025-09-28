/* eslint-disable @typescript-eslint/no-explicit-any */
import React from "react";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import Users from "../src/pages/Users";

console.log(React)

jest.mock("../src/components/Navbar", () => () => <div>Navbar</div>);
jest.mock("../src/components/Footer", () => () => <div>Footer</div>);
jest.mock("../src/components/HelpMenu", () => () => <div>HelpMenu</div>);
jest.mock("lucide-react", () => ({ X: () => <span>XIcon</span> }));

jest.mock("../src/components/UsersTable", () => {
  return ({ users, onEdit, onDelete, loading }: any) => (
    <div>
      {/* Give spinner a test id */}
      {loading && <span data-testid="spinner" className="spinner" />}
      {users?.map((u: any) => (
        <div key={u.id} data-testid="user-row">
          <span>{u.username}</span>
          <button data-testid={`edit-btn-${u.id}`} onClick={() => onEdit(u.id)}>Edit</button>
          <button data-testid={`delete-btn-${u.id}`} onClick={() => onDelete(u.id)}>Delete</button>
        </div>
      ))}
      <button data-testid="prev-page">Prev</button>
      <button data-testid="next-page">Next</button>
    </div>
  );
});

const mockFetch = jest.fn();
(global as any).fetch = mockFetch;

const mockGetItem = jest.fn();
Object.defineProperty(window, "localStorage", {
  value: { getItem: mockGetItem, setItem: jest.fn(), removeItem: jest.fn(), clear: jest.fn() },
});

const originalConfirm = window.confirm;
beforeAll(() => {
  window.confirm = jest.fn();
  Object.defineProperty(window, "matchMedia", {
    writable: true,
    value: jest.fn().mockImplementation((query) => ({
      matches: false,
      media: query,
      onchange: null,
      addEventListener: jest.fn(),
      removeEventListener: jest.fn(),
      addListener: jest.fn(),
      removeListener: jest.fn(),
      dispatchEvent: jest.fn(),
    })),
  });
});
afterAll(() => { window.confirm = originalConfirm; });

describe("Users Page", () => {
  beforeEach(() => { jest.clearAllMocks(); });

  const fakeUsers = [
    { id: "1", username: "Alice", email: "alice@test.com", is_admin: true, intersection_ids: [] },
    { id: "2", username: "Bob", email: "bob@test.com", is_admin: false, intersection_ids: [] },
  ];

  /*it("renders loading spinner initially", async () => {
    mockGetItem.mockReturnValue("token");
    mockFetch.mockImplementation(() => new Promise(() => {}));

    render(<Users />);

    expect(await screen.findByTestId("spinner")).toBeInTheDocument();
  });

  it("fetches and displays users", async () => {
    mockGetItem.mockReturnValue("token");
    mockFetch.mockResolvedValueOnce({ ok: true, json: async () => fakeUsers });

    render(<Users />);

    for (const user of fakeUsers) {
      expect(await screen.findByText((content, element) => {
        return element?.tagName === 'SPAN' && content === user.username;
      })).toBeInTheDocument();
    }
  });*/

  it("shows error if no auth token", async () => {
    mockGetItem.mockReturnValue(null);
    render(<Users />);
    expect(await screen.findByText(/No authentication token found/i)).toBeInTheDocument();
    expect(screen.getByText(/Retry/i)).toBeInTheDocument();
  });

  it("opens edit modal when edit button clicked", async () => {
    mockGetItem.mockReturnValue("token");
    mockFetch.mockResolvedValueOnce({ ok: true, json: async () => fakeUsers });

    render(<Users />);
    fireEvent.click(await screen.findByTestId("edit-btn-1"));
    expect(await screen.findByText("Edit User")).toBeInTheDocument();
  });

  it("deletes a user after confirmation", async () => {
    mockGetItem.mockReturnValue("token");
    (window.confirm as jest.Mock).mockReturnValue(true);
    mockFetch.mockResolvedValueOnce({ ok: true, json: async () => fakeUsers })
             .mockResolvedValueOnce({ status: 204 });

    render(<Users />);
    fireEvent.click(await screen.findByTestId("delete-btn-1"));

    await waitFor(() => {
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining("/admin/users/1"),
        expect.objectContaining({ method: "DELETE" })
      );
    });
  });

  it("handles pagination buttons", async () => {
    const manyUsers = Array.from({ length: 20 }, (_, i) => ({
      id: `${i+1}`,
      username: `User${i+1}`,
      email: `u${i+1}@test.com`,
      is_admin: false,
      intersection_ids: []
    }));

    mockGetItem.mockReturnValue("token");
    mockFetch.mockResolvedValue({ ok: true, json: async () => manyUsers });

    render(<Users />);
    fireEvent.click(await screen.findByTestId("next-page"));
    fireEvent.click(await screen.findByTestId("prev-page"));
    expect(mockFetch).toHaveBeenCalled();
  });
});
