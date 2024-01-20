select word from words inner join translations on words.id = translations.word_id where words_views is null ;
-- Replace 'specific_user_id' with the actual user_id you want to check against
SELECT 
    w.word,
    t.translation AS associated_translation,
    t_random.translation AS random_translation
FROM 
    words AS w
JOIN translations AS t ON w.id = t.word_id
JOIN user_current_theme AS uct ON w.theme_id = uct.theme_id
LEFT JOIN LATERAL (
    SELECT 
        tr.translation 
    FROM 
        translations AS tr
    JOIN words AS w2 ON tr.word_id = w2.id AND w2.theme_id = uct.theme_id
    WHERE 
        tr.word_id <> w.id
    ORDER BY 
        RANDOM()
    LIMIT 1
) t_random ON true
WHERE 
    w.id NOT IN (
        SELECT word_id 
        FROM words_views 
        WHERE user_id = 'ba566903-f42c-4171-9af6-52194a0f5da5'
    )
AND uct.user_id = 'ba566903-f42c-4171-9af6-52194a0f5da5'
