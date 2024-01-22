CREATE TABLE category (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE product (
    id INTEGER PRIMARY KEY,
    category_id INTEGER NOT NULL,
    ordering INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    price TEXT NOT NULL,
    discounted_price TEXT NOT NULL,
    FOREIGN KEY (category_id) REFERENCES category(id),
    UNIQUE (category_id, ordering)
);