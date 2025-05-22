import React from 'react';

// TypeScript interface for user data
interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  lastLogin: string;
}

interface UsersTableProps {
  users: User[];
  onEdit: (id: number) => void;
  onDelete: (id: number) => void;
}

const UsersTable: React.FC<UsersTableProps> = ({ users, onEdit, onDelete }) => {
  return (
    <div className="bg-white rounded-lg shadow-sm overflow-hidden">
      <table className="w-full border-collapse">
        <thead className="border-b bg-gray-50">
          <tr>
            <th className="px-4 py-3 font-bold text-black text-center">ID</th>
            <th className="px-4 py-3 font-bold text-black text-center">Name</th>
            <th className="px-4 py-3 font-bold text-black text-center">Email</th>
            <th className="px-4 py-3 font-bold text-black text-center">Role</th>
            <th className="px-4 py-3 font-bold text-black text-center">Last Login</th>
            <th className="px-4 py-3 font-bold text-black text-center">Actions</th>
          </tr>
        </thead>
        <tbody>
          {users.map((user) => (
            <tr key={user.id} className="border-b hover:bg-gray-50">
              <td className="px-4 py-3 text-black text-center">{user.id}</td>
              <td className="px-4 py-3 text-black">{user.name}</td>
              <td className="px-4 py-3 text-black">{user.email}</td>
              <td className="px-4 py-3 text-black text-center">{user.role}</td>
              <td className="px-4 py-3 text-black text-center">{user.lastLogin}</td>
              <td className="px-4 py-3">
                <div className="flex gap-2 justify-center">
                  <button
                    onClick={() => onEdit(user.id)}
                    className="p-2 bg-green-500 text-white rounded-full flex items-center justify-center hover:bg-green-600 transition-colors"
                    aria-label="Edit user"
                  >
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                      strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                      <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
                    </svg>
                  </button>
                  <button
                    onClick={() => onDelete(user.id)}
                    className="p-2 bg-red-500 text-white rounded-full flex items-center justify-center hover:bg-red-600 transition-colors"
                    aria-label="Delete user"
                  >
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                      strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <polyline points="3 6 5 6 21 6" />
                      <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2 2h4a2 2 0 0 1 2 2v2" />
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default UsersTable;
