import { useState, useEffect } from "react";
import Navbar from "../components/Navbar";
import UsersTable from "../components/UsersTable";
import "../styles/Users.css";
import Footer from "../components/Footer";
import HelpMenu from "../components/HelpMenu";

// TypeScript interface for user data
interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  lastLogin: string;
}

const Users = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(9);

  const users: User[] = [
    {
      id: 1,
      name: "John Doe",
      email: "email@email.com",
      role: "Admin",
      lastLogin: "2025-05-13 09:00",
    },
    {
      id: 2,
      name: "Jane Smith",
      email: "email@email.com",
      role: "Engineer",
      lastLogin: "2025-05-13 09:00",
    },
    {
      id: 3,
      name: "John Calvin",
      email: "email@email.com",
      role: "Viewer",
      lastLogin: "2025-05-13 09:00",
    },
    {
      id: 4,
      name: "Paul Washer",
      email: "email@email.com",
      role: "Viewer",
      lastLogin: "2025-05-13 09:00",
    },
    {
      id: 5,
      name: "Joshua Garner",
      email: "email@email.com",
      role: "Viewer",
      lastLogin: "2025-05-13 09:00",
    },
    {
      id: 6,
      name: "Chris Xides",
      email: "email@email.com",
      role: "Engineer",
      lastLogin: "2025-05-13 09:00",
    },
    {
      id: 7,
      name: "Kgosi Segale",
      email: "email@email.com",
      role: "Viewer",
      lastLogin: "2025-05-13 09:00",
    },
    {
      id: 8,
      name: "John Flavel",
      email: "email@email.com",
      role: "Viewer",
      lastLogin: "2025-05-13 09:00",
    },
    {
      id: 9,
      name: "John Owen",
      email: "email@email.com",
      role: "Viewer",
      lastLogin: "2025-05-13 09:00",
    },
    {
      id: 10,
      name: "John Doe",
      email: "email@email.com",
      role: "Admin",
      lastLogin: "2025-05-13 09:00",
    },
    {
      id: 11,
      name: "Jane Smith",
      email: "email@email.com",
      role: "Engineer",
      lastLogin: "2025-05-13 09:00",
    },
    {
      id: 12,
      name: "John Calvin",
      email: "email@email.com",
      role: "Viewer",
      lastLogin: "2025-05-13 09:00",
    },
    {
      id: 13,
      name: "Paul Washer",
      email: "email@email.com",
      role: "Viewer",
      lastLogin: "2025-05-13 09:00",
    },
    {
      id: 14,
      name: "Joshua Garner",
      email: "email@email.com",
      role: "Viewer",
      lastLogin: "2025-05-13 09:00",
    },
    {
      id: 15,
      name: "Chris Xides",
      email: "email@email.com",
      role: "Engineer",
      lastLogin: "2025-05-13 09:00",
    },
    {
      id: 16,
      name: "Kgosi Segale",
      email: "email@email.com",
      role: "Viewer",
      lastLogin: "2025-05-13 09:00",
    },
    {
      id: 17,
      name: "John Flavel",
      email: "email@email.com",
      role: "Viewer",
      lastLogin: "2025-05-13 09:00",
    },
    {
      id: 18,
      name: "John Owen",
      email: "email@email.com",
      role: "Viewer",
      lastLogin: "2025-05-13 09:00",
    },
  ];

  useEffect(() => {
    const mediaQuery = window.matchMedia(
      "(max-width: 1400px) and (max-height: 800px)",
    );
    const handleMediaChange = (e: MediaQueryListEvent | MediaQueryList) => {
      setRowsPerPage(e.matches ? 7 : 9);
      setCurrentPage(1);
    };

    handleMediaChange(mediaQuery);
    mediaQuery.addEventListener("change", handleMediaChange);

    return () => mediaQuery.removeEventListener("change", handleMediaChange);
  }, []);

  const totalPages = Math.ceil(users.length / rowsPerPage);

  const startIndex = (currentPage - 1) * rowsPerPage;
  const endIndex = startIndex + rowsPerPage;
  const currentUsers = users.slice(startIndex, endIndex);

  const handleEdit = (id: number) => console.log(`Edit user ${id}`);
  const handleDelete = (id: number) => console.log(`Delete user ${id}`);

  const goToPage = (page: number) => {
    if (page >= 1 && page <= totalPages) setCurrentPage(page);
  };

  const getPageNumbers = () => {
    const pageNumbers: (number | string)[] = [];
    const maxPagesToShow = 5;

    if (totalPages <= maxPagesToShow) {
      return Array.from({ length: totalPages }, (_, i) => i + 1);
    }

    const leftBound = Math.max(2, currentPage - 1);
    const rightBound = Math.min(totalPages - 1, currentPage + 1);

    pageNumbers.push(1);
    if (leftBound > 2) pageNumbers.push("...");

    for (let i = leftBound; i <= rightBound; i++) {
      pageNumbers.push(i);
    }

    if (rightBound < totalPages - 1) pageNumbers.push("...");
    if (totalPages > 1) pageNumbers.push(totalPages);

    return pageNumbers;
  };

  return (
    <div className="userBody min-h-screen bg-gray-100">
      <Navbar />
      <div className="user-main-content flex-grow">
        <div className="max-w-6xl mx-auto px-4 py-8">
          <UsersTable
            users={currentUsers}
            onEdit={handleEdit}
            onDelete={handleDelete}
          />

          <div className="usersPaging flex justify-center items-center py-4 gap-2 mt-4">
            <button
              onClick={() => goToPage(currentPage - 1)}
              disabled={currentPage === 1}
              className="px-4 py-2 bg-gray-200 text-gray-700 rounded-full hover:bg-gray-300 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
              aria-label="Previous page"
            >
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M15 19l-7-7 7-7"
                />
              </svg>
            </button>

            {getPageNumbers().map((page, index) =>
              typeof page === "string" ? (
                <span key={index} className="px-4 py-2 text-gray-500">
                  ...
                </span>
              ) : (
                <button
                  key={page}
                  onClick={() => goToPage(page)}
                  className={`px-4 py-2 rounded-full transition-colors duration-200 ${
                    currentPage === page
                      ? "bg-blue-500 text-white"
                      : "bg-gray-200 text-gray-700 hover:bg-gray-300"
                  }`}
                >
                  {page}
                </button>
              ),
            )}

            <button
              onClick={() => goToPage(currentPage + 1)}
              disabled={currentPage === totalPages}
              className="px-4 py-2 bg-gray-200 text-gray-700 rounded-full hover:bg-gray-300 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
              aria-label="Next page"
            >
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M9 5l7 7-7 7"
                />
              </svg>
            </button>
          </div>
        </div>
      </div>
      <Footer />
      <HelpMenu />
    </div>
  );
};

export default Users;
