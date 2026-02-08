/**
 * SEO utility functions for managing page meta tags
 */

/**
 * Update page title
 * @param {string} title - Page title
 * @param {string} siteName - Site name (optional)
 */
export function updateTitle(title, siteName = '') {
  if (siteName) {
    document.title = `${title} - ${siteName}`
  } else {
    document.title = title
  }
}

/**
 * Update or create a meta tag
 * @param {string} name - Meta tag name or property
 * @param {string} content - Meta tag content
 * @param {boolean} isProperty - Whether to use property instead of name
 */
export function updateMetaTag(name, content, isProperty = false) {
  const attr = isProperty ? 'property' : 'name'
  let element = document.querySelector(`meta[${attr}="${name}"]`)
  
  if (!element) {
    element = document.createElement('meta')
    element.setAttribute(attr, name)
    document.head.appendChild(element)
  }
  
  element.setAttribute('content', content)
}

/**
 * Update all SEO meta tags for a page
 * @param {Object} options - SEO options
 * @param {string} options.title - Page title
 * @param {string} options.description - Page description
 * @param {string} options.keywords - Page keywords (comma separated)
 * @param {string} options.url - Page URL
 * @param {string} options.image - Page image URL
 * @param {string} options.type - Open Graph type (default: website)
 * @param {string} options.siteName - Site name
 */
export function updatePageSEO(options) {
  const {
    title,
    description,
    keywords,
    url,
    image,
    type = 'website',
    siteName = ''
  } = options

  // Update title
  if (title) {
    updateTitle(title, siteName)
  }

  // Update description
  if (description) {
    updateMetaTag('description', description)
    updateMetaTag('og:description', description, true)
    updateMetaTag('twitter:description', description)
  }

  // Update keywords
  if (keywords) {
    updateMetaTag('keywords', keywords)
  }

  // Update Open Graph tags
  if (title) {
    const fullTitle = siteName ? `${title} - ${siteName}` : title
    updateMetaTag('og:title', fullTitle, true)
    updateMetaTag('twitter:title', fullTitle)
  }

  if (url) {
    updateMetaTag('og:url', url, true)
  }

  if (image) {
    updateMetaTag('og:image', image, true)
    updateMetaTag('twitter:image', image)
  }

  if (type) {
    updateMetaTag('og:type', type, true)
  }

  // Twitter card
  updateMetaTag('twitter:card', 'summary_large_image')
}

/**
 * Get default SEO data for different page types
 * @param {string} pageType - Type of page
 * @param {Object} data - Additional data for the page
 * @param {Object} siteSettings - Site settings
 * @returns {Object} SEO data
 */
export function getPageSEO(pageType, data = {}, siteSettings = {}) {
  const siteName = siteSettings.site_name || 'GitStore'
  const baseUrl = window.location.origin
  const defaultDescription = siteSettings.site_subtitle || 'Premium Discourse and Flarum plugins with permanent licenses'
  
  const seoData = {
    home: {
      title: siteName,
      description: defaultDescription,
      keywords: 'discourse plugins, flarum plugins, forum plugins, marketplace, git store',
      url: baseUrl,
      type: 'website'
    },
    plugins: {
      title: 'All Plugins',
      description: 'Browse our collection of premium Discourse and Flarum plugins. Permanent licenses with professional support.',
      keywords: 'discourse plugins, flarum plugins, plugin marketplace, forum extensions',
      url: `${baseUrl}/plugins`,
      type: 'website'
    },
    'plugin-detail': {
      title: data.name || 'Plugin Details',
      description: data.short_description || data.description || 'View plugin details, features, and purchase options',
      keywords: `${data.name || 'plugin'}, discourse, flarum, forum plugin, ${data.category || ''}`,
      url: `${baseUrl}/plugins/${data.slug || ''}`,
      image: data.icon_url || data.image_url,
      type: 'product'
    },
    dashboard: {
      title: 'Dashboard',
      description: 'Manage your plugins, licenses, and orders',
      keywords: 'dashboard, my plugins, my licenses',
      url: `${baseUrl}/dashboard`,
      type: 'website'
    },
    licenses: {
      title: 'My Licenses',
      description: 'View and manage your plugin licenses',
      keywords: 'licenses, plugin licenses, my purchases',
      url: `${baseUrl}/licenses`,
      type: 'website'
    },
    'license-detail': {
      title: data.plugin_name ? `${data.plugin_name} License` : 'License Details',
      description: `View license details and installation instructions for ${data.plugin_name || 'your plugin'}`,
      keywords: 'license, plugin license, installation',
      url: `${baseUrl}/licenses/${data.id || ''}`,
      type: 'website'
    },
    orders: {
      title: 'My Orders',
      description: 'View your order history and payment details',
      keywords: 'orders, order history, purchases',
      url: `${baseUrl}/orders`,
      type: 'website'
    },
    purchase: {
      title: data.plugin_name ? `Purchase ${data.plugin_name}` : 'Purchase Plugin',
      description: `Purchase ${data.plugin_name || 'this plugin'} with a permanent license`,
      keywords: 'purchase, buy plugin, payment',
      url: `${baseUrl}/purchase/${data.plugin_id || ''}`,
      type: 'website'
    },
    page: {
      title: data.title || 'Page',
      description: data.description || data.excerpt || 'Custom page',
      keywords: data.keywords || 'information, page',
      url: `${baseUrl}/pages/${data.slug || ''}`,
      type: 'article'
    },
    'not-found': {
      title: '404 Not Found',
      description: 'The page you are looking for does not exist',
      keywords: '404, not found',
      url: window.location.href,
      type: 'website'
    }
  }

  return {
    ...seoData[pageType],
    siteName
  }
}
