// Background service worker for ChromLLM extension
class ChromLLMBackground {
  constructor() {
    this.setupEventListeners();
    this.initializeStorage();
  }

  setupEventListeners() {
    // Handle extension installation
    chrome.runtime.onInstalled.addListener((details) => {
      this.handleInstallation(details);
    });

    // Handle messages from content scripts and popup
    chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
      this.handleMessage(request, sender, sendResponse);
      return true; // Keep message channel open for async responses
    });

    // Handle extension startup
    chrome.runtime.onStartup.addListener(() => {
      this.handleStartup();
    });

    // Handle tab updates
    chrome.tabs.onUpdated.addListener((tabId, changeInfo, tab) => {
      if (changeInfo.status === 'complete' && tab.url) {
        this.handleTabUpdate(tabId, tab);
      }
    });
  }

  async initializeStorage() {
    try {
      // Set default settings if not already set
      const result = await chrome.storage.local.get(['settings']);
      if (!result.settings) {
        const defaultSettings = {
          backendURL: 'http://localhost:8080',
          maxContentLength: 10000,
          autoExtract: false,
          theme: 'light'
        };
        await chrome.storage.local.set({ settings: defaultSettings });
      }
    } catch (error) {
      console.error('Failed to initialize storage:', error);
    }
  }

  async handleInstallation(details) {
    console.log('Extension installed:', details.reason);
    
    if (details.reason === 'install') {
      // Open welcome page on first install
      chrome.tabs.create({
        url: chrome.runtime.getURL('options.html?welcome=true')
      });
    }

    // Set up initial storage
    await this.initializeStorage();
  }

  handleStartup() {
    console.log('Extension started');
    // Perform any startup tasks here
  }

  async handleMessage(request, sender, sendResponse) {
    try {
      switch (request.action) {
        case 'openPopup':
          await this.openPopup();
          sendResponse({ success: true });
          break;

        case 'extractContent':
          const content = await this.extractContentFromTab(sender.tab);
          sendResponse(content);
          break;

        case 'getSettings':
          const settings = await this.getSettings();
          sendResponse(settings);
          break;

        case 'updateSettings':
          await this.updateSettings(request.settings);
          sendResponse({ success: true });
          break;

        case 'checkBackendHealth':
          const health = await this.checkBackendHealth();
          sendResponse(health);
          break;

        case 'createSession':
          const session = await this.createSession();
          sendResponse(session);
          break;

        case 'createConversation':
          const conversation = await this.createConversation(request.data);
          sendResponse(conversation);
          break;

        case 'sendMessage':
          const message = await this.sendMessage(request.data);
          sendResponse(message);
          break;

        default:
          console.warn('Unknown action:', request.action);
          sendResponse({ error: 'Unknown action' });
      }
    } catch (error) {
      console.error('Error handling message:', error);
      sendResponse({ error: error.message });
    }
  }

  async openPopup() {
    try {
      // Get the current active tab
      const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
      
      // Open the extension popup
      chrome.action.openPopup();
    } catch (error) {
      console.error('Failed to open popup:', error);
    }
  }

  async extractContentFromTab(tab) {
    try {
      if (!tab) {
        throw new Error('No active tab found');
      }

      // Inject content script if not already injected
      await chrome.scripting.executeScript({
        target: { tabId: tab.id },
        files: ['content/content.js']
      });

      // Send message to content script
      const response = await chrome.tabs.sendMessage(tab.id, {
        action: 'extractContent'
      });

      return response;
    } catch (error) {
      console.error('Failed to extract content:', error);
      return {
        error: error.message,
        title: tab ? tab.title : 'Unknown',
        url: tab ? tab.url : '',
        content: ''
      };
    }
  }

  async getSettings() {
    try {
      const result = await chrome.storage.local.get(['settings']);
      return result.settings || {};
    } catch (error) {
      console.error('Failed to get settings:', error);
      return {};
    }
  }

  async updateSettings(settings) {
    try {
      await chrome.storage.local.set({ settings });
    } catch (error) {
      console.error('Failed to update settings:', error);
      throw error;
    }
  }

  async checkBackendHealth() {
    try {
      const settings = await this.getSettings();
      const response = await fetch(`${settings.backendURL}/v1/health`);
      
      if (response.ok) {
        const health = await response.json();
        return { status: 'healthy', ...health };
      } else {
        return { status: 'unhealthy', error: 'Backend returned error' };
      }
    } catch (error) {
      console.error('Backend health check failed:', error);
      return { status: 'unhealthy', error: error.message };
    }
  }

  async createSession() {
    try {
      const settings = await this.getSettings();
      const response = await fetch(`${settings.backendURL}/v1/sessions`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({})
      });

      if (response.ok) {
        const session = await response.json();
        // Store session data
        await chrome.storage.local.set({
          sessionID: session.session.id,
          authToken: session.token
        });
        return session;
      } else {
        throw new Error('Failed to create session');
      }
    } catch (error) {
      console.error('Session creation failed:', error);
      throw error;
    }
  }

  async createConversation(data) {
    try {
      const settings = await this.getSettings();
      const result = await chrome.storage.local.get(['authToken']);
      
      const response = await fetch(`${settings.backendURL}/v1/conversations`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${result.authToken}`
        },
        body: JSON.stringify(data)
      });

      if (response.ok) {
        const conversation = await response.json();
        await chrome.storage.local.set({ conversationID: conversation.id });
        return conversation;
      } else {
        throw new Error('Failed to create conversation');
      }
    } catch (error) {
      console.error('Conversation creation failed:', error);
      throw error;
    }
  }

  async sendMessage(data) {
    try {
      const settings = await this.getSettings();
      const result = await chrome.storage.local.get(['authToken']);
      
      const response = await fetch(`${settings.backendURL}/v1/conversations/${data.conversationId}/messages`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${result.authToken}`
        },
        body: JSON.stringify(data)
      });

      if (response.ok) {
        return await response.json();
      } else {
        throw new Error('Failed to send message');
      }
    } catch (error) {
      console.error('Message sending failed:', error);
      throw error;
    }
  }

  handleTabUpdate(tabId, tab) {
    // Handle tab updates if needed
    // For example, clear conversation context when navigating to a new page
    if (tab.url && tab.url !== 'about:blank' && tab.url !== 'chrome://newtab/') {
      // Optionally clear conversation data when tab changes
      // chrome.storage.local.remove(['conversationID', 'extractedContent']);
    }
  }
}

// Initialize the background service
const chromllmBackground = new ChromLLMBackground();

// Export for testing (if needed)
if (typeof module !== 'undefined' && module.exports) {
  module.exports = ChromLLMBackground;
}