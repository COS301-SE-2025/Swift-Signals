import { PencilLine, Trash2 } from "lucide-react";
import React from "react";
import "../styles/UsersTable.css";

interface User {
  displayId: number;
  id: string;
  name: string;
  email: string;
  role: string;
  lastLogin: string;
}

interface UsersTableProps {
  users: User[];
  onEdit: (id: string) => void;
  onDelete: (id: string) => void;
}

const UsersTable: React.FC<UsersTableProps> = ({ users, onEdit, onDelete }) => {
  return (
    <div className="usersTablePage bg-white dark:bg-[#161B22] rounded-lg shadow-sm overflow-hidden">
      <table className="usersTable w-full border-collapse">
        <thead className="border-b bg-gray-50 dark:bg-[#21262D] dark:border-[#30363D]">
          <tr>
            <th className="px-4 py-3 font-bold text-black text-center">ID</th>
            <th className="px-4 py-3 font-bold text-black text-center">Name</th>
            <th className="px-4 py-3 font-bold text-black text-center">
              Email
            </th>
            <th className="px-4 py-3 font-bold text-black text-center">Role</th>
            <th className="px-4 py-3 font-bold text-black text-center">
              Last Login
            </th>
            <th className="px-4 py-3 font-bold text-black text-center">
              Actions
            </th>
          </tr>
        </thead>
        <tbody>
          {users.map((user) => (
            <tr
              key={user.id}
              className="border-b hover:bg-gray-50 dark:border-[#30363D]"
            >
              {/* MODIFIED: Display the simple number ID */}
              <td className="px-4 py-3 text-black text-center">
                {user.displayId}
              </td>
              <td className="px-4 py-3 text-black">{user.name}</td>
              <td className="px-4 py-3 text-black">{user.email}</td>
              <td className="px-4 py-3 text-black text-center">{user.role}</td>
              <td className="px-4 py-3 text-black text-center">
                {user.lastLogin}
              </td>
              <td className="px-4 py-3">
                <div className="flex gap-2 justify-center">
                  <button
                    onClick={() => onEdit(user.id)}
                    className="editUser p-2 bg-[#2B9348] dark:bg-[#2DA44E] text-white rounded-full flex items-center justify-center hover:bg-green-600 transition-colors"
                    aria-label="Edit user"
                  >
                    <PencilLine size={18} strokeWidth={2} />
                  </button>
                  <button
                    onClick={() => onDelete(user.id)}
                    className="deleteUser p-2 bg-red-500 dark:bg-[#DA3633] text-white rounded-full flex items-center justify-center hover:bg-red-600 transition-colors"
                    aria-label="Delete user"
                  >
                    <Trash2 size={18} strokeWidth={2} />
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
