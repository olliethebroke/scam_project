-- +goose Up
alter table users rename column username to firstname;
alter table leaderboard rename column username to firstname;
alter table friends rename column user_id to invited_user_id;
alter table friends rename column friend_id to inviting_user_id;

-- +goose Down
alter table users rename column firstname to username;
alter table leaderboard rename column firstname to username;
alter table friends rename column invited_user_id to user_id;
alter table friends rename column inviting_user_id to friend_id;
