
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import Optional, Dict, List
import logging
from contextlib import asynccontextmanager

from loader import ReccomendationModel

logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")
logger = logging.getLogger(__name__)

class RecommendationRequest(BaseModel):
    duration: str
    company: str
    has_car: bool
    budget: int
    interests: List[str]
    with_pets: bool
    query: str

    
    

recommender: Optional[ReccomendationModel] = None

@asynccontextmanager
async def lifespan(app: FastAPI):
    global recommender
    try:
        logger.info("🧠 Инициализация ML-модели...")
        recommender = ReccomendationModel()
        logger.info("✅ Модель готова к работе!")
    except Exception as e:
        logger.error(f"❌ Ошибка при инициализации модели: {e}")
        recommender = None
    yield
    recommender = None

app = FastAPI(title="Udmurtia Tourism Recommender", lifespan=lifespan)

@app.post("/predict", response_model=List[int])
def get_recommendations(req: RecommendationRequest):
    if recommender is None:
        raise HTTPException(status_code=503, detail="Модель не загружена.")
    
    try:
        return recommender.predict(req)
    except Exception as e:
        logger.error(f"Ошибка: {e}", exc_info=True)
        raise HTTPException(status_code=500, detail="Внутренняя ошибка")


if __name__ == "__main__":
    import uvicorn
    uvicorn.run("server:app", host="0.0.0.0", port=8000, reload=False)