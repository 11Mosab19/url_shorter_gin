CREATE OR REPLACE FUNCTION update_at_url_update()
RETURNS TRIGGER AS $$
BEGIN 
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURNS NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER urls_update_at_trigger
BEFORE UPDATE ON urls
FOR EACH ROW
EXECUTE FUNCTION update_at_url_update();

------------------------------------------------

CREATE OR REPLACE FUNCTION increase_clicks_urls()
RETURN TRIGGER AS $$
BEGIN 
    UPDATE urls
    SET total_clicks = total_clicks + 1
    WHERE id = NEW.url_id;
    RETURNS NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER url_total_clicks_update
AFTER INSERT analytics
FOR EACH ROW 
EXECUTE FUNCTION increase_clicks_urls();

-----------------------------------------------

CREATE OR REPLACE FUNCTION update_at_users_update()
RETURNS TRIGGER AS $$
BEGIN 
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURNS NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_update_at_trigger
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_at_users_update();