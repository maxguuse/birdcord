-- name: CreatePollMessage :exec
INSERT INTO poll_messages (
    poll_id, 
    message_id
) VALUES (
    $1, $2
);

-- name: GetMessagesForPollById :many
SELECT * FROM poll_messages WHERE poll_id = $1;

-- name: GetPollMessageById :one
SELECT * FROM poll_messages WHERE message_id = $1;

-- name: DeletePollMessageById :exec
DELETE FROM poll_messages WHERE id = $1;