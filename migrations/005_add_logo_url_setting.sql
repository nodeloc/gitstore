-- Add logo URL setting
INSERT INTO system_settings (key, value, description) VALUES
    ('logo_url', '', 'Logo image URL for site header and favicon')
ON CONFLICT (key) DO NOTHING;
