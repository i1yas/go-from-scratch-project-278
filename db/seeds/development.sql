INSERT INTO links (original_url, shortcode) VALUES
    ('https://yandex.ru', 'яндкекс'),
    ('https://github.com', 'github'),
    ('http://gone', 'gone'),
    ('https://google.com', 'google'),
    ('https://youtube.com', 'youtube'),
    ('https://en.wikipedia.org/wiki/Special:Random', 'random-wiki-article'),
    ('https://go.dev', 'golang'),
    ('https://test.com', 'very-very-long-code-101202303'),
    ('https://test.com', 'TEST'),
    ('https://example.com', 'example link')
ON CONFLICT DO NOTHING;
