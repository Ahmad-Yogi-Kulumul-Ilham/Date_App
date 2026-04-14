import { useState, useEffect } from 'react';
import { User, Camera, Save, CheckCircle, AlertCircle, Calendar, Heart, Loader } from 'lucide-react';

export default function ProfilePage() {
  const [profile, setProfile] = useState(null);
  const [user, setUser] = useState(null);
  const [form, setForm] = useState({ bio: '', gender: '', interested_in: '', birthdate: '', photo_url: '' });
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [status, setStatus] = useState(null); // { type: 'success' | 'error', message: string }

  useEffect(() => {
    fetchProfile();
  }, []);

  const fetchProfile = async () => {
    setLoading(true);
    try {
      const res = await fetch('http://localhost:8000/api/profile', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}`, 'Accept': 'application/json' }
      });
      const data = await res.json();
      setUser(data);
      setProfile(data.profile);
      setForm({
        bio: data.profile?.bio || '',
        gender: data.profile?.gender || '',
        interested_in: data.profile?.interested_in || '',
        birthdate: data.profile?.birthdate || '',
        photo_url: data.profile?.photo_url || ''
      });
      // Update localStorage with fresh user data
      localStorage.setItem('user', JSON.stringify(data));
    } catch(e) {
      console.error(e);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async (e) => {
    e.preventDefault();
    setSaving(true);
    setStatus(null);
    try {
      const res = await fetch('http://localhost:8000/api/profile', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Accept': 'application/json'
        },
        body: JSON.stringify(form)
      });
      if (!res.ok) throw new Error('Failed to update profile');
      const data = await res.json();
      setProfile(data);
      setStatus({ type: 'success', message: 'Profile updated successfully!' });
      setTimeout(() => setStatus(null), 3000);
    } catch(e) {
      setStatus({ type: 'error', message: e.message || 'Something went wrong' });
    } finally {
      setSaving(false);
    }
  };

  const getAge = (birthdate) => {
    if (!birthdate) return null;
    return new Date().getFullYear() - new Date(birthdate).getFullYear();
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full w-full">
        <Loader size={32} style={{ color: 'var(--accent-primary)', animation: 'spin 1s linear infinite' }} />
      </div>
    );
  }

  const photoUrl = form.photo_url || `https://ui-avatars.com/api/?name=${user?.name || 'User'}&size=200&background=ec4899&color=fff`;

  return (
    <div style={{ height: '100%', overflowY: 'auto', padding: '2rem' }}>
      <div style={{ maxWidth: '680px', margin: '0 auto' }}>

        {/* Header */}
        <div className="animate-slide-up" style={{ marginBottom: '2rem' }}>
          <h1 style={{ fontSize: '2rem', margin: 0 }}>My Profile</h1>
          <p style={{ color: 'var(--text-secondary)', marginTop: '0.25rem' }}>Manage how others see you</p>
        </div>

        {/* Photo & Name Card */}
        <div className="glass-panel animate-slide-up" style={{ padding: '2rem', marginBottom: '1.5rem', display: 'flex', alignItems: 'center', gap: '2rem' }}>
          <div style={{ position: 'relative', flexShrink: 0 }}>
            <div style={{ width: '100px', height: '100px', borderRadius: '50%', overflow: 'hidden', border: '3px solid var(--accent-primary)', boxShadow: '0 0 20px rgba(236,72,153,0.4)' }}>
              <img src={photoUrl} alt="Profile" style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
            </div>
            <div style={{ position: 'absolute', bottom: 0, right: 0, background: 'var(--accent-gradient)', borderRadius: '50%', width: 28, height: 28, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
              <Camera size={14} color="white" />
            </div>
          </div>
          <div>
            <h2 style={{ margin: 0, fontSize: '1.5rem' }}>{user?.name}</h2>
            <p style={{ color: 'var(--text-secondary)', margin: '0.25rem 0', fontSize: '0.875rem' }}>{user?.email}</p>
            {getAge(form.birthdate) && (
              <div style={{ display: 'flex', alignItems: 'center', gap: '6px', marginTop: '0.5rem' }}>
                <span style={{ background: 'rgba(236,72,153,0.15)', color: 'var(--accent-primary)', padding: '4px 12px', borderRadius: '999px', fontSize: '0.8rem', fontWeight: '600' }}>
                  {getAge(form.birthdate)} years old
                </span>
                {form.gender && (
                  <span style={{ background: 'rgba(139,92,246,0.15)', color: 'var(--accent-secondary)', padding: '4px 12px', borderRadius: '999px', fontSize: '0.8rem', fontWeight: '600' }}>
                    {form.gender}
                  </span>
                )}
              </div>
            )}
          </div>
        </div>

        {/* Status message */}
        {status && (
          <div className="animate-slide-up" style={{
            background: status.type === 'success' ? 'rgba(16, 185, 129, 0.1)' : 'rgba(239, 68, 68, 0.1)',
            border: `1px solid ${status.type === 'success' ? 'rgba(16,185,129,0.3)' : 'rgba(239,68,68,0.3)'}`,
            color: status.type === 'success' ? 'var(--success)' : 'var(--danger)',
            padding: '1rem 1.25rem',
            borderRadius: '12px',
            marginBottom: '1.5rem',
            display: 'flex',
            alignItems: 'center',
            gap: '8px',
            fontSize: '0.875rem'
          }}>
            {status.type === 'success' ? <CheckCircle size={18} /> : <AlertCircle size={18} />}
            {status.message}
          </div>
        )}

        {/* Edit Form */}
        <div className="glass-panel animate-slide-up" style={{ padding: '2rem' }}>
          <h3 style={{ margin: '0 0 1.5rem 0', display: 'flex', alignItems: 'center', gap: '8px' }}>
            <User size={20} style={{ color: 'var(--accent-primary)' }} />
            Edit Profile
          </h3>

          <form onSubmit={handleSave} className="flex-col">
            {/* Photo URL */}
            <div className="input-group">
              <label className="input-label">
                <Camera size={14} style={{ display: 'inline', marginRight: '6px' }} />
                Photo URL
              </label>
              <input
                type="url"
                value={form.photo_url}
                onChange={e => setForm({...form, photo_url: e.target.value})}
                className="input-field"
                placeholder="https://your-photo-url.jpg"
              />
              <p style={{ color: 'var(--text-secondary)', fontSize: '0.75rem', marginTop: '0.5rem' }}>Paste a direct link to your photo</p>
            </div>

            {/* Bio */}
            <div className="input-group">
              <label className="input-label">About Me</label>
              <textarea
                value={form.bio}
                onChange={e => setForm({...form, bio: e.target.value})}
                className="input-field"
                placeholder="Tell others about yourself..."
                rows={3}
                style={{ resize: 'vertical', fontFamily: 'inherit' }}
                maxLength={300}
              />
              <p style={{ color: 'var(--text-secondary)', fontSize: '0.75rem', marginTop: '0.4rem', textAlign: 'right' }}>{form.bio.length}/300</p>
            </div>

            {/* Birthdate */}
            <div className="input-group">
              <label className="input-label">
                <Calendar size={14} style={{ display: 'inline', marginRight: '6px' }} />
                Date of Birth
              </label>
              <input
                type="date"
                value={form.birthdate}
                onChange={e => setForm({...form, birthdate: e.target.value})}
                className="input-field"
                style={{ colorScheme: 'dark' }}
              />
            </div>

            {/* Gender */}
            <div className="input-group">
              <label className="input-label">Gender</label>
              <div style={{ display: 'flex', gap: '0.75rem', flexWrap: 'wrap' }}>
                {['Male', 'Female', 'Non-binary', 'Prefer not to say'].map(g => (
                  <button
                    key={g}
                    type="button"
                    onClick={() => setForm({...form, gender: g.toLowerCase()})}
                    style={{
                      padding: '8px 20px',
                      borderRadius: '999px',
                      border: '1px solid',
                      borderColor: form.gender === g.toLowerCase() ? 'var(--accent-primary)' : 'var(--border-color)',
                      background: form.gender === g.toLowerCase() ? 'rgba(236,72,153,0.15)' : 'transparent',
                      color: form.gender === g.toLowerCase() ? 'var(--accent-primary)' : 'var(--text-secondary)',
                      cursor: 'pointer',
                      fontSize: '0.875rem',
                      fontWeight: '500',
                      transition: 'all 0.2s'
                    }}
                  >
                    {g}
                  </button>
                ))}
              </div>
            </div>

            {/* Interested In */}
            <div className="input-group">
              <label className="input-label">
                <Heart size={14} style={{ display: 'inline', marginRight: '6px' }} />
                Interested In
              </label>
              <div style={{ display: 'flex', gap: '0.75rem', flexWrap: 'wrap' }}>
                {['Men', 'Women', 'Everyone'].map(opt => (
                  <button
                    key={opt}
                    type="button"
                    onClick={() => setForm({...form, interested_in: opt.toLowerCase()})}
                    style={{
                      padding: '8px 20px',
                      borderRadius: '999px',
                      border: '1px solid',
                      borderColor: form.interested_in === opt.toLowerCase() ? 'var(--accent-secondary)' : 'var(--border-color)',
                      background: form.interested_in === opt.toLowerCase() ? 'rgba(139,92,246,0.15)' : 'transparent',
                      color: form.interested_in === opt.toLowerCase() ? 'var(--accent-secondary)' : 'var(--text-secondary)',
                      cursor: 'pointer',
                      fontSize: '0.875rem',
                      fontWeight: '500',
                      transition: 'all 0.2s'
                    }}
                  >
                    {opt}
                  </button>
                ))}
              </div>
            </div>

            <button type="submit" className="btn btn-primary" style={{ width: '100%', marginTop: '0.5rem', gap: '8px' }} disabled={saving}>
              {saving ? (
                <><Loader size={18} style={{ animation: 'spin 1s linear infinite' }} /> Saving...</>
              ) : (
                <><Save size={18} /> Save Changes</>
              )}
            </button>
          </form>
        </div>

        {/* Stats Card */}
        <div className="glass-panel animate-slide-up" style={{ padding: '2rem', marginTop: '1.5rem', marginBottom: '2rem' }}>
          <h3 style={{ margin: '0 0 1.25rem 0', color: 'var(--text-secondary)', fontSize: '0.875rem', textTransform: 'uppercase', letterSpacing: '0.1em' }}>Profile Preview</h3>
          <div style={{ display: 'flex', gap: '1rem', alignItems: 'flex-start' }}>
            <img src={photoUrl} alt="Preview" style={{ width: 70, height: 70, borderRadius: '12px', objectFit: 'cover' }} />
            <div>
              <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                <h3 style={{ margin: 0 }}>{user?.name}</h3>
                {getAge(form.birthdate) && <span style={{ color: 'var(--text-secondary)', fontSize: '1.1rem' }}>{getAge(form.birthdate)}</span>}
              </div>
              <p style={{ color: 'var(--text-secondary)', fontSize: '0.875rem', marginTop: '0.5rem', lineHeight: 1.5 }}>
                {form.bio || 'No bio yet. Add something about yourself!'}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
