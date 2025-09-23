import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import './App.css';
import Layout from './components/Layout';
import CompressPage from './pages/CompressPage';
import ConvertPage from './pages/ConvertPage';

function App() {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<CompressPage />} />
          <Route path="/convert" element={<ConvertPage />} />
          {/* Add routes for other pages here later */}
        </Routes>
      </Layout>
    </Router>
  );
}

export default App;