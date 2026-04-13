import { useState, useEffect } from 'react';
import { Heart, X, Sparkles } from 'lucide-react';

export default function SwipePage() {
  const [profiles, setProfiles] = useState([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [loading, setLoading] = useState(true);
  const [matchAlert, setMatchAlert] = useState(false);

  useEffect(() => {
    fetchProfiles();
  }, []);

  const fetchProfiles = async () => {
    setLoading(true);
    try {
      const res = await fetch('http://localhost:8000/api/discover', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}`, 'Accept': 'application/json' }
      });
      const data = await res.json();
      setProfiles(data);
    } catch(e) {
      console.error(e);
    } finally {
      setLoading(false);
    }
  };

  const handleSwipe = async (type) => {
    if (currentIndex >= profiles.length) return;
    
    const profile = profiles[currentIndex];
    setCurrentIndex(prev => prev + 1);

    try {
      const res = await fetch('http://localhost:8000/api/swipe', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Accept': 'application/json'
        },
        body: JSON.stringify({ swipee_id: profile.id, type })
      });
      const data = await res.json();
      if (data.match) {
        setMatchAlert(true);
        setTimeout(() => setMatchAlert(false), 3000);
      }
    } catch (e) {
      console.error(e);
    }
  };

  if (loading) return <div className="flex items-center justify-center h-full w-full">Loading...</div>;

  return (
    <div className="flex-col items-center justify-center h-full w-full relative">
      {/* Match Alert Overlay */}
      {matchAlert && (
        <div style={{ position: 'absolute', inset: 0, zIndex: 50, background: 'rgba(236,72,153, 0.4)', backdropFilter: 'blur(10px)', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center' }}>
          <Sparkles size={64} className="animate-slide-up" color="white" style={{ marginBottom: '1rem' }} />
          <h1 className="animate-slide-up" style={{ fontSize: '3rem', color: 'white', textShadow: '0 4px 10px rgba(0,0,0,0.5)' }}>It's a Match!</h1>
        </div>
      )}

      {currentIndex < profiles.length ? (
        <div style={{ position: 'relative', width: '360px', height: '600px', paddingBottom: '100px' }}>
          {/* Card */}
          <div className="glass-panel" style={{ width: '100%', height: '100%', position: 'absolute', overflow: 'hidden', display: 'flex', flexDirection: 'column' }}>
            <img 
              src={profiles[currentIndex].profile?.photo_url || `https://ui-avatars.com/api/?name=${profiles[currentIndex].name}&size=400`} 
              alt={profiles[currentIndex].name}
              style={{ width: '100%', height: '70%', objectFit: 'cover' }}
            />
            <div style={{ padding: '1.5rem', flex: 1, background: 'linear-gradient(to top, var(--bg-secondary) 80%, transparent)' }}>
              <div style={{ display: 'flex', alignItems: 'flex-end', gap: '8px' }}>
                <h2 style={{ fontSize: '1.5rem', margin: 0 }}>{profiles[currentIndex].name}</h2>
                {profiles[currentIndex].profile?.birthdate && (
                   <span style={{ fontSize: '1.25rem', color: 'var(--text-secondary)' }}>
                     {new Date().getFullYear() - new Date(profiles[currentIndex].profile.birthdate).getFullYear()}
                   </span>
                )}
              </div>
              <p style={{ color: 'var(--text-secondary)', marginTop: '0.5rem', fontSize: '0.875rem' }}>
                {profiles[currentIndex].profile?.bio || 'No bio yet.'}
              </p>
            </div>
          </div>
          
          {/* Action Buttons */}
          <div style={{ position: 'absolute', bottom: '0', left: 0, right: 0, display: 'flex', justifyContent: 'center', gap: '2rem' }}>
            <button className="btn-icon btn-pass" onClick={() => handleSwipe('pass')}>
              <X size={32} />
            </button>
            <button className="btn-icon btn-like" onClick={() => handleSwipe('like')}>
              <Heart size={32} />
            </button>
          </div>
        </div>
      ) : (
        <div className="flex-col items-center justify-center animate-slide-up text-center">
          <div style={{ width: 80, height: 80, borderRadius: '50%', background: 'rgba(255,255,255,0.05)', display: 'flex', alignItems: 'center', justifyContent: 'center', marginBottom: '1rem' }}>
            <Sparkles size={40} color="var(--text-secondary)" />
          </div>
          <h2 style={{ color: 'var(--text-secondary)' }}>You've seen everyone!</h2>
          <button className="btn btn-primary" onClick={fetchProfiles} style={{ marginTop: '1.5rem' }}>Refresh Profiles</button>
        </div>
      )}
    </div>
  );
}
