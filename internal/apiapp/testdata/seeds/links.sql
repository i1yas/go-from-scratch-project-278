INSERT INTO links
    (id, original_url, shortcode)
OVERRIDING SYSTEM VALUE
VALUES
    (1, 'https://google.com', 'google'),
    (2, 'https://yandex.ru', 'yandex'),
    (3, 'https://github.com', 'github'),
    (4, 'http://gone', 'gone'),
    (5, 'https://youtube.com', 'youtube'),
    (6, 'https://en.wikipedia.org/wiki/Special:Random', 'random-wiki-article'),
    (7, 'https://go.dev', 'golang'),
    (8, 'https://test.com', 'very-very-long-code-101202303'),
    (9, 'https://test.com', 'TEST'),
    (10, 'https://example.com', 'example link');