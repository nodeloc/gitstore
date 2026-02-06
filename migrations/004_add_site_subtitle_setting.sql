-- Add site subtitle setting
INSERT INTO system_settings (key, value, description) VALUES
    ('site_subtitle', 'Plugin Marketplace', 'Site subtitle displayed in header')
ON CONFLICT (key) DO NOTHING;
