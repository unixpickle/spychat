class LoadableView {
  constructor(elementId) {
    this.element = document.getElementById(elementId);
    this.error = this.element.getElementsByClassName('error')[0];
    this.abortController = null;
  }

  loadPath(path) {
    if (this.abortController) {
      this.abortController.abort();
    }
    this.abortController = new AbortController();
    fetch(path, { signal: this.abortController.signal })
      .then((response) => response.json())
      .then((jsonResponse) => {
        this.abortController = null;
        if (jsonResponse.error) {
          this.showError(jsonResponse.error);
        } else {
          this.showResult(jsonResponse.result);
        }
      })
      .catch(() => {
        this.abortController = null;
        this.showError('request failed')
      });
    this.showLoader();
  }

  showError(message) {
    this.element.className = 'failed';
    this.error.textContent = message;
  }

  showLoader() {
    this.element.className = 'loading';
  }

  showResult(result) {
    this.element.className = 'loaded';
  }
}