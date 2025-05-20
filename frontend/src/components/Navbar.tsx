import './Navbar.css';
import { FaCircleUser } from "react-icons/fa6";
import { IoIosLogOut } from "react-icons/io";
//import ThemeToggle from './ThemeToggle';

function Navbar() {
  return (
    <nav className="navbar">
      <div className="navbar-left">
        <img src="/src/assets/logo.png" alt="Logo" className="logo-image" />
        <div className="logo">Swift Signals</div>
      </div>
      <div className="navbar-center">
        <ul className="nav-links">
          <li><a href="/dashboard">Dashboard</a></li>
          <li><a href="/intersections">Intersections</a></li>
          <li><a href="/simulations">Simulations</a></li>
          <li><a href="/users">Users</a></li>
        </ul>
      </div>
      <div className="navbar-right">
        {/* <ThemeToggle /> */}
        <div className="user-profile">
          <FaCircleUser size={45}/>
          <span>John Doe</span>
          <a href="/" className="logout-icon">
            <IoIosLogOut size={35} color='#991915 dark: color=#fff'/>
          </a>
        </div>
      </div>
    </nav>
  );
}

export default Navbar;