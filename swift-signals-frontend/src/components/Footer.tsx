import './Footer.css';

const Footer: React.FC = () => {
  return (
    <footer className="footer bg-gray-800 text-white py-6">
      <div className="container mx-auto flex items-center">
        {/* Logo Section */}
        <div className="footer-logo flex items-center">
          <img src="/src/assets/scs-logo.png" alt="Logo" className="h-10 mr-3" />
          <span className="text-lg font-semibold">Southern Cross Solutions</span>
        </div>

        {/* Copyright Section
        <div className="footer-copyright text-sm">
          © {new Date().getFullYear()} Swift Signals. All rights reserved.
        </div> */}
      </div>
    </footer>
  );
};

export default Footer;