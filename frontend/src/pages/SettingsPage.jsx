import { useState, useEffect } from 'react';
import { Settings as SettingsIcon, Palette, Save, CheckCircle } from 'lucide-react';

export default function SettingsPage() {
  const [accentColor, setAccentColor] = useState(localStorage.getItem('theme-color') || '#ec4899');
  const [saved, setSaved] = useState(false);

  const colors = [
    { name: 'Pink', hex: '#ec4899' },
    { name: 'Purple', hex: '#a855f7' },
    { name: 'Blue', hex: '#3b82f6' },
    { name: 'Green', hex: '#10b981' },
    { name: 'Orange', hex: '#f97316' },
    { name: 'Red', hex: '#ef4444' }
  ];

  const handleSave = () => {
    localStorage.setItem('theme-color', accentColor);
    document.documentElement.style.setProperty('--accent-primary', accentColor);
    setSaved(true);
    setTimeout(() => setSaved(false), 3000);
  };

  return (
    <div style={{ padding: '2rem', height: '100%', overflowY: 'auto' }}>
      <div style={{ maxWidth: '680px', margin: '0 auto' }}>
        <div className="animate-slide-up" style={{ marginBottom: '2rem' }}>
          <h1 style={{ fontSize: '2rem', margin: 0, display: 'flex', alignItems: 'center', gap: '12px' }}>
            <SettingsIcon size={32} style={{ color: 'var(--accent-primary)' }} />
            Settings
          </h1>
          <p style={{ color: 'var(--text-secondary)', marginTop: '0.25rem' }}>Customize your app experience</p>
        </div>

        {saved && (
          <div className="animate-slide-up" style={{
            background: 'rgba(16, 185, 129, 0.1)',
            border: '1px solid rgba(16,185,129,0.3)',
            color: 'var(--success)',
            padding: '1rem 1.25rem',
            borderRadius: '12px',
            marginBottom: '1.5rem',
            display: 'flex',
            alignItems: 'center',
            gap: '8px',
            fontSize: '0.875rem'
          }}>
            <CheckCircle size={18} />
            Settings saved successfully!
          </div>
        )}

        <div className="glass-panel animate-slide-up" style={{ padding: '2rem' }}>
          <h3 style={{ margin: '0 0 1.5rem 0', display: 'flex', alignItems: 'center', gap: '8px' }}>
            <Palette size={20} style={{ color: 'var(--accent-primary)' }} />
            Appearance
          </h3>

          <div className="input-group" style={{ marginBottom: '2rem' }}>
            <label className="input-label">Theme Color</label>
            <div style={{ display: 'flex', gap: '1rem', flexWrap: 'wrap', marginTop: '1rem' }}>
              {colors.map(c => (
                <button
                  key={c.name}
                  onClick={() => setAccentColor(c.hex)}
                  style={{
                    width: '48px',
                    height: '48px',
                    borderRadius: '50%',
                    background: c.hex,
                    border: accentColor === c.hex ? '3px solid white' : '3px solid transparent',
                    boxShadow: accentColor === c.hex ? `0 0 15px ${c.hex}80` : 'none',
                    cursor: 'pointer',
                    transition: 'all 0.2s',
                    position: 'relative'
                  }}
                  title={c.name}
                >
                  {accentColor === c.hex && (
                    <div style={{ position: 'absolute', inset: 0, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                      <CheckCircle size={20} color="white" />
                    </div>
                  )}
                </button>
              ))}
            </div>
          </div>

          <button onClick={handleSave} className="btn btn-primary" style={{ width: '100%', display: 'flex', justifyContent: 'center', gap: '8px' }}>
            <Save size={18} /> Save Settings
          </button>
        </div>
      </div>
    </div>
  );
}
