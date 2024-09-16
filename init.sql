CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE employee (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);

CREATE TABLE organization (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organization_responsible (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUID REFERENCES employee(id) ON DELETE CASCADE
);


INSERT INTO employee (id, username, first_name, last_name, created_at, updated_at)
VALUES ('bff35258-2e34-4458-8445-e8bbbc4d1c71', 'testuser1', 'Алексей', 'Алексеев', '2024-09-11T19:35:04+03:00', '2024-09-11T19:35:04+03:00'),
       ('2064757a-28b9-4ae4-aec9-d8a1f5c39ba0', 'testuser2', 'Павел', 'Павлов', '2024-09-11T19:49:02+03:00', '2024-09-11T19:49:02+03:00');

INSERT INTO organization (id, name, description, type, created_at, updated_at)
VALUES ('f7586a5b-de63-4c20-990c-a6bec950ec56', 'organization1', 'description1', 'LLC', '2024-09-11T19:41:44+03:00', '2024-09-11T19:41:44+03:00');

INSERT INTO organization_responsible (organization_id, user_id)
VALUES ('f7586a5b-de63-4c20-990c-a6bec950ec56', 'bff35258-2e34-4458-8445-e8bbbc4d1c71');


CREATE TYPE tender_status_enum AS ENUM ('Created', 'Published', 'Closed');

CREATE TABLE tender (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    status tender_status_enum NOT NULL,
    current_version INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id UUID NOT NULL REFERENCES employee(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organization(id) ON DELETE CASCADE
);

CREATE TYPE service_type_enum AS ENUM ('Construction', 'Delivery', 'Manufacture');

CREATE TABLE tender_versions (
    id SERIAL PRIMARY KEY,
    version INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    service_type service_type_enum NOT NULL,
    tender_id UUID NOT NULL REFERENCES tender(id) ON DELETE CASCADE,
    UNIQUE (tender_id, version)
);

CREATE TYPE bid_status_enum AS ENUM ('Created', 'Published', 'Canceled');

CREATE TYPE author_type_enum AS ENUM ('Organization', 'User');

CREATE TABLE bid (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    status bid_status_enum NOT NULL,
    tender_id UUID UNIQUE NOT NULL REFERENCES tender(id) ON DELETE CASCADE,
    author_type author_type_enum NOT NULL,
    author_id UUID NOT NULL,
    current_version INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
);

CREATE TABLE bid_versions (
    id SERIAL PRIMARY KEY,
    version INTEGER PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    bid_id UUID REFERENCES bid(id) ON DELETE CASCADE,
    UNIQUE (bid_id, version)
);


CREATE TYPE decision_enum AS ENUM ('Approved', 'Rejected');

CREATE bid_decisions (
    id SERIAL PRIMARY KEY,
    decision decision_enum,
    organization_id UUID NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES employee(id) ON DELETE CASCADE,
    bid_id UUID REFERENCES bid(id) ON DELETE CASCADE
);


CREATE TABLE bid_review (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    description VARCHAR(1000) NOT NULL,
    organization_id UUID NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES employee(id) ON DELETE CASCADE,
    bid_id UUID REFERENCES bid(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);