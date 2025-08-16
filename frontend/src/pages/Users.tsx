import { useState, useEffect } from "react";
import Navbar from "../components/Navbar";
import UsersTable from "../components/UsersTable";
import "../styles/Users.css";
import Footer from "../components/Footer";
import HelpMenu from "../components/HelpMenu";

// TypeScript interface for user data from API
interface ApiUser {
  id: string;
  username: string;
  email: string;
  is_admin: boolean;
  intersection_ids: string[];
}

// TypeScript interface for display user data
interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  lastLogin: string;
}

// API response interface
interface ApiResponse {
  users?: ApiUser[];
  message?: string;
}

const Users = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(9);
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [totalUsers, setTotalUsers] = useState(0);

  // Function to get auth token from localStorage or wherever it's stored
  const getAuthToken = () => {
    return localStorage.getItem('authToken') || sessionStorage.getItem('authToken');
  };

  // Function to fetch users from API
  const fetchUsers = async (page: number, pageSize: number) => {
    try {
      setLoading(true);
      setError(null);
      
      const token = getAuthToken();
      if (!token) {
        throw new Error('No authentication token found');
      }

      // Create URL with query parameters
      const url = new URL('/admin/users', window.location.origin);
      url.searchParams.append('page', page.toString());
      url.searchParams.append('page_size', pageSize.toString());

      const response = await fetch(url.toString(), {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        }
      });

      if (!response.ok) {
        if (response.status === 401) {
          throw new Error('Unauthorized - Please log in again');
        } else if (response.status === 403) {
          throw new Error('Forbidden - Admin access required');
        } else {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
      }

      const data: ApiUser[] = await response.json();
      
      // Transform API data to display format
      const transformedUsers: User[] = data.map((apiUser, index) => ({
        id: parseInt(apiUser.id),
        name: apiUser.username,
        email: apiUser.email,
        role: apiUser.is_admin ? 'Admin' : 'User',
        lastLogin: 'N/A' // API doesn't provide last login, you might need to add this field
      }));

      setUsers(transformedUsers);
      setTotalUsers(data.length); // Note: You might need total count from API for proper pagination
      
    } catch (err) {
      console.error('Error fetching users:', err);
      setError(err instanceof Error ? err.message : 'Failed to fetch users');
      setUsers([]);
    } finally {
      setLoading(false);
    }
  };

  // Fetch users when component mounts or pagination changes
  useEffect(() => {
    fetchUsers(currentPage, rowsPerPage);
  }, [currentPage, rowsPerPage]);

  // Handle responsive design
  useEffect(() => {
    const mediaQuery = window.matchMedia(
      "(max-width: 1400px) and (max-height: 800px)",
    );
    const handleMediaChange = (e: MediaQueryListEvent | MediaQueryList) => {
      const newRowsPerPage = e.matches ? 7 : 9;
      setRowsPerPage(newRowsPerPage);
      setCurrentPage(1);
    };

    handleMediaChange(mediaQuery);
    mediaQuery.addEventListener("change", handleMediaChange);

    return () => mediaQuery.removeEventListener("change", handleMediaChange);
  }, []);

  // Calculate pagination
  const totalPages = Math.ceil(totalUsers / rowsPerPage);

  // Handle user actions
  const handleEdit = async (id: number) => {
    try {
      const token = getAuthToken();
      if (!token) {
        throw new Error('No authentication token found');
      }

      // You can implement edit functionality here
      console.log(`Edit user ${id}`);
      
      // Example: Navigate to edit page or open modal
      // For now, just log the action
    } catch (err) {
      console.error('Error editing user:', err);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this user?')) {
      return;
    }

    try {
      const token = getAuthToken();
      if (!token) {
        throw new Error('No authentication token found');
      }

      const response = await fetch(`/admin/users/${id}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        if (response.status === 401) {
          throw new Error('Unauthorized - Please log in again');
        } else if (response.status === 403) {
          throw new Error('Forbidden - Admin access required');
        } else if (response.status === 404) {
          throw new Error('User not found');
        } else {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
      }

      // Refresh the users list after successful deletion
      await fetchUsers(currentPage, rowsPerPage);
      
      // If current page is empty after deletion, go to previous page
      if (users.length === 1 && currentPage > 1) {
        setCurrentPage(currentPage - 1);
      }

    } catch (err) {
      console.error('Error deleting user:', err);
      alert(err instanceof Error ? err.message : 'Failed to delete user');
    }
  };

  const goToPage = (page: number) => {
    if (page >= 1 && page <= totalPages) {
      setCurrentPage(page);
    }
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
        <div className="usersDisp max-w-6xl mx-auto px-4 py-8">
          
          {/* Loading State */}
          {loading && (
            <div className="flex justify-center items-center py-8">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-[#0F5BA7]"></div>
              <span className="ml-2 text-gray-600">Loading users...</span>
            </div>
          )}

          {/* Error State */}
          {error && (
            <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
              <strong>Error:</strong> {error}
              <button
                onClick={() => fetchUsers(currentPage, rowsPerPage)}
                className="ml-4 bg-red-500 text-white px-3 py-1 rounded text-sm hover:bg-red-600"
              >
                Retry
              </button>
            </div>
          )}

          {/* Users Table */}
          {!loading && !error && (
            <>
              <UsersTable
                users={users}
                onEdit={handleEdit}
                onDelete={handleDelete}
              />

              {/* Pagination */}
              {totalPages > 1 && (
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
                            ? "bg-[#0F5BA7] text-white"
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
              )}
            </>
          )}

          {/* Empty State */}
          {!loading && !error && users.length === 0 && (
            <div className="text-center py-8">
              <p className="text-gray-500 text-lg">No users found.</p>
            </div>
          )}
        </div>
      </div>
      <Footer />
      <HelpMenu />
    </div>
  );
};

export default Users;