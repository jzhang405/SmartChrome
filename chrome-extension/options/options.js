class ChromLLMOptions {
  constructor() {
    this.defaultSettings = {
      backendURL: 'http://localhost:8080',
      maxContentLength: 10000,
      autoExtract: false,
      theme: 'light',
      debugMode: false
    };

    this.initializeElements();
    this.setupEventListeners();
    this.loadSettings();
    this.checkWelcomePage();
  }

  initializeElements() {
    this.backendURL = document.getElementById('backendURL');
    this.maxContentLength = document.getElementById('maxContentLength');
    this.autoExtract = document.getElementById('autoExtract');
    this.theme = document.getElementById('theme');
    this.debugMode = document.getElementById('debugMode');
    this.saveButton = document.getElementById('saveSettings');
    this.resetButton = document.getElementById('resetSettings');
    this.testConnectionButton = document.getElementById('testConnection');
    this.connectionStatus = document.getElementById('connectionStatus');
    this.statusMessage = document.getElementById('statusMessage');
    this.welcomeBanner = document.getElementById('welcomeBanner');
    this.exportDataLink = document.getElementById('exportData');
    this.clearDataLink = document.getElementById('clearData');
  }

  setupEventListeners() {
    this.saveButton.addEventListener('click', () => this.saveSettings());
    this.resetButton.addEventListener('click', () => this.resetSettings());
    this.testConnectionButton.addEventListener('click', () => this.testConnection());
    this.exportDataLink.addEventListener('click', (e) => {
      e.preventDefault();
      this.exportData();
    });
    this.clearDataLink.addEventListener('click', (e) => {
      e.preventDefault();
      this.clearData();
    });
  }

  async checkWelcomePage() {
    const urlParams = new URLSearchParams(window.location.search);
    if (urlParams.get('welcome') === 'true') {
      this.welcomeBanner.style.display = 'block';
    }
  }

  async loadSettings() {
    try {
      const result = await chrome.storage.local.get(['settings']);
      const settings = result.settings || this.defaultSettings;
      
      this.backendURL.value = settings.backendURL || this.defaultSettings.backendURL;
      this.maxContentLength.value = settings.maxContentLength || this.defaultSettings.maxContentLength;
      this.autoExtract.checked = settings.autoExtract || this.defaultSettings.autoExtract;
      this.theme.value = settings.theme || this.defaultSettings.theme;
      this.debugMode.checked = settings.debugMode || this.defaultSettings.debugMode;
    } catch (error) {
      console.error('Failed to load settings:', error);
      this.showStatus('Failed to load settings', 'error');
    }
  }

  async saveSettings() {
    try {
      const settings = {
        backendURL: this.backendURL.value.trim(),
        maxContentLength: parseInt(this.maxContentLength.value),
        autoExtract: this.autoExtract.checked,
        theme: this.theme.value,
        debugMode: this.debugMode.checked
      };

      // Validate settings
      if (!this.validateSettings(settings)) {
        return;
      }

      await chrome.storage.local.set({ settings });
      this.showStatus('Settings saved successfully!', 'success');
      
      // Test connection after saving
      setTimeout(() => this.testConnection(), 1000);
    } catch (error) {
      console.error('Failed to save settings:', error);
      this.showStatus('Failed to save settings', 'error');
    }
  }

  validateSettings(settings) {
    if (!settings.backendURL) {
      this.showStatus('Backend URL is required', 'error');
      return false;
    }

    try {
      new URL(settings.backendURL);
    } catch {
      this.showStatus('Invalid Backend URL', 'error');
      return false;
    }

    if (settings.maxContentLength < 1000 || settings.maxContentLength > 50000) {
      this.showStatus('Maximum content length must be between 1000 and 50000', 'error');
      return false;
    }

    return true;
  }

  async resetSettings() {
    if (confirm('Are you sure you want to reset all settings to defaults?')) {
      try {
        await chrome.storage.local.set({ settings: this.defaultSettings });
        await this.loadSettings();
        this.showStatus('Settings reset to defaults', 'success');
      } catch (error) {
        console.error('Failed to reset settings:', error);
        this.showStatus('Failed to reset settings', 'error');
      }
    }
  }

  async testConnection() {
    this.testConnectionButton.disabled = true;
    this.connectionStatus.textContent = 'Testing...';

    try {
      const response = await fetch(`${this.backendURL.value}/v1/health`);
      
      if (response.ok) {
        const health = await response.json();
        this.connectionStatus.textContent = '✅ Connected';
        this.connectionStatus.style.color = 'green';
      } else {
        this.connectionStatus.textContent = '❌ Connection failed';
        this.connectionStatus.style.color = 'red';
      }
    } catch (error) {
      console.error('Connection test failed:', error);
      this.connectionStatus.textContent = '❌ Connection failed';
      this.connectionStatus.style.color = 'red';
    } finally {
      this.testConnectionButton.disabled = false;
    }
  }

  showStatus(message, type) {
    this.statusMessage.textContent = message;
    this.statusMessage.className = `status ${type}`;
    this.statusMessage.style.display = 'block';

    setTimeout(() => {
      this.statusMessage.style.display = 'none';
    }, 3000);
  }

  async exportData() {
    try {
      const result = await chrome.storage.local.get(null);
      const dataStr = JSON.stringify(result, null, 2);
      const dataBlob = new Blob([dataStr], { type: 'application/json' });
      
      const url = URL.createObjectURL(dataBlob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `chromllm-backup-${new Date().toISOString().split('T')[0]}.json`;
      link.click();
      
      URL.revokeObjectURL(url);
      this.showStatus('Data exported successfully', 'success');
    } catch (error) {
      console.error('Failed to export data:', error);
      this.showStatus('Failed to export data', 'error');
    }
  }

  async clearData() {
    if (confirm('Are you sure you want to clear all extension data? This action cannot be undone.')) {
      try {
        await chrome.storage.local.clear();
        await this.loadSettings();
        this.showStatus('All data cleared successfully', 'success');
      } catch (error) {
        console.error('Failed to clear data:', error);
        this.showStatus('Failed to clear data', 'error');
      }
    }
  }
}

// Initialize the options page when the DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
  new ChromLLMOptions();
});