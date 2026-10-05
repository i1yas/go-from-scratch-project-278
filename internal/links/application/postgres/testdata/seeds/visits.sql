INSERT INTO visits
    (id, link_id, created_at, ip, referer, user_agent, status)
OVERRIDING SYSTEM VALUE
VALUES
    (1, 1, '2026-09-19T07:56:07.602Z', '36.8.244.30', 'https://example.org/nested/page', 'very-very-long-user-agent', 302),
    (2, 4, '2026-08-18T10:36:27.247Z', '144.190.53.110', 'https://example.net/', 'very-very-long-user-agent', 302),
    (3, 8, '2026-09-16T19:34:05.582Z', '77.85.246.33', 'https://test.com/', 'very-very-long-user-agent', 302),
    (4, 2, '2026-09-22T08:30:56.491Z', '221.116.74.94', 'https://test.net/page', NULL, 302),
    (5, 4, '2026-08-18T20:10:34.658Z', '65.97.24.6', 'https://example.com/page', 'very-very-long-user-agent', 500),
    (6, 8, '2026-08-13T08:27:13.026Z', '127.237.221.252', 'https://test.net/page', 'very-very-long-user-agent', 500),
    (7, 5, '2026-08-17T06:42:16.056Z', '0.49.240.5', NULL, NULL, 302),
    (8, 8, '2026-09-06T21:10:48.619Z', '255.231.53.10', 'http://test.net/page?param=testparam', 'user-agent-2', 302),
    (9, 5, '2026-09-25T10:24:02.104Z', '62.39.156.128', 'http://test.net/page?param=testparam', 'very-very-long-user-agent', 500),
    (10, 1, '2026-09-20T16:16:00.913Z', '232.231.65.143', 'https://test.net/page?param=testparam', 'user-agent-1', 302),
    (11, 7, '2026-09-11T18:34:01.049Z', '36.8.244.30', 'https://test.org/page', 'user-agent-2', 302),
    (12, 5, '2026-08-16T17:03:54.917Z', '144.190.53.110', 'http://example.org/nested/page', NULL, 302),
    (13, 5, '2026-09-11T22:51:24.345Z', '77.85.246.33', 'http://example.com/page', 'very-very-long-user-agent', 500),
    (14, 9, '2026-09-26T16:39:17.733Z', '221.116.74.94', 'http://website.org/nested/page', NULL, 302),
    (15, 2, '2026-09-16T17:38:19.730Z', '65.97.24.6', 'https://example.net/page?param=testparam', 'user-agent-1', 302);