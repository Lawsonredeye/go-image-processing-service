import React from 'react';
import './Header.css';

const Header: React.FC = () => {
  return (
    <div className="header">
      <nav className="nav">
        <div className="nav-brand">GIPS</div>

        {/* --- Desktop Navigation --- */}
        <div className="desktop-nav">
          <a href="#">Compress</a>
          <a href="#">Resize</a>
          <a href="#">Convert</a>
        </div>

        {/* --- Mobile Navigation --- */}
        <div className="mobile-nav">
          <details className="hamburger-menu">
            <summary>
              <svg
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <line x1="3" y1="12" x2="21" y2="12"></line>
                <line x1="3" y1="6" x2="21" y2="6"></line>
                <line x1="3" y1="18" x2="21" y2="18"></line>
              </svg>
            </summary>
            <ul>
              <li><a href="#">Compress</a></li>
              <li><a href="#">Resize</a></li>
              <li><a href="#">Convert</a></li>
            </ul>
          </details>
        </div>
      </nav>
    </div>
  );
};

export default Header;