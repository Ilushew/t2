-- +goose Up
-- +goose StatementBegin

-- Обновляем человеко-читаемые названия и описания для существующих мест

-- 1. Kiyasovo - Эко-прогулка
UPDATE places SET 
    name_label = 'Экотропа в Киясово',
    description_label = 'Прогуляйтесь по живописной экотропе в Киясово. Уникальная возможность познакомиться с природой Удмуртии вдали от городского шума.'
WHERE name = 'Kiyasovo_ecowalk';

-- 2. Bright Colors - Выставка
UPDATE places SET 
    name_label = 'Выставка «Яркие краски Удмуртии»',
    description_label = 'Погрузитесь в мир удмуртского искусства. Выставка знакомит с традиционными ремёслами и современным творчеством региона.'
WHERE name = 'bright_colors_udm';

-- 3. Чайковский - Музей
UPDATE places SET 
    name_label = 'Музей П.И. Чайковского',
    description_label = 'Посетите дом-музей великого композитора в Воткинске. Здесь начинался творческий путь автора «Лебединого озера».'
WHERE name = 'chaikovsky';

-- 4. Diadem Kama - Исторический маршрут
UPDATE places SET 
    name_label = 'По берегам Камы',
    description_label = 'Путешествие вдоль великой реки Кама. Откройте для себя историю и красоту прикамских земель.'
WHERE name = 'diadem_kama';

-- 5. Eco Kok - Экоцентр
UPDATE places SET 
    name_label = 'Экотропа в Игринском районе',
    description_label = 'Прогулка по живописным тропам Игринского района. Чистый воздух, вековые деревья и тишина удмуртской природы.'
WHERE name = 'ecowalk_kok';

-- 6. Fall Love - Гастрономия
UPDATE places SET 
    name_label = 'Осенняя гастрономия Удмуртии',
    description_label = 'Попробуйте лучшие блюда удмуртской кухни. Субботи — национальные пельмени, табани — лепёшки и многое другое.'
WHERE name = 'fall_love_udm';

-- 7. Глазов - Исторический маршрут
UPDATE places SET 
    name_label = 'Исторический Глазов',
    description_label = 'Прогулка по старинному городу Глазов. Уникальная архитектура и богатая история финно-угорского края.'
WHERE name = 'glazov';

-- 8. Музей Глазова
UPDATE places SET 
    name_label = 'Краеведческий музей Глазова',
    description_label = 'Музей хранит историю края от древности до наших дней. Уникальные экспонаты культуры коми-пермяков и удмуртов.'
WHERE name = 'glazov_town_museum';

-- 9. Волонтёрство
UPDATE places SET 
    name_label = 'Волонтёрская программа',
    description_label = 'Примите участие в волонтёрских проектах Удмуртии. Помогите сохранить природу и культурное наследие региона.'
WHERE name = 'help_in_sep';

-- 10. Иднакар
UPDATE places SET 
    name_label = 'Городище Иднакар',
    description_label = 'Древнее удмуртское городище IX–XIII веков близ Глазова. Археологический памятник федерального значения.'
WHERE name = 'idnakar';

-- 11. Ижик
UPDATE places SET 
    name_label = 'Прогулка по Ижевску',
    description_label = 'Обзорная экскурсия по столице Удмуртии. Ижевский пруд, набережная, главные достопримечательности города.'
WHERE name = 'izhik';

-- 12. Лудорвай
UPDATE places SET 
    name_label = 'Этнографический комплекс «Лудорвай»',
    description_label = 'Музей-заповедник под открытым небом. Уникальный ветряная мельница, крестьянские усадьбы и живая история удмуртского быта.'
WHERE name = 'ludorvai';

-- 13. Мобильный экоцентр
UPDATE places SET 
    name_label = 'Мобильный экоцентр в Игре',
    description_label = 'Экологический центр в посёлке Игра. Программы по защите природы и волонтёрские инициативы.'
WHERE name = 'mob_ecocenter';

-- 14. Нечкино
UPDATE places SET 
    name_label = 'Национальный парк «Нечкино»',
    description_label = 'Живописный национальный парк на берегу Камы. Пешие тропы, смотровые площадки и богатый животный мир.'
WHERE name = 'nechkino';

-- 15. Парфюмерный музей
UPDATE places SET 
    name_label = 'Музей парфюмерии в Сарапуле',
    description_label = 'Уникальная коллекция исторических ароматов. Узнайте о традициях производства духов в Прикамье.'
WHERE name = 'parfume_museum';

-- 16. Пельмень-тур
UPDATE places SET 
    name_label = 'Пельмень-тур по Удмуртии',
    description_label = 'Гастрономическое путешествие по лучшим пельменным Удмуртии. Попробуйте суботи, пельмени и другие национальные блюда.'
WHERE name = 'pelmen_thour';

-- 17. Лес посадки
UPDATE places SET 
    name_label = 'Экологическая акция «Посади лес»',
    description_label = 'Присоединяйтесь к посадке деревьев в Удмуртии. Волонтёрская акция по восстановлению лесов региона.'
WHERE name = 'plant_forest';

-- 18. Сарапул
UPDATE places SET 
    name_label = 'Купеческий Сарапул',
    description_label = 'Прогулка по историческому центру Сарапула. Уникальная купеческая архитектура и музей истории города.'
WHERE name = 'sarapul';

-- 19. Вкусная Удмуртия
UPDATE places SET 
    name_label = 'Гастрономический тур «Вкусная Удмуртия»',
    description_label = 'Три дня знакомства с удмуртской кухой. Лучшие рестораны, мастер-классы и дегустации национальных блюд.'
WHERE name = 'tasty_udm';

-- 20. Воткинск эко
UPDATE places SET 
    name_label = 'Экотропа в Воткинске',
    description_label = 'Прогулка по экологической тропе Воткинского пруда. Красивые пейзажи и чистый воздух родины Чайковского.'
WHERE name = 'votkinsk_ecowalk';

-- 21. Уикенд в Игре
UPDATE places SET 
    name_label = 'Уикенд в Игре',
    description_label = 'Два дня в посёлке Игра. Исторические места, природа и аутентичная культура удмуртского края.'
WHERE name = 'weekend_in_igra';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

UPDATE places SET name_label = NULL, description_label = NULL;

-- +goose StatementEnd
