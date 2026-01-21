import os
import json
import uvicorn
from fastapi import FastAPI
from pydantic import BaseModel
from typing import List
import google.generativeai as genai
from dotenv import load_dotenv

load_dotenv()

# 1. Setup Model
api_key = os.getenv("GEMINI_API_KEY")
if not api_key:
    print("⚠️ WARNING: GEMINI_API_KEY not set!")

genai.configure(api_key=api_key)

# We use the same model for both, but you could use a smaller/faster one for profiling
model = genai.GenerativeModel(
    'gemini-2.5-flash', # Using 2.5-flash for speed and reliability
    generation_config={"response_mime_type": "application/json"}
)

app = FastAPI()

# --- DATA STRUCTURES ---

# In-memory session profiling memory (ephemeral)
PROFILE_MEMORY = {}

class BrainRequest(BaseModel):
    session_id: str
    command: str
    cwd: str
    history: List[str]

# The Profiler's Output
class ProfileResponse(BaseModel):
    motivation: str  # e.g., "Miner", "Ransomware", "Scout"
    reasoning: str   # Why do we think this?
    suggested_bait: str # What should we show them?

# --- THE SHADOW AGENT (PROFILER) ---

def run_shadow_agent(req: BrainRequest) -> dict:
    """
    Analyzes the attacker's behavior to determine their intent.
    Does not generate output for the user, only for the System.
    """
    prior_profile = PROFILE_MEMORY.get(req.session_id)
    prior_profile_str = json.dumps(prior_profile) if prior_profile else "None"

    prompt = f"""
    You are a Cybersecurity Analyst (Shadow Agent) monitoring a honeypot.
    
    ANALYSIS TARGET:
    - Command: "{req.command}"
    - History: {req.history}
    - CWD: {req.cwd}
    - Prior Profile: {prior_profile_str}
    
    TASK:
    Classify the attacker's likely motivation based on their commands.
    
    Categories:
    1. "Miner": Checking CPU/GPU (lscpu, nvidia-smi), downloading mining tools (xmrig).
    2. "Ransomware": Searching for high-value data (backups, databases, financial docs).
    3. "APT/Espionage": Looking for network config, ssh keys, lateral movement targets.
    4. "Bot/Script": Running generic, fast reconnaissance scripts.
    5. "Scout": Just looking around (ls, whoami, id).
    
    JSON RESPONSE FORMAT:
    {{
        "motivation": "One of the categories above",
        "reasoning": "Brief explanation",
        "suggested_bait": "A specific file, process, or system trait to hallucinate that entices this specific attacker."
    }}
    """
    
    try:
        response = model.generate_content(prompt)
        profile = json.loads(response.text)
        PROFILE_MEMORY[req.session_id] = profile
        return profile
    except Exception:
        # Fallback profile if AI fails
        fallback = {"motivation": "Unknown", "reasoning": "Error", "suggested_bait": "Standard weak password file"}
        PROFILE_MEMORY[req.session_id] = fallback
        return fallback

# --- THE DIRECTOR (GENERATOR) ---

@app.post("/hallucinate")
async def hallucinate(req: BrainRequest):
    print(f"Incoming Command: {req.command}")
    
    # STEP 1: RUN THE SHADOW AGENT
    # We analyze them *before* we respond.
    profile = run_shadow_agent(req)
    print(f"🕵️ Shadow Profile: {profile['motivation']} | Bait: {profile['suggested_bait']}")

    # STEP 2: GENERATE THE TRAP
    # We feed the Profile into the Director so the response is tailored.
    director_prompt = f"""
    You are a high-interaction Ubuntu honeypot.
    
    INTELLIGENCE REPORT:
    - Attacker Motivation: {profile['motivation']}
    - Suggested Bait Strategy: {profile['suggested_bait']}
    
    CONTEXT:
    - Dir: {req.cwd}
    - Command: "{req.command}"

    TASK:
    Generate 3 distinct terminal outputs in JSON.
    
    1. "standard": Realistic, boring output.
    2. "bait": HIGH PRIORITY. You MUST include the '{profile['suggested_bait']}' in this output to hook the {profile['motivation']}. 
       - If they are Miners, show a fake NVIDIA GPU or high CPU count.
       - If they are Ransomware, show a 'backup.tar.gz' or 'customers.sql'.
       - If they are APT, show a '.ssh/config' or 'id_rsa'.
    3. "stall": A realistic error/delay (e.g., 'Mirror unreachable', 'Permission denied (try sudo?)').

    Response Format:
    {{
        "standard": "...",
        "bait": "...",
        "stall": "..."
    }}
    """

    try:
        response = model.generate_content(director_prompt)
        data = json.loads(response.text)
        
        # --- ACTIVE DEFENSE LOGIC ---
        # For now, we still return "standard" to keep the shell usable.
        # But we log the "bait" option which the Dashboard will eventually use.
        
        # NOTE: If the command is purely informational (like 'lscpu' or 'df'), 
        # we might want to default to 'bait' immediately if the profile matches!
        
        selected_output = data.get("standard", "")
        
        # AUTO-ESCALATION:
        # If the shadow agent detects high-value intent (Miner/Ransomware),
        # we can automatically swap to the 'bait' response to keep them interested.
        if profile['motivation'] in ["Miner", "Ransomware", "APT"]:
             if data.get("bait"):
                 print("⚠️  AUTO-ESCALATION: Switching to BAIT response to hook attacker.")
                 selected_output = data.get("bait")

        return {"output": selected_output}

    except Exception as e:
        print(f"❌ Error: {e}")
        return {"output": "bash: command not found"}

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=5000)