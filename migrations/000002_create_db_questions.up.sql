CREATE TABLE IF NOT EXISTS questions (
    id bigserial PRIMARY KEY,
    category_id bigserial REFERENCES categories ON DELETE CASCADE ,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    web_index smallint,
    text text NOT NULL ,
    ans_options text[] NOT NULL,
    code_block text,
    correct_ans_opt text,
    correct_ans_explanation text,
    url text
);