-- migrations/003_create_pages_table.sql
-- Create pages table for static content management

-- Pages table for About, Contact, Privacy Policy, Terms of Service, etc.
CREATE TABLE pages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    slug VARCHAR(255) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'published')),
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_pages_slug ON pages(slug);
CREATE INDEX idx_pages_status ON pages(status);
CREATE INDEX idx_pages_sort_order ON pages(sort_order);

-- Insert default pages
INSERT INTO pages (slug, title, content, status) VALUES
('about', 'About Us', '# About Us

Welcome to our plugin store.

## Who We Are

We are dedicated to providing high-quality plugins for your projects.

## Our Mission

To make development easier and more efficient.', 'published'),

('contact', 'Contact', '# Contact Us

Get in touch with us.

## Email
support@example.com

## Office Hours
Monday - Friday: 9:00 AM - 5:00 PM (UTC)', 'published'),

('privacy-policy', 'Privacy Policy', '# Privacy Policy

**Last updated: ' || TO_CHAR(CURRENT_DATE, 'Month DD, YYYY') || '**

## Information We Collect

We collect information you provide directly to us when you:
- Create an account
- Make a purchase
- Contact us

## How We Use Your Information

We use the information we collect to:
- Provide and maintain our services
- Process your transactions
- Send you technical notices and updates

## Data Security

We implement appropriate security measures to protect your personal information.

## Contact Us

If you have questions about this Privacy Policy, please contact us at privacy@example.com', 'published'),

('terms-of-service', 'Terms of Service', '# Terms of Service

**Last updated: ' || TO_CHAR(CURRENT_DATE, 'Month DD, YYYY') || '**

## Agreement to Terms

By accessing our service, you agree to be bound by these Terms.

## Use License

Permission is granted to temporarily download one copy of the materials for personal, non-commercial transitory viewing only.

## Disclaimer

The materials on our website are provided on an ''as is'' basis. We make no warranties, expressed or implied.

## Limitations

In no event shall we or our suppliers be liable for any damages arising out of the use or inability to use our materials.

## Modifications

We may revise these terms of service at any time without notice.

## Contact

For questions about these Terms, contact us at legal@example.com', 'published');
