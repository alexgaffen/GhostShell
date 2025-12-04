import os
import uvicorn
from fastapi import FastAPI
from pydantic import BaseModel
import google.generativeai as genai

from dotenv import load_dotenv
load_dotenv()  # This loads the .env file into os.getenv

# 1. Setup the Model
# Best practice: Read from Environment Variable (we will set this in Docker/Terminal)
api_key = os.getenv("GEMINI_API_KEY")

genai.configure(api_key=api_key)
model = genai.GenerativeModel('gemini-2.5-flash') # Flash is faster/cheaper for this

app = FastAPI()

# 2. Define the Data Structure (Must match Go's JSON)
class BrainRequest(BaseModel):
    command: str

@app.post("/hallucinate")
async def hallucinate(req: BrainRequest):
    print(f"🧠 Brain received: {req.command}")
    
    # 3. The "System Prompt" - This is where the magic happens
    # We tell Gemini it is NOT an AI, but a Ubuntu server.
    prompt = f"""
    You are a high-interaction honeypot simulating a Ubuntu Linux server.
    The attacker just typed this command: "{req.command}"
    
    RULES:
    1. Output ONLY the standard terminal output for this command.
    2. Do not explain anything. Do not say "Here is the output".
    3. If the command returns nothing (like 'cp' or 'rm'), output nothing.
    4. If the command is 'ls', invent some realistic files (maybe a file named 'passwords.txt' or 'wallet.dat').
    5. Be consistent.
    
    Command: {req.command}
    Output:
    """

    try:
        response = model.generate_content(prompt)
        # .strip() helps remove accidental newlines added by the model
        return {"output": response.text.strip()}
    except Exception as e:
        return {"output": f"Error: {str(e)}"}

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=5000)