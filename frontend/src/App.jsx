import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { useState, useEffect } from 'react';
import AuthPage from './pages/AuthPage';
import SwipePage from './pages/SwipePage';
import MatchesPage from './pages/MatchesPage';
import MainLayout from './layouts/MainLayout';

function App() {
  const [token, setToken] = useState(localStorage.getItem('token'));

  if (!token) {
    return <AuthPage setToken={setToken} />;
  }

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<MainLayout setToken={setToken} />}>
          <Route index element={<SwipePage />} />
          <Route path="matches" element={<MatchesPage />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;
