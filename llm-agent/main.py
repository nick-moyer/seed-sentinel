from flask import Flask, request, jsonify
from datetime import datetime
import ollama
import json

app = Flask(__name__)

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

    print(f"[Brain] Analyzing {plant_name}: Moisture={moisture}%")

    # 2. Construct the Prompt
    prompt = f"""
    You are an expert botanist growing {plant_name} from seed.
    The current soil moisture is {moisture}%.

    Determine if this is dangerous for this specific plant type.

    Return ONLY a JSON object with this format (do not include markdown formatting):
    {{
        "alert_needed": yes/no,
        "advice": "Short, actionable advice here."
    }}
    """

    # 3. Call the Local LLM
    response = ollama.chat(model='llama3', messages=[
        {'role': 'user', 'content': prompt},
    ], format='json')

    # 4. Extract the alert flag and advice
    content = json.loads(response['message']['content'])
    alert_needed = content['alert_needed']
    advice = content['advice']

    print(f"Alert Needed: {alert_needed}")
    print(f"Llama says: {advice}")

    # 5. Return the analysis
    return jsonify({
        "timestamp": datetime.now().isoformat(),
        "plant_name": plant_name,
        "alert_needed": True if alert_needed == 'yes' else False,
        "advice": advice
    })

if __name__ == '__main__':
    # Listen on all interfaces so Docker can find it
    print("Seed Sentinel Brain listening on :5000...")
    app.run(host='0.0.0.0', port=5000)