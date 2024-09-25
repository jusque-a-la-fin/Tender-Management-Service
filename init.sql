DO $CREATION$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tender_status_enum') THEN
        CREATE TYPE tender_status_enum AS ENUM ('Created', 'Published', 'Closed');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'service_type_enum') THEN
        CREATE TYPE service_type_enum AS ENUM ('Construction', 'Delivery', 'Manufacture');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'bid_status_enum') THEN
        CREATE TYPE bid_status_enum AS ENUM ('Created', 'Published', 'Canceled');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'author_type_enum') THEN
        CREATE TYPE author_type_enum AS ENUM ('Organization', 'User');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'decision_enum') THEN
        CREATE TYPE decision_enum AS ENUM ('Approved', 'Rejected');
    END IF;
END $CREATION$;

-- tender - это тендер
CREATE TABLE IF NOT EXISTS tender (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    status tender_status_enum NOT NULL,
    current_version INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id UUID NOT NULL REFERENCES employee(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organization(id) ON DELETE CASCADE
);

-- tender_versions - это версии тендеров
CREATE TABLE IF NOT EXISTS tender_versions (
    id SERIAL PRIMARY KEY,
    version INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    service_type service_type_enum NOT NULL,
    tender_id UUID NOT NULL REFERENCES tender(id) ON DELETE CASCADE,
    UNIQUE (tender_id, version)
);

-- bid - это предложение для тендера
CREATE TABLE IF NOT EXISTS bid (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    status bid_status_enum NOT NULL,
    tender_id UUID NOT NULL REFERENCES tender(id) ON DELETE CASCADE,
    author_type author_type_enum NOT NULL,
    author_id UUID NOT NULL,
    current_version INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- bid_versions - это версии предложений для тендеров
CREATE TABLE IF NOT EXISTS bid_versions (
    id SERIAL PRIMARY KEY,
    version INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    bid_id UUID REFERENCES bid(id) ON DELETE CASCADE,
    UNIQUE (bid_id, version)
);

-- bid_decisions - решения по предложениям
CREATE TABLE IF NOT EXISTS bid_decisions (
    id SERIAL PRIMARY KEY,
    decision decision_enum,
    organization_id UUID NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES employee(id) ON DELETE CASCADE,
    bid_id UUID REFERENCES bid(id) ON DELETE CASCADE
);

-- bid_review - отзыв по предложению
CREATE TABLE IF NOT EXISTS bid_review (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    description VARCHAR(1000) NOT NULL,
    organization_id UUID NOT NULL REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES employee(id) ON DELETE CASCADE,
    bid_id UUID REFERENCES bid(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
