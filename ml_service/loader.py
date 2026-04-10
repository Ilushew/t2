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

df = pd.read_csv('places.csv', encoding='utf-8')

# Оставляем только нужные колонки, переименовываем и конвертируем в список словарей
df_places = (
    df
    .rename(columns={'place_id': 'id', 'description': 'desc'})
    .astype({'id': int})
)


if os.path.exists(CACHE_FILE):
    print(f"✅ Найдены готовые эмбеддинги ({CACHE_FILE}). Загружаем...")
    with open(CACHE_FILE, "rb") as f:
        place_embeddings = pickle.load(f)
else:
    print("⏳ Эмбеддингов нет. Генерируем...")
    place_embeddings = model.encode(df_places['desc'].tolist(), show_progress_bar=True)
    
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
        print(result[['name', 'interest_score']])
        return result[['name', 'interest_score']]
    
    def tabular_model_predict(self):
        pass

    def predict(self, request):
        pred = self.embedding_model_predict(request.query)
        indices = pred.index.tolist()
        return list(indices)




