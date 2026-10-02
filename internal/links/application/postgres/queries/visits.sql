-- name: GetVisitsTotalCount :one
SELECT count(*)
FROM visits;