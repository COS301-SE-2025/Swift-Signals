import { X } from "lucide-react";
import { useState, useEffect } from "react";

import Footer from "../components/Footer";
import HelpMenu from "../components/HelpMenu";
import Navbar from "../components/Navbar";
import UsersTable from "../components/UsersTable";
import "../styles/Users.css";
import { API_BASE_URL } from "../config";

interface ApiUser {
  id: string;
  username: string;
  email: string;
  is_admin: boolean;
  intersection_ids: string[];
}

interface User {
  displayId: number;
  id: string;
  name: string;
  email: string;
  role: string;
  lastLogin: string;
}

interface UpdateUserPayload {
  username?: string;
  email?: string;
}

interface EditUserFormData {
  username: string;
  email: string;
}

interface EditUserModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (id: string, data: EditUserFormData) => void;
  user: User | null;
  isLoading: boolean;
  error: string | null;
}

const EditUserModal: React.FC<EditUserModalProps> = ({
  isOpen,
  onClose,
  onSubmit,
  user,
  isLoading,
  error,
}) => {
  const [formData, setFormData] = useState<EditUserFormData>({
    username: "",
    email: "",
  });

  useEffect(() => {
    if (user) {
      setFormData({
        username: user.name,
        email: user.email,
      });
    }
  }, [user]);

  if (!isOpen) return null;

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (user) {
      onSubmit(user.id, formData);
    }
  };

  const inputClasses =
    "mt-1 block w-full px-3 py-2 bg-gray-50 dark:bg-[#0D1117] border-2 border-gray-300 dark:border-[#30363D] rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2 dark:focus:ring-offset-[#161B22] focus:ring-[#2da44e] sm:text-sm text-gray-900 dark:text-[#C9D1D9]";

  return (
    <div className="fixed inset-0 bg-black bg-opacity-60 backdrop-blur-sm flex justify-center items-center z-50 p-4">
      <div className="bg-white dark:bg-[#161B22] p-8 rounded-xl shadow-2xl w-full max-w-lg relative border border-gray-200 dark:border-[#30363D]">
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-gray-400 dark:text-[#7D8590] hover:text-gray-600 dark:hover:text-[#E6EDF3] transition-colors duration-150"
        >
          <X size={24} />
        </button>
        <h2 className="text-2xl font-bold mb-6 text-center text-gray-900 dark:text-[#E6EDF3]">
          Edit User
        </h2>
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label
              htmlFor="username"
              className="block text-sm font-medium text-gray-700 dark:text-[#C9D1D9] mb-1"
            >
              Username
            </label>
            <input
              type="text"
              name="username"
              id="username"
              required
              minLength={3}
              maxLength={32}
              className={inputClasses}
              value={formData.username}
              onChange={handleChange}
            />
            <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
              Must be between 3 and 32 characters.
            </p>
          </div>
          <div>
            <label
              htmlFor="email"
              className="block text-sm font-medium text-gray-700 dark:text-[#C9D1D9] mb-1"
            >
              Email
            </label>
            <input
              type="email"
              name="email"
              id="email"
              required
              className={inputClasses}
              value={formData.email}
              onChange={handleChange}
            />
          </div>
          {error && <p className="text-red-500 text-sm text-center">{error}</p>}
          <div className="flex justify-end space-x-4 pt-4 border-t border-gray-200 dark:border-[#30363D] mt-6">
            <button
              type="button"
              onClick={onClose}
              className="px-6 py-2 bg-gray-100 dark:bg-[#21262D] border-2 border-gray-300 dark:border-[#30363D] text-gray-700 dark:text-[#C9D1D9] rounded-lg font-medium hover:bg-gray-200 dark:hover:bg-[#30363D] transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isLoading}
              className="px-6 py-2 bg-[#2da44e] text-white rounded-lg font-medium hover:bg-[#288c42] disabled:opacity-60 disabled:cursor-not-allowed flex items-center justify-center transition-colors"
            >
              {isLoading ? "Saving..." : "Save Changes"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

const Users = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(9);
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [totalUsers, setTotalUsers] = useState(0);

  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [isUpdating, setIsUpdating] = useState(false);
  const [updateError, setUpdateError] = useState<string | null>(null);

  const getAuthToken = () => {
    return localStorage.getItem("authToken");
  };

  const fetchUsers = async (page: number, pageSize: number) => {
    setLoading(true);
    setError(null);
    try {
      const token = getAuthToken();
      if (!token)
        throw new Error("No authentication token found. Please log in.");

      const url = new URL(`${API_BASE_URL}/admin/users`);
      url.searchParams.append("page", page.toString());
      url.searchParams.append("page_size", pageSize.toString());

      const response = await fetch(url.toString(), {
        method: "GET",
        headers: { Authorization: `Bearer ${token}` },
      });

      if (!response.ok) {
        if (response.status === 403)
          throw new Error(
            "Forbidden: You do not have admin rights to view this page.",
          );
        if (response.status === 401)
          throw new Error(
            "Unauthorized: Your session has expired. Please log in again.",
          );
        throw new Error(`Failed to fetch users: ${response.statusText}`);
      }

      const data: ApiUser[] = await response.json();

      const startingId = (page - 1) * pageSize + 1;
      const transformedUsers: User[] = data.map((apiUser, index) => ({
        displayId: startingId + index,
        id: apiUser.id,
        name: apiUser.username,
        email: apiUser.email,
        role: apiUser.is_admin ? "Admin" : "User",
        lastLogin: "N/A",
      }));

      setUsers(transformedUsers);

      const newTotal = (page - 1) * pageSize + data.length;
      if (data.length === pageSize) {
        setTotalUsers(newTotal + 1);
      } else {
        setTotalUsers(newTotal);
      }
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "Failed to fetch users");
      setUsers([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers(currentPage, rowsPerPage);
  }, [currentPage, rowsPerPage]);

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

  const handleEdit = (id: string) => {
    const userToEdit = users.find((u) => u.id === id);
    if (userToEdit) {
      setEditingUser(userToEdit);
      setUpdateError(null);
      setIsEditModalOpen(true);
    }
  };

  const handleUpdateUser = async (id: string, data: EditUserFormData) => {
    if (!editingUser) {
      setUpdateError("Cannot update: Original user data not found.");
      return;
    }

    if (
      data.username === editingUser.name &&
      data.email === editingUser.email
    ) {
      setIsEditModalOpen(false);
      return;
    }

    setIsUpdating(true);
    setUpdateError(null);

    try {
      const token = getAuthToken();
      if (!token) throw new Error("Authentication token not found.");

      const payload: UpdateUserPayload = {
        username: data.username,
        email: data.email,
      };

      const response = await fetch(`${API_BASE_URL}/admin/users/${id}`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || `Failed to update user.`);
      }

      setIsEditModalOpen(false);
      setEditingUser(null);
      fetchUsers(currentPage, rowsPerPage);
    } catch (err: unknown) {
      setUpdateError(
        err instanceof Error ? err.message : "Failed to update user",
      );
    } finally {
      setIsUpdating(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (
      !window.confirm("Are you sure you want to permanently delete this user?")
    )
      return;

    try {
      const token = getAuthToken();
      if (!token) throw new Error("Authentication token not found.");

      const response = await fetch(`${API_BASE_URL}/admin/users/${id}`, {
        method: "DELETE",
        headers: { Authorization: `Bearer ${token}` },
      });

      if (response.status === 204) {
        if (users.length === 1 && currentPage > 1) {
          setCurrentPage(currentPage - 1);
        } else {
          fetchUsers(currentPage, rowsPerPage);
        }
      } else {
        const errorData = await response.json().catch(() => null);
        throw new Error(errorData?.message || "Failed to delete user.");
      }
    } catch (err: unknown) {
      alert(
        `Error: ${err instanceof Error ? err.message : "Failed to delete user"}`,
      );
    }
  };

  const totalPages = Math.ceil(totalUsers / rowsPerPage);

  const goToPage = (page: number) => {
    if (page >= 1 && page <= totalPages) {
      setCurrentPage(page);
    }
  };

  const getPageNumbers = () => {
    const pageNumbers: (number | string)[] = [];
    if (totalPages <= 1) return [];

    if (totalPages <= 5) {
      return Array.from({ length: totalPages }, (_, i) => i + 1);
    }
    pageNumbers.push(1);
    if (currentPage > 3) pageNumbers.push("...");
    if (currentPage > 2) pageNumbers.push(currentPage - 1);
    if (currentPage > 1 && currentPage < totalPages)
      pageNumbers.push(currentPage);
    if (currentPage < totalPages - 1) pageNumbers.push(currentPage + 1);
    if (currentPage < totalPages - 2) pageNumbers.push("...");
    pageNumbers.push(totalPages);

    return [...new Set(pageNumbers)];
  };

  return (
    <>
      <div className="userBody min-h-screen bg-gray-100 dark:bg-[#0D1117]">
        <Navbar />
        <main className="user-main-content flex-grow">
          <div className="usersDisp max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-6">
              User Management
            </h1>

            {loading && (
              <div className="text-center py-10">
                <span className="animate-spin inline-block h-8 w-8 border-b-2 border-[#0F5BA7] rounded-full"></span>
              </div>
            )}

            {error && (
              <div
                className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4"
                role="alert"
              >
                <strong>Error:</strong> {error}
                <button
                  onClick={() => fetchUsers(currentPage, rowsPerPage)}
                  className="ml-4 bg-red-500 text-white px-3 py-1 rounded text-sm hover:bg-red-600"
                >
                  Retry
                </button>
              </div>
            )}

            {!loading &&
              !error &&
              (users.length > 0 ? (
                <>
                  <UsersTable
                    users={users}
                    onEdit={handleEdit}
                    onDelete={handleDelete}
                  />
                  <div className="usersPaging flex justify-center items-center py-4 gap-2 mt-4">
                    <button
                      onClick={() => goToPage(currentPage - 1)}
                      disabled={currentPage === 1}
                      className="px-3 py-2 bg-white dark:bg-[#21262D] border border-gray-300 dark:border-[#30363D] text-gray-700 dark:text-gray-300 rounded-md hover:bg-gray-50 dark:hover:bg-[#30363D] disabled:opacity-50 transition-colors"
                    >
                      Prev
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
                          className={`px-4 py-2 rounded-md transition-colors ${currentPage === page ? "bg-[#0F5BA7] text-white" : "bg-white dark:bg-[#21262D] border border-gray-300 dark:border-[#30363D] text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-[#30363D]"}`}
                        >
                          {page}
                        </button>
                      ),
                    )}
                    <button
                      onClick={() => goToPage(currentPage + 1)}
                      disabled={
                        currentPage === totalPages || users.length < rowsPerPage
                      }
                      className="px-3 py-2 bg-white dark:bg-[#21262D] border border-gray-300 dark:border-[#30363D] text-gray-700 dark:text-gray-300 rounded-md hover:bg-gray-50 dark:hover:bg-[#30363D] disabled:opacity-50 transition-colors"
                    >
                      Next
                    </button>
                  </div>
                </>
              ) : (
                <div className="text-center py-10 bg-white dark:bg-[#161B22] rounded-lg shadow-sm">
                  <p className="text-gray-500 text-lg">No users found.</p>
                </div>
              ))}
          </div>
        </main>
        <Footer />
        <HelpMenu />
      </div>

      <EditUserModal
        isOpen={isEditModalOpen}
        onClose={() => setIsEditModalOpen(false)}
        onSubmit={handleUpdateUser}
        user={editingUser}
        isLoading={isUpdating}
        error={updateError}
      />
    </>
  );
};

export default Users;
