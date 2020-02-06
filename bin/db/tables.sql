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
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    --  updated_at timestamp NOT NULL DEFAULT (DATETIME('now', 'localtime')),
    --  created_at timestamp NOT NULL DEFAULT (DATETIME('now', 'localtime')),
    UNIQUE (region, weight, deadline)
); 

CREATE TRIGGER freight_region_trigger_updated_time
AFTER UPDATE On freight_region
BEGIN
   UPDATE freight_region SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
   --  UPDATE freight_region SET timestamp = STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW') WHERE id = NEW.id;
END;
--  CREATE UNIQUE INDEX t1b ON freight_region(region, weight, deadline);
