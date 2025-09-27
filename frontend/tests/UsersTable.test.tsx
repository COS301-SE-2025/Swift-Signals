// tests/UsersTable.test.tsx
import React from "react";
import { render, screen, fireEvent } from "@testing-library/react";
import UsersTable from "../src/components/UsersTable";

console.log(React)

describe("UsersTable", () => {
  const mockUsers = [
    {
      displayId: 1,
      id: "u1",
      name: "Alice Johnson",
      email: "alice@example.com",
      role: "Admin",
      lastLogin: "2025-09-25",
    },
    {
      displayId: 2,
      id: "u2",
      name: "Bob Smith",
      email: "bob@example.com",
      role: "User",
      lastLogin: "2025-09-26",
    },
  ];

  const onEdit = jest.fn();
  const onDelete = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("renders table headers correctly", () => {
    render(<UsersTable users={mockUsers} onEdit={onEdit} onDelete={onDelete} />);
    expect(screen.getByText("ID")).toBeInTheDocument();
    expect(screen.getByText("Name")).toBeInTheDocument();
    expect(screen.getByText("Email")).toBeInTheDocument();
    expect(screen.getByText("Role")).toBeInTheDocument();
    expect(screen.getByText("Last Login")).toBeInTheDocument();
    expect(screen.getByText("Actions")).toBeInTheDocument();
  });

  it("renders user rows with correct data", () => {
    render(<UsersTable users={mockUsers} onEdit={onEdit} onDelete={onDelete} />);

    expect(screen.getByText("1")).toBeInTheDocument();
    expect(screen.getByText("Alice Johnson")).toBeInTheDocument();
    expect(screen.getByText("alice@example.com")).toBeInTheDocument();
    expect(screen.getByText("Admin")).toBeInTheDocument();
    expect(screen.getByText("2025-09-25")).toBeInTheDocument();

    expect(screen.getByText("2")).toBeInTheDocument();
    expect(screen.getByText("Bob Smith")).toBeInTheDocument();
    expect(screen.getByText("bob@example.com")).toBeInTheDocument();
    expect(screen.getByText("User")).toBeInTheDocument();
    expect(screen.getByText("2025-09-26")).toBeInTheDocument();
  });

  it("calls onEdit when edit button is clicked", () => {
    render(<UsersTable users={mockUsers} onEdit={onEdit} onDelete={onDelete} />);
    const editButtons = screen.getAllByRole("button", { name: /edit user/i });

    fireEvent.click(editButtons[0]);
    expect(onEdit).toHaveBeenCalledTimes(1);
    expect(onEdit).toHaveBeenCalledWith("u1");

    fireEvent.click(editButtons[1]);
    expect(onEdit).toHaveBeenCalledWith("u2");
  });

  it("calls onDelete when delete button is clicked", () => {
    render(<UsersTable users={mockUsers} onEdit={onEdit} onDelete={onDelete} />);
    const deleteButtons = screen.getAllByRole("button", { name: /delete user/i });

    fireEvent.click(deleteButtons[0]);
    expect(onDelete).toHaveBeenCalledTimes(1);
    expect(onDelete).toHaveBeenCalledWith("u1");

    fireEvent.click(deleteButtons[1]);
    expect(onDelete).toHaveBeenCalledWith("u2");
  });

  it("renders no rows when users list is empty", () => {
    render(<UsersTable users={[]} onEdit={onEdit} onDelete={onDelete} />);
    // Table still renders headers, but no user data
    expect(screen.getByText("ID")).toBeInTheDocument();
    expect(screen.queryByText("Alice Johnson")).not.toBeInTheDocument();
  });
});
