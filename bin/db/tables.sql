-- drop table if exists entrance;
-- enable foreign keys
-- not working, reset to off when back to db
pragma foreign_keys = on;

-- Freight by region.
CREATE TABLE IF NOT EXISTS freight_region  (
    id integer primary key autoincrement,
    region varchar(64) CHECK(region IN ('north', 'northeast', 'midwest', 'southeast', 'south')) not null,
    weight integer CHECK(weight>=100) not null,    -- g
    deadline integer CHECK(deadline>0) not null,  -- days
    price integer CHECK(price>0) not null,     -- R$ X 100
    createdAt datetime not null,
    updatedAt datetime not null,
    UNIQUE (region, weight, deadline, price)
);
