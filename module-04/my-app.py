from fastapi import FastAPI
import datetime

app = FastAPI()

@app.get("/")
def root():
    return {
        "message": "Hello from Python and Kate with hot-reload!",
        "time": str(datetime.datetime.now())
    }

@app.get("/test")
def test():
    return {"status": "ok", "data": "you can change me!"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=5000)