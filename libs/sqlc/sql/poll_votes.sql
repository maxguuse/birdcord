-- name: GetVote :one
SELECT COUNT(*) FROM poll_votes 
WHERE user_id = $1 AND poll_id = $2;

-- name: GetPollVotes :many
SELECT * FROM poll_votes
WHERE poll_id = $1;

-- name: AddVote :exec
INSERT INTO poll_votes (
    user_id, 
    poll_id, 
    option_id
) VALUES (
    $1, $2, $3
);