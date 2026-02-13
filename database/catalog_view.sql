-- View for catalog with full hierarchy (for optimized queries)
CREATE OR REPLACE VIEW catalog_view AS
SELECT 
    t.trim_id,
    t.name AS trim_name,
    t.base_price,
    t.is_available,
    g.generation_id,
    g.name AS generation_name,
    g.year_from,
    g.year_to,
    m.model_id,
    m.name AS model_name,
    m.segment,
    m.description AS model_description,
    b.brand_id,
    b.name AS brand_name,
    b.country AS brand_country,
    et.engine_type_id,
    et.name AS engine_type_name,
    tr.transmission_id,
    tr.name AS transmission_name,
    dt.drive_type_id,
    dt.name AS drive_type_name
FROM trims t
JOIN generations g ON t.generation_id = g.generation_id
JOIN models m ON g.model_id = m.model_id
JOIN brands b ON m.brand_id = b.brand_id
JOIN engine_types et ON t.engine_type_id = et.engine_type_id
JOIN transmissions tr ON t.transmission_id = tr.transmission_id
JOIN drive_types dt ON t.drive_type_id = dt.drive_type_id;