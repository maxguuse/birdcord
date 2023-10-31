-- name: GetUserByIdForPoll :one
SELECT * FROM voted_users
         WHERE discord_id = $1 AND poll_id = $2;

-- name: AddVotedUser :exec
INSERT INTO voted_users (
    discord_id,
    option_id,
    poll_id
) VALUES (
    $1, $2, $3
) RETURNING *;