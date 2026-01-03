import { useState, useEffect } from 'react'
import axios from 'axios'
import logo from '../../assets/logo.png'

function App() {
  const [loading, setLoading] = useState(true)
  const [sensor, setSensor] = useState(null)
  const [plant, setPlant] = useState(null)

  // Form State
  const [newName, setNewName] = useState('')
  const [newDate, setNewDate] = useState('') // New State for Date

  const [isSaving, setIsSaving] = useState(false)
  const [saveStatus, setSaveStatus] = useState(null)

  const fetchData = async () => {
    try {
      const sensorsRes = await axios.get('/api/sensors')
      if (!sensorsRes.data?.[0]?.ID) {
        setLoading(false)
        return
      }

      const sensor = sensorsRes.data[0]
      setSensor(sensor) // Assume single device for now

      const plantRes = await axios.get(`/api/plant?sensorID=${sensor.ID}`)
      if (plantRes.data.length === 0) {
        setLoading(false)
        return
      }

      const plant = plantRes.data[0]
      setPlant(plant)
      setNewName(plant.plant_name || '')
      setNewDate(plant.date_planted ? plant.date_planted.split('T')[0] : '') // Set date in YYYY-MM-DD format

      setLoading(false)
    } catch (err) {
      console.error(err)
      setLoading(false)
    }
  }

  const handleSave = async () => {
    if (!sensor) return
    setIsSaving(true)
    setSaveStatus(null)

    try {
      // Send both Name and Date
      await axios.put(`/api/configure`, {
        SensorID: sensor.ID,
        Name: newName,
        DatePlanted: newDate
      })

      setSaveStatus('success')
      fetchData() // Refresh to sync state
      setTimeout(() => setSaveStatus(null), 3000)
    } catch (err) {
      console.error(err)
      setSaveStatus('error')
    } finally {
      setIsSaving(false)
    }
  }

  useEffect(() => {
    fetchData()
    // const interval = setInterval(fetchData, 2000)
    // return () => clearInterval(interval)
  }, [])

  if (loading) return <div className="loading-screen">Connecting...</div>
  if (!sensor) return <div className="loading-screen">No Sensors Found</div>

  return (
    <div className="app-layout">
      <nav>
        <div className="nav-content">
          <span className="brand">🌱 Seed Sentinel</span>
          <span className="status-indicator online">Online</span>
        </div>
      </nav>

      <main>
     {/* Data Card */}
        <div className="card hero-card">
          <img src={logo} alt="Plant Sentinel Logo" className="plant-logo" />
          <h1 className="plant-title">{plant?.Name ?? 'Unnamed Plant'}</h1>
          <p className="mac-label">{sensor.ID}</p>

          {/* Calibration Footer */}
          <div className="card-footer-stats">
            <div className="stat-item">
              <span className="stat-label">Dry Ref (0%)</span>
              <span className="stat-value dry">{sensor.DryReference ?? '--'}</span>
            </div>
            <div className="stat-separator"></div>
            <div className="stat-item">
              <span className="stat-label">Wet Ref (100%)</span>
              <span className="stat-value wet">{sensor.WetReference ?? '--'}</span>
            </div>
          </div>
        </div>

        {/* Configuration Card */}
        <div className="card settings-card">
          <div className="card-header">
            <h3>Configuration</h3>
          </div>

          <div className="form-grid">
            {/* Input 1: Name */}
            <div className="input-group">
              <label>Plant Name</label>
              <input
                type="text"
                value={newName}
                onChange={(e) => setNewName(e.target.value)}
                placeholder="e.g. Tomato"
                disabled={isSaving}
              />
            </div>

            {/* Input 2: Date Planted */}
            <div className="input-group">
              <label>Date Planted</label>
              <input
                type="date"
                value={newDate}
                onChange={(e) => setNewDate(e.target.value)}
                disabled={isSaving}
              />
            </div>

            {/* Save Button */}
            <button
              onClick={handleSave}
              disabled={isSaving}
              className={`full-width ${saveStatus}`}
            >
              {isSaving ? 'Saving...' : (saveStatus === 'success' ? 'Settings Saved' : 'Save Changes')}
            </button>
          </div>
        </div>
      </main>
    </div>
  )
}

export default App
