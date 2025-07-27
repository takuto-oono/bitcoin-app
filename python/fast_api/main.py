import os
import argparse
from fastapi import FastAPI
import uvicorn
from dotenv import load_dotenv

parser = argparse.ArgumentParser(description="FastAPI Server")
parser.add_argument("--env", type=str, default="env/.env.local",
                    help="Path to .env file")
args = parser.parse_args()

load_dotenv(dotenv_path=args.env)

app = FastAPI()


@app.get("/")
def read_root():
    return {"message": "Hello Fast API!"}


def main():
    host = os.getenv("HOST", "0.0.0.0")
    port = int(os.getenv("PORT", 8001))
    reload = os.getenv("RELOAD", "false").lower() == "true"

    uvicorn.run("main:app", host=host, port=port, reload=reload)


if __name__ == "__main__":
    main()
