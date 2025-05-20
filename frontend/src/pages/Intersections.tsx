import { useState } from 'react';
import Navbar from '../components/Navbar';

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
  const [totalPages] = useState(1);

  // Sample user data
  const users: User[] = [
    { id: 1, name: 'John Doe', email: 'email@email.com', role: 'Admin', lastLogin: '2025-05-13 09:00' },
    { id: 2, name: 'Jane Smith', email: 'email@email.com', role: 'Engineer', lastLogin: '2025-05-13 09:00' },
    { id: 3, name: 'John Calvin', email: 'email@email.com', role: 'Viewer', lastLogin: '2025-05-13 09:00' },
    { id: 4, name: 'Paul Washer', email: 'email@email.com', role: 'Viewer', lastLogin: '2025-05-13 09:00' },
    { id: 5, name: 'Joshua Garner', email: 'email@email.com', role: 'Viewer', lastLogin: '2025-05-13 09:00' },
    { id: 6, name: 'Chris Xides', email: 'email@email.com', role: 'Engineer', lastLogin: '2025-05-13 09:00' },
    { id: 7, name: 'Kgosi Segale', email: 'email@email.com', role: 'Viewer', lastLogin: '2025-05-13 09:00' },
    { id: 8, name: 'John Flavel', email: 'email@email.com', role: 'Viewer', lastLogin: '2025-05-13 09:00' },
    { id: 9, name: 'John Owen', email: 'email@email.com', role: 'Viewer', lastLogin: '2025-05-13 09:00' },
  ];

  const handleEdit = (id: number) => console.log(`Edit user ${id}`);
  const handleDelete = (id: number) => console.log(`Delete user ${id}`);

  const goToNextPage = () => {
    if (currentPage < totalPages) setCurrentPage(currentPage + 1);
  };

  const goToPreviousPage = () => {
    if (currentPage > 1) setCurrentPage(currentPage - 1);
  };

  return (
    <div className="min-h-screen bg-gray-100">
      <Navbar />
      <div className="max-w-6xl mx-auto px-4 py-8">
        <div className="bg-white rounded-lg shadow-sm overflow-hidden">
          <table className="w-full border-collapse">
            <thead className="text-left border-b">
              <tr>
                <th className="px-4 py-3 font-bold">ID</th>
                <th className="px-4 py-3 font-bold">Name</th>
                <th className="px-4 py-3 font-bold">Email</th>
                <th className="px-4 py-3 font-bold">Role</th>
                <th className="px-4 py-3 font-bold">Last Login</th>
                <th className="px-4 py-3 font-bold">Actions</th>
              </tr>
            </thead>
            <tbody>
              {users.map((user) => (
                <tr key={user.id} className="border-b hover:bg-gray-50">
                  <td className="px-4 py-3">{user.id}</td>
                  <td className="px-4 py-3">{user.name}</td>
                  <td className="px-4 py-3">{user.email}</td>
                  <td className="px-4 py-3">{user.role}</td>
                  <td className="px-4 py-3">{user.lastLogin}</td>
                  <td className="px-4 py-3">
                    <div className="flex gap-2 justify-center">
                      <button
                        onClick={() => handleEdit(user.id)}
                        className="p-2 bg-green-500 text-white rounded-full flex items-center justify-center"
                        aria-label="Edit user"
                      >
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                          strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                          <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
                        </svg>
                      </button>
                      <button
                        onClick={() => handleDelete(user.id)}
                        className="p-2 bg-red-500 text-white rounded-full flex items-center justify-center"
                        aria-label="Delete user"
                      >
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                          strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                          <polyline points="3 6 5 6 21 6" />
                          <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
                        </svg>
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>

          <div className="flex justify-center items-center py-4 gap-4">
            <button
              onClick={goToPreviousPage}
              disabled={currentPage === 1}
              className="px-3 py-1 border rounded-md hover:bg-gray-100 disabled:opacity-50"
            >
              &lt;
            </button>
            <span className="text-gray-700">Page {currentPage} of {totalPages}</span>
            <button
              onClick={goToNextPage}
              disabled={currentPage === totalPages}
              className="px-3 py-1 border rounded-md hover:bg-gray-100 disabled:opacity-50"
            >
              &gt;
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Users;
