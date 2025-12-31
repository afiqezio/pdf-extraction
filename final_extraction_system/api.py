import uvicorn
from fastapi import FastAPI, BackgroundTasks, HTTPException
from pydantic import BaseModel
import uuid
from typing import Dict, Any, Optional
import logging
from extractor import extract_tables_smart, logger
import os

app = FastAPI(title="Oil & Gas Extraction Service")

# Simple in-memory job store
jobs: Dict[str, Dict[str, Any]] = {}

class ExtractionRequest(BaseModel):
    file_path: str

class ExtractionResponse(BaseModel):
    job_id: str
    status: str

def run_extraction(job_id: str, file_path: str):
    """
    Background task to run the extraction process.
    """
    logger.info(f"Starting background extraction for job {job_id}, file: {file_path}")
    try:
        if not os.path.exists(file_path):
             jobs[job_id]["status"] = "failed"
             jobs[job_id]["error"] = f"File not found: {file_path}"
             return

        # Run the extraction
        # extract_tables_smart returns List[Dict] (pages)
        results = extract_tables_smart(file_path, save_output=True)
        
        jobs[job_id]["status"] = "completed"
        jobs[job_id]["data"] = results
        logger.info(f"Job {job_id} completed successfully")
        
    except Exception as e:
        logger.error(f"Job {job_id} failed: {e}")
        jobs[job_id]["status"] = "failed"
        jobs[job_id]["error"] = str(e)

@app.post("/extract", response_model=ExtractionResponse)
async def start_extraction(request: ExtractionRequest, background_tasks: BackgroundTasks):
    job_id = str(uuid.uuid4())
    
    jobs[job_id] = {
        "id": job_id,
        "status": "pending",
        "file_path": request.file_path,
        "created_at": str(os.path.getctime(request.file_path) if os.path.exists(request.file_path) else None)
    }
    
    # Start background task
    jobs[job_id]["status"] = "processing"
    background_tasks.add_task(run_extraction, job_id, request.file_path)
    
    return {"job_id": job_id, "status": "processing"}

@app.get("/result/{job_id}")
async def get_result(job_id: str):
    if job_id not in jobs:
        raise HTTPException(status_code=404, detail="Job not found")
    
    return jobs[job_id]

@app.get("/health")
async def health_check():
    return {"status": "ok"}

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)

