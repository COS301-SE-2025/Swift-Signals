import "./Footer.css";
import ThemeToggle from "./ThemeToggle";

const Footer: React.FC = () => {
  return (
    <footer className="footer bg-gray-800 text-white py-3">
      <div className="container mx-auto flex items-center justify-between px-4">
        {/* Logo Section */}
        <div className="footer-logo flex items-center">
          <img src="/src/assets/scs-logo.png" alt="Logo" className="h-6 mr-2" />
          <span className="footerText text-sm font-medium">
            A Southern Cross Solutions Product
          </span>
        </div>

        <div className="footer-copyright text-xs opacity-75">
          © {new Date().getFullYear()} Swift Signals. All rights reserved.
        </div>

        <div className="footer-toggle">
          <ThemeToggle />
        </div>
      </div>
    </footer>
  );
};

export default Footer;
