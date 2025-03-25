DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS ticket_options;
DROP TABLE IF EXISTS purchases;
DROP TABLE IF EXISTS tickets;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE ticket_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT,
    allocation INTEGER NOT NULL CHECK (allocation >= 0),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE purchases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_option_id UUID NOT NULL REFERENCES ticket_options(id),
    user_id UUID NOT NULL REFERENCES users(id),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    purchase_id UUID NOT NULL REFERENCES purchases(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);


-- Create a test user for convenience
INSERT INTO users (id) VALUES ('d6abe829-c28c-44ec-bee6-3183f2c53fef');

-- Create a test ticket option
INSERT INTO ticket_options (id, name, description, allocation) 
VALUES ('969f4317-09f4-4b15-b8be-a87d40fb56fb', 'Test Event', 'A test event for the API', 100);
