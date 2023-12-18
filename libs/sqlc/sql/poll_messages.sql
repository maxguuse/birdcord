-- name: CreatePollMessage :exec
INSERT INTO poll_messages (
    poll_id, 
    message_id
) VALUES (
    $1, $2
);