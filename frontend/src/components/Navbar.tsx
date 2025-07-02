import "./Navbar.css";
import { FaCircleUser } from "react-icons/fa6";
import { IoIosLogOut } from "react-icons/io";
import { useState } from "react";
import { GiHamburgerMenu } from "react-icons/gi";
import { IoClose } from "react-icons/io5";

function Navbar() {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  const toggleMobileMenu = () => {
    setIsMobileMenuOpen(!isMobileMenuOpen);
  };

  return (
    <nav className="navbar">
      <div className="navbar-left">
        <img src="/src/assets/logo.png" alt="Logo" className="logo-image" />
        <div className="logo">Swift Signals</div>
      </div>

      <button className="mobile-menu-toggle" onClick={toggleMobileMenu}>
        {isMobileMenuOpen ? (
          <IoClose size={30} />
        ) : (
          <GiHamburgerMenu size={30} />
        )}
      </button>

      <div className={`navbar-center ${isMobileMenuOpen ? "active" : ""}`}>
        <ul className="nav-links">
          <li>
            <a href="/dashboard" onClick={toggleMobileMenu}>
              Dashboard
            </a>
          </li>
          <li>
            <a href="/intersections" onClick={toggleMobileMenu}>
              Intersections
            </a>
          </li>
          <li>
            <a href="/simulations" onClick={toggleMobileMenu}>
              Simulations
            </a>
          </li>
          <li>
            <a href="/users" onClick={toggleMobileMenu}>
              Users
            </a>
          </li>
        </ul>
        <div className="mobile-user-profile">
          <FaCircleUser size={45} />
          <span>John Doe</span>
          <a href="/" className="logout-icon" onClick={toggleMobileMenu}>
            <IoIosLogOut size={35} />
          </a>
        </div>
      </div>

      <div className="navbar-right">
        <div className="user-profile">
          <FaCircleUser size={45} />
          <span>John Doe</span>
          <a href="/" className="logout-icon">
            <IoIosLogOut size={35} />
          </a>
        </div>
      </div>
    </nav>
  );
}

export default Navbar;
