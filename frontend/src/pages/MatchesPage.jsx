import { useState, useEffect } from 'react';
import { Send } from 'lucide-react';

export default function MatchesPage() {
  const [matches, setMatches] = useState([]);
  const [selectedMatch, setSelectedMatch] = useState(null);
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState('');
  const currentUser = JSON.parse(localStorage.getItem('user'));

  useEffect(() => {
    fetchMatches();
  }, []);

  useEffect(() => {
    if (selectedMatch) {
      fetchMessages(selectedMatch.id);
    }
  }, [selectedMatch]);

  const fetchMatches = async () => {
    try {
      const res = await fetch('http://localhost:8000/api/matches', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}`, 'Accept': 'application/json' }
      });
      const data = await res.json();
      setMatches(data);
    } catch(e) { console.error(e); }
  };

  const fetchMessages = async (matchId) => {
    try {
      const res = await fetch(`http://localhost:8000/api/matches/${matchId}/messages`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}`, 'Accept': 'application/json' }
      });
      const data = await res.json();
      setMessages(data);
    } catch(e) { console.error(e); }
  };

  const sendMessage = async (e) => {
    e.preventDefault();
    if (!newMessage.trim() || !selectedMatch) return;

    try {
      const res = await fetch(`http://localhost:8000/api/matches/${selectedMatch.id}/messages`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Accept': 'application/json'
        },
        body: JSON.stringify({ content: newMessage })
      });
      const data = await res.json();
      setMessages([...messages, data]);
      setNewMessage('');
    } catch(e) { console.error(e); }
  };

  const getOtherUser = (match) => {
    return match.user1.name === currentUser.name ? match.user2 : match.user1;
  };

  return (
    <div className="flex h-full w-full">
      {/* Matches List */}
      <div className="glass-panel items-center" style={{ width: '320px', borderRadius: 0, borderRight: '1px solid rgba(255,255,255,0.05)', display: 'flex', flexDirection: 'column' }}>
        <div style={{ padding: '1.5rem', width: '100%', borderBottom: '1px solid rgba(255,255,255,0.05)' }}>
          <h2 style={{ fontSize: '1.5rem', margin: 0, color: 'white' }}>Messages</h2>
        </div>
        <div style={{ flex: 1, overflowY: 'auto', width: '100%' }}>
          {matches.length === 0 ? <p style={{ padding: '1.5rem', color: 'var(--text-secondary)' }}>No matches yet.</p> : null}
          {matches.map(match => {
            const otherUser = getOtherUser(match);
            const isSelected = selectedMatch?.id === match.id;
            return (
              <div 
                key={match.id}
                onClick={() => setSelectedMatch(match)}
                style={{ 
                  padding: '1rem 1.5rem', 
                  display: 'flex', 
                  alignItems: 'center', 
                  gap: '12px', 
                  cursor: 'pointer',
                  background: isSelected ? 'rgba(236, 72, 153, 0.1)' : 'transparent',
                  borderLeft: isSelected ? '4px solid var(--accent-primary)' : '4px solid transparent',
                  transition: 'background 0.2s'
                }}
              >
                <div style={{ width: '48px', height: '48px', borderRadius: '50%', overflow: 'hidden', flexShrink: 0 }}>
                  <img src={otherUser.profile?.photo_url || `https://ui-avatars.com/api/?name=${otherUser.name}`} style={{ width: '100%', height: '100%', objectFit: 'cover' }} alt="Avatar" />
                </div>
                <div>
                  <h4 style={{ margin: 0, fontSize: '1rem', color: isSelected ? 'white' : 'var(--text-primary)' }}>{otherUser.name}</h4>
                  <p style={{ margin: 0, fontSize: '0.875rem', color: 'var(--text-secondary)' }}>Chat now</p>
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* Chat Area */}
      {selectedMatch ? (
        <div className="flex-col h-full" style={{ flex: 1, background: 'rgba(15, 23, 42, 0.3)' }}>
          <div style={{ padding: '1.5rem', borderBottom: '1px solid rgba(255,255,255,0.05)', background: 'var(--bg-secondary)', display: 'flex', alignItems: 'center', gap: '1rem' }}>
             <img src={getOtherUser(selectedMatch).profile?.photo_url || `https://ui-avatars.com/api/?name=${getOtherUser(selectedMatch).name}`} style={{ width: 40, height: 40, borderRadius: '50%' }} />
             <h3 style={{ margin: 0 }}>{getOtherUser(selectedMatch).name}</h3>
          </div>
          
          <div style={{ flex: 1, overflowY: 'auto', padding: '1.5rem', display: 'flex', flexDirection: 'column', gap: '1rem' }}>
            {messages.map(msg => {
              const isMine = msg.sender_id === (currentUser.id || (matches.find(m => m.id === selectedMatch.id).user_id_1 === currentUser.id ? null : null)); 
              // Basic check, better to use proper user ID comparison
              // Here we assume sender_id isn't the other user's ID
              const otherUser = getOtherUser(selectedMatch);
              const actuallyMine = msg.sender_id !== otherUser.id;

              return (
                <div key={msg.id} style={{ display: 'flex', justifyContent: actuallyMine ? 'flex-end' : 'flex-start' }}>
                  <div style={{ 
                    background: actuallyMine ? 'var(--accent-gradient)' : 'var(--card-bg)',
                    padding: '12px 16px',
                    borderRadius: '16px',
                    borderBottomRightRadius: actuallyMine ? '4px' : '16px',
                    borderBottomLeftRadius: actuallyMine ? '16px' : '4px',
                    maxWidth: '70%',
                    boxShadow: '0 4px 6px -1px rgba(0,0,0,0.1)'
                  }}>
                    {msg.content}
                  </div>
                </div>
              )
            })}
          </div>

          <form onSubmit={sendMessage} style={{ padding: '1rem 1.5rem', background: 'var(--bg-secondary)', display: 'flex', gap: '1rem', alignItems: 'center' }}>
            <input 
              type="text" 
              value={newMessage} 
              onChange={(e) => setNewMessage(e.target.value)} 
              className="input-field" 
              style={{ marginBottom: 0, borderRadius: '999px', background: 'var(--bg-primary)' }}
              placeholder="Type your message..."
            />
            <button type="submit" className="btn btn-primary" style={{ padding: '12px', borderRadius: '50%' }}>
              <Send size={20} />
            </button>
          </form>
        </div>
      ) : (
        <div className="flex-col h-full w-full items-center justify-center" style={{ flex: 1 }}>
          <div style={{ width: 100, height: 100, borderRadius: '50%', background: 'rgba(255,255,255,0.02)', display: 'flex', alignItems: 'center', justifyContent: 'center', marginBottom: '1.5rem' }}>
             <MessageCircle size={48} color="var(--text-secondary)" />
          </div>
          <h2 style={{ color: 'var(--text-secondary)' }}>Select a match to start chatting</h2>
        </div>
      )}
    </div>
  );
}
