/** @type {import('@eventcatalog/core/bin/eventcatalog.config').Config} */
export default {
  title: 'Mills Mess',
  tagline: 'Licensed under the Mills Mess License Agreement. Free for non-commercial use or first earned $500K, then 5% and 3% after $2M. See LICENSE.md in the root of this repository',
  organizationName: 'Staring, LLC',
  homepageLink: 'https://eventcatalog.dev/',
  editUrl: 'https://github.com/staringfun/millsmess',
  // Supports static or server. Static renders a static site, server renders a server side rendered site
  // large catalogs may benefit from server side rendering
  output: 'static',
  // By default set to false, add true to get urls ending in /
  trailingSlash: false,
  // Change to make the base url of the site different, by default https://{website}.com/docs,
  // changing to /company would be https://{website}.com/company/docs,
  base: '/',
  // Customize the logo, add your logo to public/ folder
  logo: {
    alt: 'Mills Mess Logo',
    src: '/logo.png',
    text: 'Mills Mess'
  },
  // Enable RSS feed for your eventcatalog
  rss: {
    enabled: false,
    // number of items to include in the feed per resource (event, service, etc)
    // limit: 20
  },
  // This lets you copy markdown contents from EventCatalog to your clipboard
  // Including schemas for your events and services
  llmsTxt: {
    enabled: true,
  },
  docs: {
    sidebar: {
      // TREE_VIEW will render the DOCS as a tree view and map your file system folder structure
      // LIST_VIEW will render the DOCS that look familiar to API documentation websites
      type: 'LIST_VIEW'
    },
  },
  generators: [
    [
      './go.js',
      {
        outPath: '../libs/types/generated.go',
        schemaPath: './schemas/root.json',
        preamble: [
            '// Mills Mess',
            '// Licensed under the Mills Mess License Agreement',
            '// See LICENSE.md in the root of this repository.',
            '',
        ]
      }
    ],
  ],
  sidebar: [{
    id: '/docs/custom',
    visible: false
  }, {
    id: '/chat',
    visible: false
  }],
  // required random generated id used by eventcatalog
  cId: 'cbface1a-347f-4495-9c44-9c2ade59c428'
}
