import { useState, useEffect, useRef } from 'react';
import { Heart, X, Sparkles, Star, Info, RotateCcw } from 'lucide-react';

export default function SwipePage() {
  const [profiles, setProfiles] = useState([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [loading, setLoading] = useState(true);
  const [matchAlert, setMatchAlert] = useState(null); // null | { name, photo }
  const [swipeDirection, setSwipeDirection] = useState(null); // 'left' | 'right'
  const [showDetail, setShowDetail] = useState(false);

  // Drag state
  const [isDragging, setIsDragging] = useState(false);
  const [dragX, setDragX] = useState(0);
  const [dragY, setDragY] = useState(0);
  const dragStart = useRef(null);
  const cardRef = useRef(null);

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
      setCurrentIndex(0);
    } catch(e) {
      console.error(e);
    } finally {
      setLoading(false);
    }
  };

  const handleSwipe = async (type) => {
    if (currentIndex >= profiles.length) return;
    const profile = profiles[currentIndex];

    // Animate out
    setSwipeDirection(type === 'like' ? 'right' : 'left');
    setTimeout(async () => {
      setCurrentIndex(prev => prev + 1);
      setSwipeDirection(null);
      setDragX(0);
      setDragY(0);
      setShowDetail(false);
    }, 350);

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
        setTimeout(() => {
          setMatchAlert({
            name: profile.name,
            photo: profile.profile?.photo_url || `https://ui-avatars.com/api/?name=${profile.name}&background=ec4899&color=fff`
          });
          setTimeout(() => setMatchAlert(null), 4000);
        }, 400);
      }
    } catch (e) {
      console.error(e);
    }
  };

  // Drag handlers
  const onPointerDown = (e) => {
    if (e.target.closest('button')) return;
    dragStart.current = { x: e.clientX, y: e.clientY };
    setIsDragging(true);
    e.currentTarget.setPointerCapture(e.pointerId);
  };

  const onPointerMove = (e) => {
    if (!isDragging || !dragStart.current) return;
    const dx = e.clientX - dragStart.current.x;
    const dy = e.clientY - dragStart.current.y;
    setDragX(dx);
    setDragY(dy);
  };

  const onPointerUp = () => {
    if (!isDragging) return;
    setIsDragging(false);
    const threshold = 100;
    if (dragX > threshold) {
      handleSwipe('like');
    } else if (dragX < -threshold) {
      handleSwipe('pass');
    } else {
      // Snap back
      setDragX(0);
      setDragY(0);
    }
    dragStart.current = null;
  };

  const getAge = (birthdate) => {
    if (!birthdate) return null;
    return new Date().getFullYear() - new Date(birthdate).getFullYear();
  };

  const rotation = dragX / 15;
  const likeOpacity = Math.min(dragX / 80, 1);
  const passOpacity = Math.min(-dragX / 80, 1);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full w-full">
        <div style={{ textAlign: 'center' }}>
          <div style={{ width: 60, height: 60, borderRadius: '50%', background: 'var(--accent-gradient)', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 1rem', animation: 'pulse 1.5s ease-in-out infinite' }}>
            <Heart size={28} color="white" />
          </div>
          <p style={{ color: 'var(--text-secondary)' }}>Finding people near you...</p>
        </div>
      </div>
    );
  }

  const currentProfile = profiles[currentIndex];

  return (
    <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', height: '100%', width: '100%', position: 'relative', userSelect: 'none' }}>

      {/* Match Alert */}
      {matchAlert && (
        <div style={{ position: 'absolute', inset: 0, zIndex: 50, background: 'rgba(15,23,42,0.85)', backdropFilter: 'blur(16px)', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center' }}>
          <div className="animate-slide-up" style={{ textAlign: 'center' }}>
            {/* Avatar pair */}
            <div style={{ display: 'flex', justifyContent: 'center', gap: '-16px', marginBottom: '1.5rem' }}>
              <div style={{ width: 90, height: 90, borderRadius: '50%', overflow: 'hidden', border: '3px solid var(--accent-primary)', boxShadow: '0 0 24px rgba(236,72,153,0.6)', zIndex: 1 }}>
                <img src={matchAlert.photo} style={{ width: '100%', height: '100%', objectFit: 'cover' }} alt="" />
              </div>
              <div style={{ width: 40, height: 40, borderRadius: '50%', background: 'var(--accent-gradient)', display: 'flex', alignItems: 'center', justifyContent: 'center', zIndex: 2, margin: 'auto -12px', boxShadow: '0 0 16px rgba(236,72,153,0.8)' }}>
                <Heart size={18} color="white" fill="white" />
              </div>
              <div style={{ width: 90, height: 90, borderRadius: '50%', overflow: 'hidden', border: '3px solid var(--accent-secondary)', boxShadow: '0 0 24px rgba(139,92,246,0.6)', zIndex: 1 }}>
                <img src={JSON.parse(localStorage.getItem('user'))?.profile?.photo_url || `https://ui-avatars.com/api/?name=${JSON.parse(localStorage.getItem('user'))?.name}&background=8b5cf6&color=fff`} style={{ width: '100%', height: '100%', objectFit: 'cover' }} alt="" />
              </div>
            </div>
            <Sparkles size={48} color="#ec4899" style={{ marginBottom: '0.75rem' }} />
            <h1 style={{ fontSize: '2.5rem', background: 'var(--accent-gradient)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent', marginBottom: '0.5rem' }}>It's a Match!</h1>
            <p style={{ color: 'var(--text-secondary)', fontSize: '1rem' }}>You and <strong style={{ color: 'white' }}>{matchAlert.name}</strong> liked each other</p>
            <div style={{ display: 'flex', gap: '1rem', justifyContent: 'center', marginTop: '2rem' }}>
              <button className="btn btn-primary" onClick={() => setMatchAlert(null)} style={{ minWidth: '140px' }}>
                Send Message
              </button>
              <button className="btn" onClick={() => setMatchAlert(null)} style={{ minWidth: '140px', background: 'var(--bg-secondary)', color: 'white', border: '1px solid var(--border-color)' }}>
                Keep Swiping
              </button>
            </div>
          </div>
        </div>
      )}

      {currentIndex < profiles.length ? (
        <div style={{ position: 'relative', width: '360px', paddingBottom: '90px' }}>
          {/* Next card preview (behind) */}
          {profiles[currentIndex + 1] && (
            <div style={{ position: 'absolute', top: 8, left: 8, right: 8, height: '560px', borderRadius: '20px', overflow: 'hidden', background: 'var(--bg-secondary)', opacity: 0.5, transform: 'scale(0.96)' }}>
              <img
                src={profiles[currentIndex + 1].profile?.photo_url || `https://ui-avatars.com/api/?name=${profiles[currentIndex + 1].name}&size=400`}
                style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                alt=""
              />
            </div>
          )}

          {/* Main Card */}
          <div
            ref={cardRef}
            onPointerDown={onPointerDown}
            onPointerMove={onPointerMove}
            onPointerUp={onPointerUp}
            onPointerLeave={onPointerUp}
            style={{
              position: 'relative',
              width: '100%',
              height: '560px',
              borderRadius: '20px',
              overflow: 'hidden',
              background: 'var(--bg-secondary)',
              cursor: isDragging ? 'grabbing' : 'grab',
              transform: swipeDirection === 'right'
                ? 'translateX(120%) rotate(20deg)'
                : swipeDirection === 'left'
                ? 'translateX(-120%) rotate(-20deg)'
                : `translateX(${dragX}px) translateY(${dragY * 0.2}px) rotate(${rotation}deg)`,
              transition: isDragging ? 'none' : swipeDirection ? 'transform 0.35s ease-in' : 'transform 0.4s cubic-bezier(0.34, 1.56, 0.64, 1)',
              boxShadow: isDragging ? '0 20px 60px rgba(0,0,0,0.5)' : '0 8px 30px rgba(0,0,0,0.3)',
              touchAction: 'none'
            }}
          >
            <img
              src={currentProfile.profile?.photo_url || `https://ui-avatars.com/api/?name=${currentProfile.name}&size=400&background=ec4899&color=fff`}
              alt={currentProfile.name}
              style={{ width: '100%', height: showDetail ? '45%' : '70%', objectFit: 'cover', transition: 'height 0.3s ease', pointerEvents: 'none' }}
            />

            {/* Like/Pass Overlays */}
            <div style={{ position: 'absolute', top: 24, left: 20, transform: 'rotate(-15deg)', opacity: likeOpacity, transition: isDragging ? 'none' : 'opacity 0.2s', pointerEvents: 'none' }}>
              <div style={{ border: '3px solid var(--success)', borderRadius: '8px', padding: '6px 16px', color: 'var(--success)', fontWeight: '800', fontSize: '1.5rem', letterSpacing: '2px', textShadow: '0 0 10px var(--success)' }}>LIKE</div>
            </div>
            <div style={{ position: 'absolute', top: 24, right: 20, transform: 'rotate(15deg)', opacity: passOpacity, transition: isDragging ? 'none' : 'opacity 0.2s', pointerEvents: 'none' }}>
              <div style={{ border: '3px solid var(--danger)', borderRadius: '8px', padding: '6px 16px', color: 'var(--danger)', fontWeight: '800', fontSize: '1.5rem', letterSpacing: '2px', textShadow: '0 0 10px var(--danger)' }}>NOPE</div>
            </div>

            {/* Profile Info */}
            <div style={{ background: 'linear-gradient(to top, rgba(15,23,42,1) 60%, transparent)', padding: '1.25rem 1.5rem', position: 'absolute', bottom: 0, left: 0, right: 0 }}>
              <div style={{ display: 'flex', alignItems: 'flex-end', gap: '8px' }}>
                <h2 style={{ fontSize: '1.5rem', margin: 0 }}>{currentProfile.name}</h2>
                {getAge(currentProfile.profile?.birthdate) && (
                  <span style={{ fontSize: '1.2rem', color: 'var(--text-secondary)', marginBottom: '2px' }}>{getAge(currentProfile.profile.birthdate)}</span>
                )}
              </div>

              {showDetail && (
                <div className="animate-slide-up">
                  {currentProfile.profile?.bio && (
                    <p style={{ color: 'var(--text-secondary)', marginTop: '0.75rem', fontSize: '0.875rem', lineHeight: 1.5 }}>
                      {currentProfile.profile.bio}
                    </p>
                  )}
                  {currentProfile.profile?.gender && (
                    <div style={{ display: 'flex', gap: '8px', marginTop: '0.75rem', flexWrap: 'wrap' }}>
                      <span style={{ background: 'rgba(236,72,153,0.15)', color: 'var(--accent-primary)', padding: '4px 12px', borderRadius: '999px', fontSize: '0.78rem', fontWeight: '600' }}>
                        {currentProfile.profile.gender}
                      </span>
                    </div>
                  )}
                </div>
              )}

              {!showDetail && currentProfile.profile?.bio && (
                <p style={{ color: 'var(--text-secondary)', marginTop: '0.5rem', fontSize: '0.85rem', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                  {currentProfile.profile.bio}
                </p>
              )}

              <button
                onClick={() => setShowDetail(!showDetail)}
                style={{ background: 'none', border: 'none', color: 'var(--text-secondary)', cursor: 'pointer', marginTop: '0.5rem', display: 'flex', alignItems: 'center', gap: '4px', fontSize: '0.8rem', padding: 0 }}
              >
                <Info size={14} />
                {showDetail ? 'Less info' : 'More info'}
              </button>
            </div>
          </div>

          {/* Action Buttons */}
          <div style={{ position: 'absolute', bottom: '0', left: 0, right: 0, display: 'flex', justifyContent: 'center', alignItems: 'center', gap: '1.5rem' }}>
            <button
              className="btn-icon btn-pass"
              onClick={() => handleSwipe('pass')}
              style={{ width: 56, height: 56 }}
              title="Pass"
            >
              <X size={26} />
            </button>
            <button
              className="btn-icon"
              onClick={fetchProfiles}
              style={{ width: 44, height: 44, background: 'rgba(255,255,255,0.05)' }}
              title="Reload"
            >
              <RotateCcw size={16} />
            </button>
            <button
              className="btn-icon btn-like"
              onClick={() => handleSwipe('like')}
              style={{ width: 56, height: 56 }}
              title="Like"
            >
              <Heart size={26} />
            </button>
          </div>

          {/* Swipe Hint */}
          <p style={{ textAlign: 'center', color: 'var(--text-secondary)', fontSize: '0.75rem', position: 'absolute', bottom: '-1.75rem', left: 0, right: 0 }}>
            Swipe or tap buttons
          </p>
        </div>
      ) : (
        <div className="flex-col items-center justify-center animate-slide-up text-center">
          <div style={{ width: 90, height: 90, borderRadius: '50%', background: 'rgba(255,255,255,0.04)', display: 'flex', alignItems: 'center', justifyContent: 'center', marginBottom: '1.25rem', border: '2px dashed var(--border-color)' }}>
            <Sparkles size={40} color="var(--text-secondary)" />
          </div>
          <h2 style={{ color: 'var(--text-primary)' }}>You've seen everyone!</h2>
          <p style={{ color: 'var(--text-secondary)', margin: '0.5rem 0 1.5rem' }}>Check back later for new profiles</p>
          <button className="btn btn-primary" onClick={fetchProfiles}>
            <RotateCcw size={16} />
            &nbsp; Refresh Profiles
          </button>
        </div>
      )}
    </div>
  );
}
