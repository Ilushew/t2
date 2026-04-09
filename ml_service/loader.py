from sentence_transformers import SentenceTransformer
import numpy as np
import pandas as pd
from sklearn.metrics.pairwise import cosine_similarity

import os
import pickle

# Модель скачивается с HuggingFace и кэшируется в ~/.cache/huggingface/
# При перезапуске контейнера с volume — используется кэш, без volume — скачивается заново
MODEL_NAME = os.getenv("MODEL_NAME", "sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2")

print(f"🧠 Загрузка модели: {MODEL_NAME}...")
model = SentenceTransformer(MODEL_NAME)
print("✅ Модель загружена!")


CACHE_FILE = "places_embeddings.pkl"

places_data = [
    {"id": 1, "name": "Музей Калашникова", "desc": "Музейно-выставочный комплекс имени Михаила Калашникова в Ижевске. Экспозиция об оружии, истории завода и биографии конструктора."},
    {"id": 2, "name": "Этнопарк Зура", "desc": "Интерактивная площадка с удмуртской культурой, традициями, национальной кухней и ремеслами. Подходит для семей с детьми."},
    {"id": 3, "name": "Нечкинский парк", "desc": "Национальный парк на берегу Камы. Горнолыжный курорт, лес, природа, активный отдых, лыжи зимой."},
    {"id": 4, "name": "Свято-Михайловский собор", "desc": "Православный храм в Ижевске, памятник архитектуры, визитная карточка города, религия и история."}
]
df_places = pd.DataFrame(places_data)
df_places['text_for_embed'] = df_places['name'] + ". " + df_places['desc']


if os.path.exists(CACHE_FILE):
    print(f"✅ Найдены готовые эмбеддинги ({CACHE_FILE}). Загружаем...")
    with open(CACHE_FILE, "rb") as f:
        place_embeddings = pickle.load(f)
else:
    print("⏳ Эмбеддингов нет. Генерируем...")
    place_embeddings = model.encode(df_places['text_for_embed'].tolist(), show_progress_bar=True)
    
    print(f"💾 Сохраняем эмбеддинги в {CACHE_FILE}...")
    with open(CACHE_FILE, "wb") as f:
        pickle.dump(place_embeddings, f)
    print("✅ Сохранено!")


class ReccomendationModel:
    def __init__(self):
        pass

    def embedding_model_predict(self, user_query, top_k=5):
        '''находит косинусное сходство эмбединга запроса пользователя с существующими вариантами из базы'''

        query_embedding = model.encode([user_query])
        
    
        cosine_sim = cosine_similarity(query_embedding, place_embeddings)[0]
    
        semantic_score = 1 / (1 + np.exp(-(cosine_sim + 1) / 2 * 10 + 5))

        df_places['interest_score'] = semantic_score
        result = df_places.sort_values(by='interest_score', ascending=False).head(top_k)
        
        return result[['name', 'interest_score']]
    
    def tabular_model_predict(self):
        pass

    def predict(self):
        self.embedding_model_predict("хочу чупа чупс")
        return [5, 7, 2, 1, 6]




mdl = ReccomendationModel()
#print(mdl.predict())
