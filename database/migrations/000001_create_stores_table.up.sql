CREATE TABLE Stores (
    ID SERIAL PRIMARY KEY,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    Name varchar(255) NOT NULL
);

-- Create a trigger function to update the UpdatedAt column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.UpdatedAt = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger that calls the function before an update
CREATE TRIGGER update_stores_updated_at
BEFORE UPDATE ON Stores
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
