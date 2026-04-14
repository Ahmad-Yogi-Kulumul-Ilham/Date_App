import { useState, useEffect, useRef } from 'react';
import { Send, MessageCircle, Clock, RefreshCw } from 'lucide-react';

const POLL_INTERVAL = 3000; // poll every 3 seconds

function timeAgo(dateStr) {
  const date = new Date(dateStr);
  const now = new Date();
  const diff = Math.floor((now - date) / 1000);
  if (diff < 60) return 'just now';
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`;
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`;
  return date.toLocaleDateString();
}

function formatTime(dateStr) {
  const date = new Date(dateStr);
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

export default function MatchesPage() {
  const [matches, setMatches] = useState([]);
  const [selectedMatch, setSelectedMatch] = useState(null);
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState('');
  const [isPolling, setIsPolling] = useState(false);
  const [lastMessageCount, setLastMessageCount] = useState(0);
  const [unreadCounts, setUnreadCounts] = useState({});
  const [sending, setSending] = useState(false);
  const messagesEndRef = useRef(null);
  const pollRef = useRef(null);
  const currentUser = JSON.parse(localStorage.getItem('user'));

  useEffect(() => {
    fetchMatches();
  }, []);

  useEffect(() => {
    if (selectedMatch) {
      fetchMessages(selectedMatch.id);
      startPolling(selectedMatch.id);
    }
    return () => stopPolling();
  }, [selectedMatch]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const startPolling = (matchId) => {
    stopPolling();
    pollRef.current = setInterval(async () => {
      try {
        const res = await fetch(`http://localhost:8000/api/matches/${matchId}/messages`, {
          headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}`, 'Accept': 'application/json' }
        });
        const data = await res.json();
        setMessages(prev => {
          if (data.length !== prev.length) {
            setIsPolling(true);
            setTimeout(() => setIsPolling(false), 800);
            return data;
          }
          return prev;
        });
      } catch(e) {}
    }, POLL_INTERVAL);
  };

  const stopPolling = () => {
    if (pollRef.current) clearInterval(pollRef.current);
  };

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
    if (!newMessage.trim() || !selectedMatch || sending) return;

    const content = newMessage;
    setNewMessage('');
    setSending(true);

    // Optimistic update
    const optimisticMsg = {
      id: `temp-${Date.now()}`,
      content,
      sender_id: currentUser?.id,
      created_at: new Date().toISOString(),
      _optimistic: true
    };
    setMessages(prev => [...prev, optimisticMsg]);

    try {
      const res = await fetch(`http://localhost:8000/api/matches/${selectedMatch.id}/messages`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Accept': 'application/json'
        },
        body: JSON.stringify({ content })
      });
      const data = await res.json();
      // Replace optimistic with real
      setMessages(prev => prev.map(m => m._optimistic ? data : m));
    } catch(e) {
      // Remove optimistic on error
      setMessages(prev => prev.filter(m => !m._optimistic));
      setNewMessage(content);
    } finally {
      setSending(false);
    }
  };

  const getOtherUser = (match) => {
    return match.user1?.name === currentUser?.name ? match.user2 : match.user1;
  };

  const handleSelectMatch = (match) => {
    setSelectedMatch(match);
    // Clear unread
    setUnreadCounts(prev => ({ ...prev, [match.id]: 0 }));
  };

  return (
    <div className="flex h-full w-full">
      {/* Matches List Sidebar */}
      <div className="glass-panel" style={{ width: '300px', minWidth: '240px', borderRadius: 0, borderRight: '1px solid rgba(255,255,255,0.05)', display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
        {/* Sidebar Header */}
        <div style={{ padding: '1.5rem', borderBottom: '1px solid rgba(255,255,255,0.05)', flexShrink: 0 }}>
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
            <h2 style={{ fontSize: '1.25rem', margin: 0, color: 'white' }}>Messages</h2>
            <span style={{ background: 'rgba(236,72,153,0.15)', color: 'var(--accent-primary)', padding: '2px 10px', borderRadius: '999px', fontSize: '0.8rem', fontWeight: '600' }}>
              {matches.length}
            </span>
          </div>
          <p style={{ color: 'var(--text-secondary)', fontSize: '0.8rem', marginTop: '0.25rem' }}>Your matches</p>
        </div>

        {/* Matches List */}
        <div style={{ flex: 1, overflowY: 'auto', width: '100%' }}>
          {matches.length === 0 ? (
            <div style={{ padding: '2rem', textAlign: 'center' }}>
              <MessageCircle size={40} style={{ color: 'var(--text-secondary)', margin: '0 auto 1rem', display: 'block' }} />
              <p style={{ color: 'var(--text-secondary)', fontSize: '0.875rem' }}>No matches yet.<br/>Start swiping!</p>
            </div>
          ) : null}
          {matches.map(match => {
            const otherUser = getOtherUser(match);
            const isSelected = selectedMatch?.id === match.id;
            return (
              <div
                key={match.id}
                onClick={() => handleSelectMatch(match)}
                style={{
                  padding: '0.875rem 1.25rem',
                  display: 'flex',
                  alignItems: 'center',
                  gap: '12px',
                  cursor: 'pointer',
                  background: isSelected ? 'rgba(236, 72, 153, 0.1)' : 'transparent',
                  borderLeft: isSelected ? '3px solid var(--accent-primary)' : '3px solid transparent',
                  transition: 'background 0.2s',
                }}
              >
                {/* Avatar */}
                <div style={{ width: '46px', height: '46px', borderRadius: '50%', overflow: 'hidden', flexShrink: 0, border: isSelected ? '2px solid var(--accent-primary)' : '2px solid transparent' }}>
                  <img
                    src={otherUser?.profile?.photo_url || `https://ui-avatars.com/api/?name=${otherUser?.name}&background=ec4899&color=fff`}
                    style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                    alt="Avatar"
                  />
                </div>
                {/* Info */}
                <div style={{ flex: 1, minWidth: 0 }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <h4 style={{ margin: 0, fontSize: '0.9rem', color: isSelected ? 'white' : 'var(--text-primary)', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>
                      {otherUser?.name}
                    </h4>
                  </div>
                  <p style={{ margin: '2px 0 0', fontSize: '0.78rem', color: 'var(--text-secondary)' }}>
                    {isSelected ? '● chatting' : 'Tap to chat'}
                  </p>
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* Chat Area */}
      {selectedMatch ? (
        <div style={{ flex: 1, display: 'flex', flexDirection: 'column', background: 'rgba(15, 23, 42, 0.3)', overflow: 'hidden' }}>
          {/* Chat Header */}
          <div style={{ padding: '1rem 1.5rem', borderBottom: '1px solid rgba(255,255,255,0.05)', background: 'var(--bg-secondary)', display: 'flex', alignItems: 'center', gap: '1rem', flexShrink: 0 }}>
            <div style={{ position: 'relative' }}>
              <img
                src={getOtherUser(selectedMatch)?.profile?.photo_url || `https://ui-avatars.com/api/?name=${getOtherUser(selectedMatch)?.name}&background=ec4899&color=fff`}
                style={{ width: 42, height: 42, borderRadius: '50%', objectFit: 'cover' }}
                alt="Avatar"
              />
              {/* Online indicator */}
              <div style={{ position: 'absolute', bottom: 2, right: 2, width: 10, height: 10, borderRadius: '50%', background: 'var(--success)', border: '2px solid var(--bg-secondary)' }} />
            </div>
            <div>
              <h3 style={{ margin: 0, fontSize: '1rem' }}>{getOtherUser(selectedMatch)?.name}</h3>
              <p style={{ margin: 0, fontSize: '0.78rem', color: 'var(--success)' }}>Active now</p>
            </div>
            <div style={{ marginLeft: 'auto', display: 'flex', alignItems: 'center', gap: '6px' }}>
              {isPolling && (
                <div style={{ display: 'flex', alignItems: 'center', gap: '4px', color: 'var(--text-secondary)', fontSize: '0.75rem' }}>
                  <RefreshCw size={12} style={{ animation: 'spin 1s linear infinite' }} />
                  New messages
                </div>
              )}
            </div>
          </div>

          {/* Messages */}
          <div style={{ flex: 1, overflowY: 'auto', padding: '1.5rem', display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
            {messages.length === 0 && (
              <div style={{ textAlign: 'center', color: 'var(--text-secondary)', marginTop: '2rem' }}>
                <MessageCircle size={40} style={{ margin: '0 auto 0.75rem', display: 'block', opacity: 0.4 }} />
                <p style={{ fontSize: '0.875rem' }}>No messages yet. Say hi! 👋</p>
              </div>
            )}
            {messages.map((msg, idx) => {
              const otherUser = getOtherUser(selectedMatch);
              const isMine = msg.sender_id !== otherUser?.id;
              const isTemp = msg._optimistic;

              return (
                <div key={msg.id || idx} style={{ display: 'flex', justifyContent: isMine ? 'flex-end' : 'flex-start', alignItems: 'flex-end', gap: '8px' }}>
                  {!isMine && (
                    <img src={otherUser?.profile?.photo_url || `https://ui-avatars.com/api/?name=${otherUser?.name}&size=32&background=ec4899&color=fff`} style={{ width: 28, height: 28, borderRadius: '50%', flexShrink: 0 }} alt="" />
                  )}
                  <div style={{ display: 'flex', flexDirection: 'column', alignItems: isMine ? 'flex-end' : 'flex-start', gap: '2px' }}>
                    <div style={{
                      background: isMine ? 'var(--accent-gradient)' : 'var(--card-bg)',
                      border: isMine ? 'none' : '1px solid rgba(255,255,255,0.05)',
                      padding: '10px 14px',
                      borderRadius: '18px',
                      borderBottomRightRadius: isMine ? '4px' : '18px',
                      borderBottomLeftRadius: isMine ? '18px' : '4px',
                      maxWidth: '65%',
                      boxShadow: '0 2px 8px rgba(0,0,0,0.15)',
                      opacity: isTemp ? 0.7 : 1,
                      fontSize: '0.9rem',
                      lineHeight: '1.4'
                    }}>
                      {msg.content}
                    </div>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '4px', color: 'var(--text-secondary)', fontSize: '0.7rem', padding: '0 4px' }}>
                      <Clock size={10} />
                      {isTemp ? 'Sending...' : formatTime(msg.created_at)}
                    </div>
                  </div>
                </div>
              );
            })}
            <div ref={messagesEndRef} />
          </div>

          {/* Message Input */}
          <form onSubmit={sendMessage} style={{ padding: '1rem 1.5rem', background: 'var(--bg-secondary)', borderTop: '1px solid rgba(255,255,255,0.05)', display: 'flex', gap: '0.75rem', alignItems: 'center', flexShrink: 0 }}>
            <input
              type="text"
              value={newMessage}
              onChange={e => setNewMessage(e.target.value)}
              className="input-field"
              style={{ marginBottom: 0, borderRadius: '999px', background: 'var(--bg-primary)', flex: 1 }}
              placeholder="Type a message..."
              disabled={sending}
            />
            <button
              type="submit"
              className="btn btn-primary"
              style={{ padding: '12px', borderRadius: '50%', flexShrink: 0, width: 46, height: 46 }}
              disabled={!newMessage.trim() || sending}
            >
              <Send size={18} />
            </button>
          </form>
        </div>
      ) : (
        <div style={{ flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center' }}>
          <div style={{ width: 100, height: 100, borderRadius: '50%', background: 'rgba(255,255,255,0.02)', display: 'flex', alignItems: 'center', justifyContent: 'center', marginBottom: '1.5rem' }}>
            <MessageCircle size={48} color="var(--text-secondary)" />
          </div>
          <h2 style={{ color: 'var(--text-secondary)', fontWeight: 400, textAlign: 'center', maxWidth: '300px', lineHeight: 1.5 }}>Select a match to start chatting</h2>
        </div>
      )}
    </div>
  );
}
