import React from 'react';
import './Header.css';

const Header: React.FC = () => {
  return (
    <div className="header">
      <nav className="nav">
        <div className="nav-brand">GIPS</div>
        <div className="nav-links">
          <a href="#">Compress</a>
          <a href="#">Resize</a>
          <a href="#">Convert</a>
        </div>
      </nav>
    </div>
  );
};

export default Header;
