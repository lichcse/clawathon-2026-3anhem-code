import pathlib

from sqlmodel import Session, SQLModel, create_engine

from . import models  # noqa: F401 — register tables

DB_PATH = pathlib.Path(__file__).resolve().parent.parent.parent / "data.db"
DATABASE_URL = f"sqlite:///{DB_PATH}"

engine = create_engine(DATABASE_URL, connect_args={"check_same_thread": False})


def init_db() -> None:
    SQLModel.metadata.create_all(engine)


def get_session():
    with Session(engine) as session:
        yield session
