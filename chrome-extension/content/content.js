// Content script for webpage extraction
(function() {
  'use strict';

  // Listen for messages from popup
  chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
    if (request.action === 'extractContent') {
      const extractedContent = extractMainContent();
      sendResponse(extractedContent);
    }
    return true; // Keep the message channel open for async response
  });

  function extractMainContent() {
    try {
      const title = document.title;
      const url = window.location.href;
      
      // Remove unwanted elements temporarily
      const unwantedSelectors = [
        'script', 'style', 'noscript', 'iframe', 'nav', 'header', 'footer',
        '.advertisement', '.ads', '.sidebar', '.menu', '.navigation',
        '.social', '.comments', '.popup', '.modal'
      ];
      
      const removedElements = [];
      unwantedSelectors.forEach(selector => {
        document.querySelectorAll(selector).forEach(el => {
          removedElements.push({
            element: el,
            parent: el.parentNode,
            nextSibling: el.nextSibling
          });
          el.remove();
        });
      });

      // Try to find main content area
      let mainContent = '';
      
      // Priority order for content selectors
      const contentSelectors = [
        'main', 
        'article', 
        '[role="main"]',
        '.content',
        '.main-content',
        '.post-content',
        '.entry-content',
        '.article-body',
        '#content',
        '#main'
      ];

      let contentElement = null;
      for (const selector of contentSelectors) {
        const element = document.querySelector(selector);
        if (element && element.textContent.trim().length > 100) {
          contentElement = element;
          break;
        }
      }

      // Fallback to body if no specific content element found
      if (!contentElement) {
        contentElement = document.body;
      }

      // Extract text content
      const textContent = contentElement.textContent || contentElement.innerText || '';
      mainContent = textContent.trim();

      // Restore removed elements
      removedElements.forEach(({ element, parent, nextSibling }) => {
        if (parent && nextSibling) {
          parent.insertBefore(element, nextSibling);
        } else if (parent) {
          parent.appendChild(element);
        }
      });

      // Clean up the content
      mainContent = cleanContent(mainContent);

      // Extract metadata
      const metadata = extractMetadata();

      return {
        title,
        url,
        content: mainContent.substring(0, 10000), // Limit content length
        metadata
      };
    } catch (error) {
      console.error('Content extraction failed:', error);
      return {
        title: document.title,
        url: window.location.href,
        content: '',
        metadata: {},
        error: error.message
      };
    }
  }

  function cleanContent(content) {
    return content
      .replace(/\s+/g, ' ')           // Replace multiple whitespace with single space
      .replace(/\n\s*\n/g, '\n')       // Replace multiple newlines with single newline
      .replace(/^\s+|\s+$/g, '')      // Trim leading and trailing whitespace
      .substring(0, 10000);           // Limit content length
  }

  function extractMetadata() {
    const metadata = {};

    // Extract description
    const descriptionMeta = document.querySelector('meta[name="description"]');
    if (descriptionMeta) {
      metadata.description = descriptionMeta.getAttribute('content');
    }

    // Extract author
    const authorMeta = document.querySelector('meta[name="author"]');
    if (authorMeta) {
      metadata.author = authorMeta.getAttribute('content');
    }

    // Extract published date
    const publishedMeta = document.querySelector('meta[name="article:published_time"], meta[property="article:published_time"]');
    if (publishedMeta) {
      metadata.publishedDate = publishedMeta.getAttribute('content');
    }

    // Extract language
    const langMeta = document.querySelector('html');
    if (langMeta) {
      metadata.language = langMeta.getAttribute('lang');
    }

    // Extract word count
    const textContent = document.body.textContent || document.body.innerText || '';
    metadata.wordCount = textContent.trim().split(/\s+/).length;

    return metadata;
  }

  // Optional: Add a floating action button for quick access
  function addFloatingButton() {
    // Check if button already exists
    if (document.querySelector('#chromllm-fab')) {
      return;
    }

    const button = document.createElement('div');
    button.id = 'chromllm-fab';
    button.innerHTML = '🤖';
    button.style.cssText = `
      position: fixed;
      bottom: 20px;
      right: 20px;
      width: 50px;
      height: 50px;
      border-radius: 50%;
      background-color: #3498db;
      color: white;
      border: none;
      cursor: pointer;
      font-size: 20px;
      display: flex;
      align-items: center;
      justify-content: center;
      box-shadow: 0 2px 10px rgba(0,0,0,0.3);
      z-index: 10000;
      transition: all 0.3s ease;
    `;

    button.addEventListener('mouseenter', () => {
      button.style.transform = 'scale(1.1)';
      button.style.backgroundColor = '#2980b9';
    });

    button.addEventListener('mouseleave', () => {
      button.style.transform = 'scale(1)';
      button.style.backgroundColor = '#3498db';
    });

    button.addEventListener('click', () => {
      // Send message to background script to open popup
      chrome.runtime.sendMessage({ action: 'openPopup' });
    });

    document.body.appendChild(button);
  }

  // Add floating button after page loads
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', addFloatingButton);
  } else {
    addFloatingButton();
  }

})();