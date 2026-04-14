const fetch = require('node-fetch'); // wait node 18+ has built-in fetch
async function run() {
  try {
    const res = await fetch('http://127.0.0.1:8000/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'Accept': 'application/json' },
      body: JSON.stringify({ email: 'alice@example.com', password: 'password' })
    });
    console.log(res.status, await res.text());
  } catch(e) { console.error(e); }
}
run();
