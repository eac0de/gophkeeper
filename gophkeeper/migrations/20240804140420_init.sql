-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    user_text_data (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL,
        name VARCHAR(255) NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        data VARCHAR(1024) NOT NULL,
        metadata JSONB
    );

CREATE TABLE
    user_auth_info (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL,
        name VARCHAR(255) NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        login VARCHAR(255) NOT NULL,
        password VARCHAR(255) NOT NULL,
        metadata JSONB
    );

CREATE TABLE
    user_file_data (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL,
        name VARCHAR(255) NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        path_to_file VARCHAR(255) NOT NULL,
        ext VARCHAR(16) NOT NULL,
        metadata JSONB
    );

CREATE TABLE
    user_bank_card (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL,
        name VARCHAR(255) NOT NULL,
        created_at TIMESTAMP NOT NULL,
        updated_at TIMESTAMP NOT NULL,
        number VARCHAR(255) NOT NULL,
        card_holder VARCHAR(255) NOT NULL,
        expire_date VARCHAR(255) NOT NULL,
        csc SMALLINT NOT NULL,
        metadata JSONB
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE user_text_data;

DROP TABLE user_auth_info;

DROP TABLE user_file_data;

DROP TABLE user_bank_card;

-- +goose StatementEnd