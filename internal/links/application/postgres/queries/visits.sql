-- name: GetVisitsTotalCount :one
SELECT count(*)
FROM visits;

-- name: CreateVisit :execrows
INSERT INTO visits
    (link_id, ip, referer, user_agent, status)
VALUES
    ($1, $2, $3, $4, $5);