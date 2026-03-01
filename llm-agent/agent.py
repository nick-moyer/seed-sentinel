from flask import Flask, request, jsonify
from datetime import datetime
import ollama
import json
import os

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
    plant_age_days = data.get('plant_age_days', 0)
    moisture = data.get('moisture_percentage', 0)

    print(f"[Brain] Analyzing {plant_name}: Moisture={moisture}%, Age={plant_age_days} days")

    # 2. Construct the Prompt
    prompt_path = os.path.join(os.path.dirname(__file__), 'prompts', 'analyze_plant.txt')
    with open(prompt_path, 'r') as f:
        prompt_template = f.read()

    prompt = prompt_template.format(
        plant_name=plant_name,
        plant_age_days=plant_age_days,
        moisture=moisture
    )

    # 3. Call the Local LLM
    model = os.getenv('OLLAMA_MODEL', 'llama3')
    response = ollama.chat(model=model, messages=[
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