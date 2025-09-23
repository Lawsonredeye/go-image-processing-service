import React from 'react';

const Footer: React.FC = () => {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="footer">
      <small>&copy; {currentYear} GIPS - Go Image Processing Service</small>
    </footer>
  );
};

export default Footer;
