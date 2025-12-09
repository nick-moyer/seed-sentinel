from flask import Flask, request, jsonify
from datetime import datetime

app = Flask(__name__)

# --- CONFIGURATION ---
# In v0.2, the LLM will decide this dynamically based on plant type.
# For v0.1 POC, we hard-code the threshold.
MOISTURE_THRESHOLD = 30

@app.route('/analyze', methods=['POST'])
def analyze_plant():
    """
    Pure Logic Endpoint.
    Receives telemetry -> Returns advice.
    """
    # 1. Parse Data
    data = request.get_json() or {}

    plant_name = data.get('plant_name', 'Unknown Plant')
    moisture = data.get('moisture', 0)

    print(f"🧠 [Brain] Analyzing {plant_name}: Moisture={moisture}%")

    # 2. Mocked "AI" Logic
    alert_needed = False
    advice = "Conditions are optimal. Keep monitoring."

    if moisture < MOISTURE_THRESHOLD:
        alert_needed = True
        # The AI "Persona" generates the text
        advice = f"CRITICAL: Your {plant_name} is critically dry ({moisture}%). Water immediately to prevent root stress."

    # 3. Return the Decision to Go
    return jsonify({
        "timestamp": datetime.now().isoformat(),
        "plant_name": plant_name,
        "alert_needed": alert_needed,  # Go looks at this boolean
        "advice": advice               # Go sends this text to the user
    })

if __name__ == '__main__':
    # Listen on all interfaces so Docker can find it
    print("🧠 Seed Sentinel Brain listening on :5000...")
    app.run(host='0.0.0.0', port=5000)