import { useState, useEffect } from 'react'
import axios from 'axios'
import './index.css'

function App() {
  const [token, setToken] = useState(localStorage.getItem('token') || '')
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  
  const [sites, setSites] = useState([])
  const [liveLogs, setLiveLogs] = useState([])

  const [newName, setNewName] = useState('')
  const [newUrl, setNewUrl] = useState('')

  const handleLogin = async (e) => {
    e.preventDefault()
    try {
      const res = await axios.post('http://localhost:8080/api/login', { username, password })
      const receivedToken = res.data.token
      setToken(receivedToken)
      localStorage.setItem('token', receivedToken)
    } catch (err) {
      alert("Błąd logowania!")
    }
  }

  const handleLogout = () => {
    setToken('')
    localStorage.removeItem('token')
  }

  useEffect(() => {
    if (token) {
      axios.get('http://localhost:8080/api/sites', {
        headers: { Authorization: `Bearer ${token}` }
      })
      .then(res => setSites(res.data))
      .catch(err => console.log("Błąd pobierania stron", err))
    }
  }, [token])

  useEffect(() => {
    if (!token) return;
    const ws = new WebSocket('ws://localhost:8080/ws')
    
    ws.onopen = () => console.log("Połączono z WebSockets!")
    
    ws.onmessage = (event) => {
      const incomingLogs = JSON.parse(event.data)
      setLiveLogs(prevLogs => {
        let updatedLogs = [...prevLogs]
        incomingLogs.forEach(newLog => {
          const existingIndex = updatedLogs.findIndex(log => log.url === newLog.url)
          if (existingIndex !== -1) {
            updatedLogs[existingIndex] = newLog
          } else {
            updatedLogs.unshift(newLog)
          }
        })
        return updatedLogs
      })
    }
    return () => ws.close()
  }, [token])

  const handleAddSite = async (e) => {
    e.preventDefault()
    try {
      const res = await axios.post('http://localhost:8080/api/sites', {
        name: newName,
        url: newUrl,
        interval: 15
      }, {
        headers: { Authorization: `Bearer ${token}` }
      })

      setSites(prev => [...prev, res.data])

      setNewName('')
      setNewUrl('')
    } catch (err) {
      console.log("Błąd dodawania", err)
      alert("Zły format danych lub błąd serwera.")
    }
  }

  const handleDeleteSite = async (id) => {
    if (!window.confirm("Na pewno usunąć tę stronę z monitoringu?")) return;

    try {
      const siteToDelete = sites.find(s => s.id === id)

      await axios.delete(`http://localhost:8080/api/sites/${id}`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      
      setSites(prevSites => prevSites.filter(site => site.id !== id));
      
      if (siteToDelete) {
        setLiveLogs(prevLogs => prevLogs.filter(log => log.url !== siteToDelete.url));
      }
      
    } catch (err) {
      console.log("Błąd usuwania", err);
      alert("Nie udało się usunąć strony. Sprawdź konsolę.");
    }
  }

  if (!token) {
    return (
      <div className="container">
        <h2>Uptime Monitor - Logowanie</h2>
        <form onSubmit={handleLogin} className="card">
          <input placeholder="Login" value={username} onChange={e => setUsername(e.target.value)} />
          <input type="password" placeholder="Hasło" value={password} onChange={e => setPassword(e.target.value)} />
          <button type="submit">Zaloguj</button>
        </form>
      </div>
    )
  }

  return (
    <div className="container">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h2>Uptime Monitor Dashboard</h2>
        <button onClick={handleLogout}>Wyloguj</button>
      </div>

      <div className="card">
        <h3>Dodaj nową stronę</h3>
        <form onSubmit={handleAddSite} style={{ display: 'flex', gap: '10px', flexWrap: 'wrap' }}>
          <input 
            placeholder="Nazwa (np. Google)" 
            value={newName} 
            onChange={e => setNewName(e.target.value)} 
            required 
          />
          <input 
            type="url" 
            placeholder="Adres URL (https://...)" 
            value={newUrl} 
            onChange={e => setNewUrl(e.target.value)} 
            required 
          />
          <button type="submit" style={{ backgroundColor: '#03dac6', color: '#000' }}>
            Dodaj
          </button>
        </form>
      </div>

      <div className="card">
        <h3>Twoje Monitorowane Strony</h3>
        {sites.length === 0 ? <p>Brak stron. Dodaj jakieś u góry!</p> : (
          <ul style={{ listStyle: 'none', padding: 0 }}>
            {sites.map(site => (
              <li key={site.id} style={{ 
                display: 'flex', justifyContent: 'space-between', alignItems: 'center',
                padding: '10px 0', borderBottom: '1px solid #333'
              }}>
                <div>
                  <strong>{site.name}</strong> <br/>
                  <small style={{ color: '#888' }}>{site.url}</small>
                </div>
                <button onClick={() => handleDeleteSite(site.id)} style={{ backgroundColor: '#cf6679', color: '#000', padding: '5px 15px' }}>
                  Usuń
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>

      <div className="card">
        <h3>Live Logs (WebSockets)</h3>
        {liveLogs.length === 0 ? <p>Czekam na dane z serwera...</p> : (
          <div>
            {liveLogs.map((log, index) => (
              <div key={index} className="log-item">
                <strong>{log.site_name}</strong> - Status: {log.status_code} ({log.response_time_ms} ms)
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

export default App