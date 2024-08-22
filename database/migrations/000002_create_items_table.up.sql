CREATE TABLE Items (
    ID SERIAL PRIMARY KEY,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    StoreID BIGINT NOT NULL,
    Name varchar(255) NOT NULL,

    FOREIGN KEY (StoreID) REFERENCES Stores(ID) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create a trigger function to update the UpdatedAt column
CREATE OR REPLACE FUNCTION update_items_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.UpdatedAt = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger that calls the function before an update
CREATE TRIGGER update_items_updated_at
BEFORE UPDATE ON Items
FOR EACH ROW
EXECUTE FUNCTION update_items_updated_at_column();
