import { Outlet, Link, useLocation } from 'react-router-dom';
import { Flame, MessageCircle, User as UserIcon, LogOut } from 'lucide-react';

export default function MainLayout({ setToken }) {
  const location = useLocation();

  const handleLogout = async () => {
    try {
      await fetch('http://localhost:8000/api/logout', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Accept': 'application/json'
        }
      });
    } catch(e) {}
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setToken(null);
  };

  const navItems = [
    { path: '/', icon: Flame, label: 'Discover' },
    { path: '/matches', icon: MessageCircle, label: 'Matches' }
  ];

  return (
    <div className="app-container">
      {/* Sidebar / Navbar */}
      <nav className="glass-panel flex-col justify-between" style={{ width: '80px', height: '100%', borderRadius: '0', borderRight: '1px solid var(--border-color)', padding: '2rem 0', alignItems: 'center' }}>
        <div className="flex-col gap-4">
          <div style={{ background: 'var(--accent-gradient)', padding: '10px', borderRadius: '12px', marginBottom: '2rem' }}>
            <Flame size={28} color="white" />
          </div>
          
          {navItems.map(item => {
            const Icon = item.icon;
            const isActive = location.pathname === item.path;
            return (
              <Link key={item.path} to={item.path} style={{ position: 'relative', display: 'flex', justifyContent: 'center', padding: '12px', borderRadius: '12px', background: isActive ? 'rgba(236, 72, 153, 0.1)' : 'transparent', color: isActive ? 'var(--accent-primary)' : 'var(--text-secondary)', transition: 'all 0.2s' }}>
                <Icon size={24} />
              </Link>
            )
          })}
        </div>

        <div className="flex-col gap-4">
          <button style={{ background: 'none', border: 'none', color: 'var(--text-secondary)', cursor: 'pointer', padding: '12px' }}>
            <UserIcon size={24} />
          </button>
          <button onClick={handleLogout} style={{ background: 'none', border: 'none', color: 'var(--danger)', cursor: 'pointer', padding: '12px' }}>
            <LogOut size={24} />
          </button>
        </div>
      </nav>

      {/* Main Content Area */}
      <main style={{ flex: 1, position: 'relative', overflow: 'hidden' }}>
        <Outlet />
      </main>
    </div>
  );
}
