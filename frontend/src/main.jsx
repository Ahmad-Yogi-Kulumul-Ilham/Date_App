import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.jsx'

const themeColor = localStorage.getItem('theme-color');
if (themeColor) {
  document.documentElement.style.setProperty('--accent-primary', themeColor);
}

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
