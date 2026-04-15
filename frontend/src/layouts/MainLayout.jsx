import { Outlet, Link, useLocation } from 'react-router-dom';
import { Flame, MessageCircle, User as UserIcon, LogOut, Settings } from 'lucide-react';

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
    { path: '/matches', icon: MessageCircle, label: 'Messages' },
    { path: '/profile', icon: UserIcon, label: 'Profile' },
    { path: '/settings', icon: Settings, label: 'Settings' }
  ];

  return (
    <div className="app-container">
      {/* Sidebar Navbar */}
      <nav
        className="glass-panel flex-col justify-between"
        style={{
          width: '80px',
          height: '100%',
          borderRadius: '0',
          borderRight: '1px solid var(--border-color)',
          padding: '1.5rem 0',
          alignItems: 'center',
          display: 'flex'
        }}
      >
        {/* Top Section: Logo + Nav */}
        <div className="flex-col gap-4" style={{ alignItems: 'center', width: '100%' }}>
          {/* Logo */}
          <div style={{ background: 'var(--accent-gradient)', padding: '10px', borderRadius: '12px', marginBottom: '1.5rem', display: 'flex', alignItems: 'center', justifyContent: 'center', boxShadow: '0 4px 14px rgba(236,72,153,0.4)' }}>
            <Flame size={26} color="white" />
          </div>

          {/* Nav Links */}
          {navItems.map(item => {
            const Icon = item.icon;
            const isActive = location.pathname === item.path;
            return (
              <Link
                key={item.path}
                to={item.path}
                title={item.label}
                style={{
                  position: 'relative',
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                  justifyContent: 'center',
                  gap: '4px',
                  padding: '12px 8px',
                  borderRadius: '12px',
                  width: '56px',
                  background: isActive ? 'rgba(236, 72, 153, 0.12)' : 'transparent',
                  color: isActive ? 'var(--accent-primary)' : 'var(--text-secondary)',
                  textDecoration: 'none',
                  transition: 'all 0.2s',
                  boxShadow: isActive ? 'inset 0 0 0 1px rgba(236,72,153,0.2)' : 'none'
                }}
              >
                <Icon size={22} />
                <span style={{ fontSize: '0.6rem', fontWeight: '600', letterSpacing: '0.03em', fontFamily: 'Outfit, sans-serif' }}>
                  {item.label}
                </span>
                {/* Active indicator */}
                {isActive && (
                  <div style={{ position: 'absolute', right: '-1px', top: '50%', transform: 'translateY(-50%)', width: '3px', height: '32px', background: 'var(--accent-gradient)', borderRadius: '3px 0 0 3px' }} />
                )}
              </Link>
            );
          })}
        </div>

        {/* Bottom Section: Logout */}
        <button
          onClick={handleLogout}
          title="Logout"
          style={{
            background: 'rgba(239,68,68,0.08)',
            border: '1px solid rgba(239,68,68,0.2)',
            color: 'var(--danger)',
            cursor: 'pointer',
            padding: '10px',
            borderRadius: '12px',
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            gap: '4px',
            width: '56px',
            transition: 'background 0.2s'
          }}
        >
          <LogOut size={20} />
          <span style={{ fontSize: '0.6rem', fontWeight: '600', letterSpacing: '0.03em', fontFamily: 'Outfit, sans-serif' }}>Logout</span>
        </button>
      </nav>

      {/* Main Content Area */}
      <main style={{ flex: 1, position: 'relative', overflow: 'hidden', display: 'flex', flexDirection: 'column' }}>
        <Outlet />
      </main>
    </div>
  );
}
