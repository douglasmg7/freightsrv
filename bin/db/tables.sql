-- drop table if exists entrance;
-- enable foreign keys
-- not working, reset to off when back to db
pragma foreign_keys = on;

-- Freight by region.
CREATE TABLE IF NOT EXISTS freight_region  (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    region VARCHAR(64) CHECK(region IN ('north', 'northeast', 'midwest', 'southeast', 'south')) NOT NULL,
    weight INTEGER CHECK(weight >= 100) NOT NULL,    -- g
    deadline INTEGER CHECK(deadline > 0) NOT NULL,  -- days
    price INTEGER CHECK(price>0) NOT NULL,     -- R$ X 100
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    --  updated_at timestamp NOT NULL DEFAULT (DATETIME('now', 'localtime')),
    --  created_at timestamp NOT NULL DEFAULT (DATETIME('now', 'localtime')),
    UNIQUE (region, weight, deadline)
); 

CREATE TRIGGER IF NOT EXISTS freight_region_trigger_updated_at
AFTER UPDATE ON freight_region
BEGIN
   UPDATE freight_region SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
   --  UPDATE freight_region SET timestamp = STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW') WHERE id = NEW.id;
END;
--  CREATE UNIQUE INDEX t1b ON freight_region(region, weight, deadline);

-- Motoboy freight.
CREATE TABLE IF NOT EXISTS motoboy_freight  (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    state VARCHAR(64) CHECK(state IN ('mg')) NOT NULL DEFAULT 'mg',
    city VARCHAR(64)  NOT NULL,
    city_norm VARCHAR(64)  NOT NULL, -- Normalized city name.
    deadline INTEGER CHECK(deadline > 0) NOT NULL,  -- days
    price INTEGER CHECK(price > 0) NOT NULL,     -- R$ X 100
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (state, city_norm)
); 

CREATE TRIGGER IF NOT EXISTS motoboy_freight_trigger_updated_at
AFTER UPDATE ON motoboy_freight
BEGIN
   UPDATE motoboy_freight SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- Dealer freight.
CREATE TABLE IF NOT EXISTS dealer_freight (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    dealer VARCHAR(64) NOT NULL,
    weight INTEGER CHECK(weight >= 100) NOT NULL,    -- g
    deadline INTEGER CHECK(deadline > 0) NOT NULL,  -- days
    price INTEGER CHECK(price>=0) NOT NULL,     -- R$ X 100
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (dealer, weight, deadline)
); 

CREATE TRIGGER IF NOT EXISTS dealer_freight_trigger_updated_at
AFTER UPDATE ON dealer_freight
BEGIN
   UPDATE dealer_freight SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
