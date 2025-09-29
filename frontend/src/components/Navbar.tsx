import { useState, useEffect, useContext } from "react";
import { FaCircleUser } from "react-icons/fa6";
import { GiHamburgerMenu } from "react-icons/gi";
import { IoIosLogOut } from "react-icons/io";
import { IoClose } from "react-icons/io5";
import { useLocation } from "react-router-dom";

import logo from "../../src/assets/logo.png";
// import { API_BASE_URL } from "../config";
import "../components/Navbar.css";
import { UserContext } from "../context/UserContext";

function Navbar() {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [username, setUsername] = useState("");
  const location = useLocation();
  const userContext = useContext(UserContext);

  const toggleMobileMenu = () => {
    setIsMobileMenuOpen(!isMobileMenuOpen);
  };

  interface IsActiveFn {
    (path: string): boolean;
  }
  const isActive: IsActiveFn = (path) => location.pathname === path;

  useEffect(() => {
    if (userContext?.user) {
      setUsername(userContext.user.username);
    }
  }, [userContext]);

  return (
    <nav className="navbar">
      <div className="navbar-left">
        <img src={logo} alt="Logo" className="logo-image" />
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
            <a
              href="/dashboard"
              className={isActive("/dashboard") ? "active" : ""}
              onClick={toggleMobileMenu}
            >
              Dashboard
            </a>
          </li>
          <li>
            <a
              href="/intersections"
              className={isActive("/intersections") ? "active" : ""}
              onClick={toggleMobileMenu}
            >
              Intersections
            </a>
          </li>
          <li>
            <a
              href="/simulations"
              className={isActive("/simulations") ? "active" : ""}
              onClick={toggleMobileMenu}
            >
              Simulations
            </a>
          </li>
          {userContext?.user?.role === "admin" && (
            <li>
              <a
                href="/users"
                className={isActive("/users") ? "active" : ""}
                onClick={toggleMobileMenu}
              >
                Users
              </a>
            </li>
          )}
        </ul>
        <div className="mobile-user-profile">
          <FaCircleUser size={45} />
          <span>{username || "Loading..."}</span>
          <a href="/" className="logout-icon" onClick={toggleMobileMenu}>
            <IoIosLogOut size={35} />
          </a>
        </div>
      </div>

      <div className="navbar-right">
        <div className="user-profile">
          <FaCircleUser size={45} />
          <span>{username || "Loading..."}</span>
          <a href="/" className="logout-icon">
            <IoIosLogOut size={35} />
          </a>
        </div>
      </div>
    </nav>
  );
}

export default Navbar;
