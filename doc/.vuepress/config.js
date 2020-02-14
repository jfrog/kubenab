const container = require('markdown-it-container')

const ogprefix = 'og: http://ogp.me/ns#'
const title = 'kubenab'
const description = 'Kubernetes Admission Webhook to enforce pulling of Docker images from the private registry.'
const color = '#2F80ED'
//const author = '§author§'
//const url = 'https://docs.fab/§prjid§'

module.exports = {
  plugins: [
    ['@vuepress/search', {
      searchMaxSuggestions: 10
    }]
  ],
  head: [
    ['link', { rel: 'icon', href: `/icon.png` }],
    ['meta', { name: 'theme-color', content: color }],
    ['meta', { prefix: ogprefix, property: 'og:title', content: title }],
    ['meta', { prefix: ogprefix, property: 'twitter:title', content: title }],
    ['meta', { prefix: ogprefix, property: 'og:type', content: 'article' }],
    ['meta', { prefix: ogprefix, property: 'og:url', content: url }],
    ['meta', { prefix: ogprefix, property: 'og:description', content: description }],
    ['meta', { prefix: ogprefix, property: 'og:image', content: `${url}icon.png` }],
    ['meta', { prefix: ogprefix, property: 'og:article:author', content: author }],
    ['meta', { name: 'apple-mobile-web-app-capable', content: 'yes' }],
    ['meta', { name: 'apple-mobile-web-app-status-bar-style', content: 'black' }],
    // ['link', { rel: 'apple-touch-icon', href: `/assets/apple-touch-icon.png` }],
    // ['link', { rel: 'mask-icon', href: '/assets/safari-pinned-tab.svg', color: color }],
    // ['meta', { name: 'msapplication-TileImage', content: '/icon.png' }],
    // ['meta', { name: 'msapplication-TileColor', content: color }],
  ],
  markdown: {
    anchor: {
      permalink: true,
    },
    config: md => {
      md
        .use(require('markdown-it-container'))
        .use(require('markdown-it-mathjax'))
        .use(require('markdown-it-attrs'))
        .use(require('markdown-it-checkbox'), {divWrap: true, divClass: 'cb', idPrefix: 'cbx_'})
        .use(require('markdown-it-mark'))
        .use(require('markdown-it-footnote'))
        .use(require('markdown-it-sup'))
        .use(require('markdown-it-anchor'), {permalink: true, permalinkBefore: true, permalinkSymbol: '§'})
        .use(require('markdown-it-toc-done-right'), {"placeholder": "[[toc]]"})
        .use(require('markdown-it-decorate'))
        .use(...createContainer('intro'))
        .use(...createContainer('note'))
    }
  },
  title,
  description,
  base: '/',
  themeConfig: {
    versions: [
      // VERSIONS
    ],
    docsDir: 'doc',
    serviceWorker: true,
    hiddenLinks: [
      // HIDDEN_LINKS
    ],
    sidebar: {
            // add links to documentation
    },
  },
}


/**
  * Function copied from Strapi (https://github.com/strapi/strapi/blob/master/docs/.vuepress/config.js#L179)
  */
function createContainer(className) {
  return [container, className, {
    render(tokens, idx) {
      const token = tokens[idx]
      if (token.nesting === 1) {
        return `<div class="${className} custom-block">\n`
      } else {
        return `</div>\n`
      }
    }
  }]
}
