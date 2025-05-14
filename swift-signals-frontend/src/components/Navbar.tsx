import './Navbar.css';

function Navbar() {
  return (
    <nav className="navbar">
      <div className="navbar-left">
        <img src="/src/assets/placeholder.png" alt="Logo" className="logo-image" />
        <div className="logo">Swift Signals</div>
      </div>
      <div className="navbar-center">
        <ul className="nav-links">
          <li><a href="#dashboard">Dashboard</a></li>
          <li><a href="#intersections">Intersections</a></li>
          <li><a href="#simulations">Simulations</a></li>
          <li><a href="#users">Users</a></li>
        </ul>
      </div>
      <div className="navbar-right">
        <div className="user-profile">
          <img src="/src/assets/userProfileIcon.png" alt="User" className="user-image" />
          <span>John Doe</span>
          <a href="#logout" className="logout-icon">
            <img src="/src/assets/logoutIcon.png" alt="Logout" />
          </a>
        </div>
      </div>
    </nav>
  );
}

export default Navbar;