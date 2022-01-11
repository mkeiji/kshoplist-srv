CREATE TABLE Items(
    ID SERIAL PRIMARY KEY,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    StoreID BIGINT UNSIGNED NOT NULL,
    Name varchar(255) NOT NULL,

    FOREIGN KEY (StoreID) REFERENCES Stores(ID) ON DELETE CASCADE ON UPDATE CASCADE
);
