from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from sqlmodel import Session

from .api import blocks, chat, runs, uploads, workflows
from .db.session import engine, init_db
from .registry import seed_blocks


@asynccontextmanager
async def lifespan(app: FastAPI):
    init_db()
    with Session(engine) as session:
        n = seed_blocks(session)
        print(f"[startup] seeded {n} blocks from disk")
    yield


app = FastAPI(title="Block Chat", lifespan=lifespan)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(blocks.router)
app.include_router(chat.router)
app.include_router(runs.router)
app.include_router(uploads.router)
app.include_router(workflows.router)


@app.get("/healthz")
def healthz() -> dict[str, bool]:
    return {"ok": True}
