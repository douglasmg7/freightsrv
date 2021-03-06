-- FREIGHT REGION
INSERT INTO freight_region(region, weight, deadline, price) VALUES ("north", 4000, 5, 10000);
INSERT INTO freight_region(region, weight, deadline, price) VALUES ("north", 100000, 5, 12000);

INSERT INTO freight_region(region, weight, deadline, price) VALUES ("northeast", 4000, 10, 10000);
INSERT INTO freight_region(region, weight, deadline, price) VALUES ("northeast", 100000, 10, 12000);

INSERT INTO freight_region(region, weight, deadline, price) VALUES ("midwest", 4000, 5, 7500);
INSERT INTO freight_region(region, weight, deadline, price) VALUES ("midwest", 100000, 5, 9500);

INSERT INTO freight_region(region, weight, deadline, price) VALUES ("southeast", 4000, 2, 8360);
INSERT INTO freight_region(region, weight, deadline, price) VALUES ("southeast", 4000, 5, 5000);
INSERT INTO freight_region(region, weight, deadline, price) VALUES ("southeast", 100000, 5, 7000);

INSERT INTO freight_region(region, weight, deadline, price) VALUES ("south", 4000, 3, 3040);
INSERT INTO freight_region(region, weight, deadline, price) VALUES ("south", 4000, 5, 5000);
INSERT INTO freight_region(region, weight, deadline, price) VALUES ("south", 100000, 10, 7000);

-- MOTOBOY FREIGHT
INSERT INTO motoboy_freight(city, city_norm, deadline, price) VALUES ("Belo Horizonte", "belo-horizonte", 1, 7520);
INSERT INTO motoboy_freight(city, city_norm, deadline, price) VALUES ("Conceição do Mato Dentro", "conceicao-do-mato-dentro", 2, 8545);
INSERT INTO motoboy_freight(city, city_norm, deadline, price) VALUES ("Guarupé", "guarupe", 3, 9545);
INSERT INTO motoboy_freight(city, city_norm, deadline, price) VALUES ("Sabará", "sabara", 1, 10000);

-- DEALER FREIGHT
INSERT INTO dealer_freight(dealer, weight, deadline, price) VALUES ("aldo", 4000, 6, 11000);
INSERT INTO dealer_freight(dealer, weight, deadline, price) VALUES ("aldo", 100000, 6, 13000);

INSERT INTO dealer_freight(dealer, weight, deadline, price) VALUES ("allnations_rj", 4000, 5, 10000);
INSERT INTO dealer_freight(dealer, weight, deadline, price) VALUES ("allnations_rj", 100000, 5, 12000);

INSERT INTO dealer_freight(dealer, weight, deadline, price) VALUES ("allnations_es", 4000, 3, 11100);
INSERT INTO dealer_freight(dealer, weight, deadline, price) VALUES ("allnations_es", 100000, 3, 12200);

INSERT INTO dealer_freight(dealer, weight, deadline, price) VALUES ("allnations_sc", 4000, 4, 12200);
INSERT INTO dealer_freight(dealer, weight, deadline, price) VALUES ("allnations_sc", 100000, 4, 13300);
