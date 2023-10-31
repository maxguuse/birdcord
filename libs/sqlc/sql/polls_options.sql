-- name: CreateOption :one
INSERT INTO polls_options (
    title,
    poll_id
) VALUES (
    $1, $2
) RETURNING *;

-- SELECT polls_options.id, polls_options.title, COUNT(voted_users.id)::integer AS vote_count
-- FROM polls_options
--          LEFT JOIN voted_users ON polls_options.id = voted_users.option_id
-- WHERE polls_options.poll_id = $1
-- GROUP BY polls_options.id, polls_options.title;

-- name: GetOptionsWithVotesCount :many
SELECT polls_options.id, polls_options.title, COUNT(voted_users.id)::integer AS vote_count
FROM polls_options
         LEFT JOIN voted_users ON polls_options.id = voted_users.option_id
WHERE polls_options.poll_id = $1
GROUP BY polls_options.id, polls_options.title
ORDER BY polls_options.id;