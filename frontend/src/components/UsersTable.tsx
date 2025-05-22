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


export default UsersTable;
